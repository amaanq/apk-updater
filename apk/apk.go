package apk

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/amaanq/sc-compression"
	"github.com/go-resty/resty/v2"
	cp "github.com/otiai10/copy"
	"github.com/withmandala/go-log"
	"golang.org/x/net/html"
)

const (
	BaseURL = "https://clash-of-clans.en.uptodown.com/android/download"
)

var (
	CurrentVersion string

	Path = flag.String("path", ".", "output directory for /assets")

	Client = resty.New()
	Log    = log.New(os.Stdout).WithColor().WithDebug().WithTimestamp()

	ValidDirectories = []string{"csv", "localization", "logic" /*"sc"*/} // Uncomment sc if you have a lot of RAM (>8GB or >4GB free)
)

func LoadCurrentVersion() {
	CurrentVersion = os.Getenv("game")
	if CurrentVersion == "" {
		panic("Please set game to the game version you'd like to view the decompressed assets of; if you haven't moved dotenv to .env please do so.")
	}
}

func FixPath() {
	if !strings.HasSuffix(*Path, "/") {
		*Path += "/"
	}
}

func UpdateAPK() error {
	LoadCurrentVersion()

	node, err := CurlAPKLink()
	if err != nil {
		Log.Error(err)
	}

	Log.Info("Checking version...")
	version, err := GetAPKVersion(node)
	if err != nil {
		Log.Error(err)
		return err
	}

	if version == CurrentVersion {
		Log.Info("You are up to date!")
		return nil
	}
	CurrentVersion = version

	Log.Info("New game version available! (" + version + ")")
	url, err := GetDownloadURL(node)
	if err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Downloading APK!")
	if err = WgetAPK(url, version); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Removing Potential Base Path Collision")
	os.RemoveAll("clash-" + version)

	Log.Info("Decompiling APK!")
	if err = DecompileAPK("clash-" + version + ".apk"); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Removing Potential Assets Path Collision")
	os.RemoveAll(*Path + "/assets" + CurrentVersion)

	Log.Info("Moving assets folder outside...")
	if err = ExtractAssets("clash-" + version); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Decompressing assets...")
	if err = WalkAndDecompressAssets(); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Done!")

	return nil
}

// Walk the assets folder and decompress each file inside
func WalkAndDecompressAssets() error {
	// Create new dir for decompressed assets
	os.RemoveAll("decompressed" + CurrentVersion)
	err := os.Mkdir("decompressed"+CurrentVersion, 0775)
	if err != nil {
		return err
	}

	for _, subdir := range ValidDirectories {
		dir := *Path + "/assets" + CurrentVersion + "/" + subdir + "/"
		Log.Info("Reading dir", dir)
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		err = os.Mkdir("decompressed"+CurrentVersion+"/"+subdir, 0775)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fileName := entry.Name()
			fullPath := dir + fileName

			Log.Info("Decompressing", fullPath)
			compFile, err := os.Open(fullPath)
			if err != nil {
				return err
			}

			decompressor := ScCompression.NewDecompressor(compFile)
			reader, err := decompressor.Decompress()
			if err != nil {
				return err
			}

			fd, err := os.Create("./decompressed" + CurrentVersion + "/" + subdir + "/" + fileName)
			if err != nil {
				return err
			}

			if _, err = io.Copy(fd, reader); err != nil {
				return err
			}
			if err = fd.Close(); err != nil {
				return err
			}
		} 
	}
	return nil
}

// Parses the uptodown HTML node for the current game version
func GetAPKVersion(node *html.Node) (string, error) {
	query := goquery.NewDocumentFromNode(node)
	var version string
	query.Find(`script[type="application/ld+json"]`).Each(func(i int, script *goquery.Selection) {
		if strings.Contains(script.Text(), "softwareVersion") {
			var metadata MetaData
			err := json.Unmarshal([]byte(script.Text()), &metadata)
			if err != nil {
				Log.Error(err)
			}
			version = metadata.MainEntity.SoftwareVersion
		}
	})
	if version == "" {
		return version, errors.New("couldn't find the version")
	}
	return version, nil
}

// Parses the uptodown HTML node for the download link
func GetDownloadURL(node *html.Node) (string, error) {
	query := goquery.NewDocumentFromNode(node)
	var url string
	query.Find("a.button.download").Each(func(i int, s *goquery.Selection) {
		n, ok := s.Attr("href")
		if ok {
			url = n
		}
	})
	if url == "" {
		return url, errors.New("couldn't find the version")
	}
	return url, nil
}

// Executes apktool and removes the apk file
func DecompileAPK(apkPath string) error {
	err := exec.Command("apktool", "d", apkPath).Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(apkPath)
	return err
}

// Moves assets folder inside apk to project root directory
func ExtractAssets(path string) error {
	var err error
	if *Path == "." {
		__path := "./assets" + CurrentVersion
		os.RemoveAll(__path)
		Log.Infof("Moving %s to %s", path+"/assets", __path)
		err = cp.Copy(path+"/assets", __path)
	} else {
		__path := *Path + "/assets" + CurrentVersion
		os.RemoveAll(__path)
		err = cp.Copy(path+"/assets", __path)
	}

	if err != nil {
		return err
	}

	err = os.RemoveAll(path)
	return err
}

// Get uptodowns HTML page
func CurlAPKLink() (*html.Node, error) {
	resp, err := Client.R().Get(BaseURL)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(resp.Body())))
	if err != nil {
		return nil, err
	}
	return doc, err
}

// Download the APK from uptodown
func WgetAPK(url, version string) error {
	resp, err := Client.R().Get(url)
	if err != nil {
		return err
	}

	fd, err := os.Create("clash-" + version + ".apk")
	if err != nil {
		return err
	}

	_, err = fd.Write(resp.Body())
	return err
}
