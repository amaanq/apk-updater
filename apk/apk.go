package apk

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	// "github.com/amaanq/sc-compression"
	"github.com/go-resty/resty/v2"
	"github.com/withmandala/go-log"
	"golang.org/x/net/html"
)

const (
	BaseURL        = "https://clash-of-clans.en.uptodown.com/android/download"
	CurrentVersion = "14.426.0"
)

var (
	Path = flag.String("path", ".", "output directory for /assets")

	Client = resty.New()
	Log    = log.New(os.Stdout).WithColor().WithDebug().WithTimestamp()

	ValidDirectories = []string{"csv", "image", "localization", "logic", "SC", ""}
)

func FixPath() {
	if !strings.HasSuffix(*Path, "/") {
		*Path += "/"
	}
}

func CheckForNewVersion() error {
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
	if version != CurrentVersion {

		Log.Info("New game version detected!")
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
		RemovePotentialCollision("clash-" + version)

		Log.Info("Decompiling APK!")
		if err = DecompileAPK("clash-" + version + ".apk"); err != nil {
			Log.Error(err)
			return err
		}

		Log.Info("Removing Potential Assets Path Collision")
		RemovePotentialCollision(*Path + "/assets")

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
	}
	return nil
}

func WalkAndDecompressAssets() error {
	for _, subdir := range ValidDirectories {
		dir := *Path + subdir
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		fmt.Println(entries)
	}

	return nil
}

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
		return version, fmt.Errorf("couldn't find the version")
	}
	return version, nil
}

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
		return url, fmt.Errorf("couldn't find the version")
	}
	return url, nil
}

func DecompileAPK(apkPath string) error {
	err := exec.Command("apktool", "d", apkPath).Run()
	if err != nil {
		return err
	}
	err = exec.Command("rm", apkPath).Run()
	if err != nil {
		return err
	}
	return nil
}

func ExtractAssets(path string) error {
	fmt.Println("mv", path+"/assets", *Path)
	err := exec.Command("mv", path+"/assets", *Path).Run()
	if err != nil {
		return err
	}
	err = exec.Command("rm", "-r", path).Run()
	if err != nil {
		return err
	}
	return nil
}

func RemovePotentialCollision(path string) {
	exec.Command("rm", "-r", path).Run()
}

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
	if err != nil {
		return err
	}
	return nil
}
