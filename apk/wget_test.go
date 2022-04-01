package apk

import (
	"fmt"
	"testing"
)

func TestWget(t *testing.T) {
	// fmt.Println("begin")
	// vers, err := GetAllVersions(ClashofClans.URL)
	// if err != nil {
	// 	Log.Errorf("Versions() error = %v", err)
	// 	return
	// }
	// // log first url
	// Log.Info(vers[0].DownloadURL)
	// resp, err := Client.Get(vers[0].DownloadURL)
	// if err != nil {
	// 	Log.Error(err)
	// 	return
	// }
	// defer resp.Body.Close()

	// total := int64(0)
	// speed := 0.0
	// start := time.Now()
	// pr := &WgetReader{resp.Body, resp.ContentLength, func(r int64) {
	// 	total += r
	// 	percent := float64(total) / float64(resp.ContentLength) * 100
	// 	t := time.Now()
	// 	year, month, day := t.Date()
	// 	hour, min, sec := t.Clock()
	// 	date := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
	// 	speed = float64(total)/float64(time.Since(start).Milliseconds())/125/8
	// 	if r > 0 && percent != 100.00 {
	// 		fmt.Printf("\033[2K\r\033[0;32m[INFO] \033[0;34m %s \033[0m%.2f%% %.2f mb/s", date, percent, speed)
	// 	} else {
	// 		fmt.Printf("\033[2K\r\033[0;32m[INFO] \033[0;34m %s \033[0m100%% %.2f mb/s took %.2f seconds", date, speed, time.Since(start).Seconds())
	// 	}
	// }}
	// io.Copy(io.Discard, pr)
	// fmt.Printf("\n")
	fmt.Println("begin")
	vers, err := GetAllVersions(ClashofClans.URL)
	if err != nil {
		Log.Errorf("Versions() error = %v", err)
		return
	}
	// log first url
	Log.Info(vers[0].DownloadURL)
	WgetAPK(&ClashofClans, vers[0].DownloadURL, "", "test.apk")
}
