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
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/amaanq/sc-compression"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/withmandala/go-log"
	"golang.org/x/net/html"
)

var (
	CurrentVersion string
	Path           = flag.String("path", ".", "output directory for /assets")
	Client         = LoadRetryClient()
	Log            = log.New(os.Stdout).WithColor().WithDebug().WithTimestamp()
)

func LoadRetryClient() *retryablehttp.Client {
	Client := retryablehttp.NewClient()
	Client.Logger = nil
	Client.RetryMax = 5
	return Client
}

func FixPath() {
	if !strings.HasSuffix(*Path, "/") {
		*Path += "/"
	}
}

func UpdateAPK() error {
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

	Log.Info("Decompressing assets...")
	if _, err = WalkAndDecompressAssets(ClashofClans.ValidDirectories, "clash-"+version, "decompressed"+version); err != nil {
		Log.Error(err)
		return err
	}

	Log.Info("Done!")

	return nil
}

// Walk the assets folder and decompress each file inside
func WalkAndDecompressAssets(validDirs []string, fpToDecompiledAPK, fpToOutputFiles string) (string, error) {
	os.RemoveAll(fpToOutputFiles)
	err := os.Mkdir(fpToOutputFiles, 0755)
	if err != nil && !os.IsExist(err) {
		Log.Error(fpToOutputFiles)
		Log.Error(err)
		return "", err
	}

	for _, subdir := range validDirs {
		entries, err := os.ReadDir("./" + fpToDecompiledAPK + "/assets/" + subdir + "/")
		if err != nil {
			continue
		}

		err = os.Mkdir(fpToOutputFiles+"/"+subdir, 0755)
		if err != nil && !os.IsExist(err) {
			Log.Error(err)
			return "", err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			fileName := entry.Name()
			fullPath := "./" + fpToDecompiledAPK + "/assets/" + subdir + "/" + fileName
			t := time.Now()
			year, month, day := t.Date()
			hour, min, sec := t.Clock()
			date := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
			fmt.Printf("\033[2K\r\033[0;32m[INFO] \033[0;34m %s \033[0mDecompressing %s", date, fullPath)
			compFile, err := os.Open(fullPath)
			if err != nil {
				return "", err
			}
			decompressor := ScCompression.NewDecompressor(compFile)
			reader, err := decompressor.Decompress()
			if err != nil {
				Log.Errorf("Failed to decompress %s: %s\n", fullPath, err)
				continue
			}

			fd, err := os.Create(fpToOutputFiles + "/" + subdir + "/" + fileName)
			if err != nil {
				return "", err
			}
			if _, err = io.Copy(fd, reader); err != nil {
				return "", err
			}
			if err = fd.Close(); err != nil {
				return "", err
			}
		}
	}
	fmt.Printf("\n")
	return fpToOutputFiles, nil
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
	return err
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

func CleanUp(assetsFP, apkFP string) error {
	outerFP := strings.Split(assetsFP, "/")
	outFP := outerFP[len(outerFP)-1]
	err := os.Rename(assetsFP, "./"+outFP)
	if err != nil {
		return err
	}
	err = os.RemoveAll(outerFP[0])
	if err != nil {
		return err
	}
	err = os.RemoveAll(apkFP)
	if err != nil {
		return err
	}
	return nil
}
