// Copyright (C) 2022 Amaan Qureshi (aq0527@pm.me)
//
//
// This file is a part of APK Updater.
//
//
// This project, APK Updater, is not to be redistributed or copied without
//
// the express permission of the copyright holder, Amaan Qureshi (amaanq).

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
	"github.com/hashicorp/go-retryablehttp"
	"github.com/withmandala/go-log"
	"golang.org/x/net/html"
)

var (
	CurrentVersion string

	Path = flag.String("path", ".", "output directory for /assets")

	//Client = http.Client{Timeout: time.Second * 15}
	Client = LoadRetryClient()

	Log = log.New(os.Stdout).WithColor().WithDebug().WithTimestamp()

	ValidDirectories = []string{"csv", "localization", "logic" /*"sc"*/} // Uncomment sc if you have a lot of RAM (>8GB or >4GB free)
)

func LoadRetryClient() *retryablehttp.Client {
	Client := retryablehttp.NewClient()
	Client.Logger = nil
	Client.RetryMax = 5
	return Client
}

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

	version, err := GetCurrentAPKVersion(true)
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
	url, err := GetDownloadURL(ClashofClans.URL)
	if err != nil {
		Log.Error(err)
		return err
	}

	if _, err = WgetAPK(&ClashofClans, url, version, ""); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Removing Potential Base Path Collision")
	os.RemoveAll("clash-" + version)

	if err = DecompileAPK("clash-" + version + ".apk"); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Removing Potential Assets Path Collision")
	os.RemoveAll(*Path + "/assets" + CurrentVersion)

	Log.Info("Moving assets folder outside...")
	if err = ExtractAssets("clash-"+version, ""); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Decompressing assets...")
	if err = WalkAndDecompressAssets("clash-"+version, "decompressed"+version); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Done!")

	return nil
}

// Walk the assets folder and decompress each file inside
func WalkAndDecompressAssets(fpToDecompiledAPK, fpToOutputFiles string) error {
	os.RemoveAll(fpToOutputFiles)
	err := os.Mkdir(fpToOutputFiles, 0755)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, subdir := range ValidDirectories {
		Log.Info("Reading dir", subdir)
		entries, err := os.ReadDir("./" + fpToDecompiledAPK + "/assets/" + subdir + "/")
		if err != nil {
			return err
		}

		err = os.Mkdir(fpToOutputFiles+"/"+subdir, 0755)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			fileName := entry.Name()
			fullPath := cwd + "/" + fpToOutputFiles + "/" + fileName
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

			fd, err := os.Create(fpToOutputFiles + "/" + subdir + "/" + fileName)
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
func GetCurrentAPKVersion(_print bool) (string, error) {
	if _print {
		Log.Info("Checking version...")
	}
	node, err := CurlAPKLink(ClashofClans.URL)
	if err != nil {
		Log.Error(err)
	}

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
func GetDownloadURL(url string) (string, error) {
	node, err := CurlAPKLink(url)
	if err != nil {
		Log.Error(err)
	}

	query := goquery.NewDocumentFromNode(node)
	var downloadUrl string
	query.Find("a.button.download").Each(func(i int, s *goquery.Selection) {
		n, ok := s.Attr("href")
		if ok {
			downloadUrl = n
		}
	})
	if downloadUrl == "" {
		return downloadUrl, errors.New("couldn't find the version")
	}
	return downloadUrl, nil
}

// Executes apktool and removes the apk file
func DecompileAPK(apkPath string) error {
	Log.Info("Decompiling APK!")
	err := exec.Command("apktool", "d", apkPath, "-f").Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(apkPath)
	return err
}

// Moves assets folder inside apk to project root directory
func ExtractAssets(gamepath, assetpath string) error {

	// var err error
	// if *Path == "." {
	// 	__path := "./assets" + CurrentVersion
	// 	os.RemoveAll(__path)
	// 	Log.Infof("Moving %s to %s", path+"/assets", __path)
	// 	err = cp.Copy(path+"/assets", __path)
	// } else {
	// 	__path := *Path + "/assets" + CurrentVersion
	// 	os.RemoveAll(__path)
	// 	err = cp.Copy(path+"/assets", __path)
	// }

	// if err != nil {
	// 	return err
	// }

	// err = os.RemoveAll(path)
	// return err
	return nil
}

// Get uptodowns HTML page
func CurlAPKLink(link string) (*html.Node, error) {
	resp, err := Client.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, err
}

// Download the APK from uptodown
func WgetAPK(game *GameLink, downloadUrl, version, fp string) (string, error) {
	resp, err := Client.Get(downloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if fp == "" {
		fp = strings.ToLower(strings.ReplaceAll(game.Name, " ", "")) + "-" + version + ".apk" // clashofclans-14.426.4.apk
	}

	__fd, err := os.Create(fp)
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	_, err = __fd.Write(bytes)
	__fd.Close()
	return fp, err
}
