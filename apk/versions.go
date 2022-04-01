package apk

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type GameLink struct {
	URL              string
	Name             string
	ValidDirectories []string
}

type VersionData struct {
	Version     string
	URL         string
	Date        string
	DownloadURL string
}

var (
	ClashofClans = GameLink{URL: "https://clash-of-clans.en.uptodown.com/android/versions/%d", Name: "Clash of Clans", ValidDirectories: []string{"csv", "localization", "logic"}}
	ClashRoyale  = GameLink{URL: "https://clash-royale.en.uptodown.com/android/versions/%d", Name: "Clash Royale"}
	BrawlStars   = GameLink{URL: "https://brawl-stars.en.uptodown.com/android/versions/%d", Name: "Brawl Stars", ValidDirectories: []string{"csv_client", "csv_logic", "localization", "logic"}}
	ClashMini    = GameLink{URL: "https://clash-mini.en.uptodown.com/android/versions/%d", Name: "Clash Mini"}
	HayDay       = GameLink{URL: "https://hay-day.en.uptodown.com/android/versions/%d", Name: "Hay Day", ValidDirectories: []string{"data", "localization"}}
	ClashQuest   = GameLink{URL: "https://clash-quest.en.uptodown.com/android/versions/%d", Name: "Clash Quest"}
	BoomBeach    = GameLink{URL: "https://boom-beach.en.uptodown.com/android/versions/%d", Name: "Boom Beach"}
	Everdale     = GameLink{URL: "https://everdale.en.uptodown.com/android/versions/%d", Name: "Everdale"}
	HayDayPop    = GameLink{URL: "https://hay-day-pop.en.uptodown.com/android/versions/%d", Name: "Hay Day Pop"}
	RushWars     = GameLink{URL: "https://rush-wars.en.uptodown.com/android/versions/%d", Name: "Rush Wars"}

	AllGameLinks = []GameLink{
		ClashofClans,
		ClashRoyale,
		BrawlStars,
		ClashMini,
		HayDay,
		ClashQuest,
		BoomBeach,
		Everdale,
		HayDayPop,
		RushWars,
	}

	ErrLastPage = fmt.Errorf("End of the Line!")
)

func GetAllVersions(gamelink string) ([]VersionData, error) {
	var page int = 0
	var wg sync.WaitGroup
	keepGoing := true
	versions := make([]VersionData, 0)
	for {
		if keepGoing {
			time.Sleep(time.Millisecond * 100)
			wg.Add(1)
			go func(page int) {
				defer wg.Done()
				vers, err := GetVersions(gamelink, page)
				if err != nil && err != ErrLastPage {
					Log.Error(err)
					keepGoing = false
				}
				if err == ErrLastPage {
					keepGoing = false
				}
				versions = append(versions, vers...)
			}(page)
			page++
		} else {
			break
		}
	}
	wg.Wait()
	return versions, nil
}

func GetVersions(gamelink string, page int) ([]VersionData, error) {
	resp, err := Client.Get(fmt.Sprintf(gamelink, page))
	if err != nil {
		Log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		Log.Error(err)
		return nil, err
	}

	n, err := html.Parse(strings.NewReader(string(bytes)))
	if err != nil {
		Log.Error(err)
		return nil, err
	}

	query := goquery.NewDocumentFromNode(n)

	var currentPage int
	query.Find("span.page-link.active").Each(func(i int, s *goquery.Selection) {
		currentPage, err = strconv.Atoi(s.Text())
		if err != nil {
			Log.Error(err)
			return
		}
	})
	if currentPage != page {
		return nil, ErrLastPage
	}

	vers := make([]VersionData, 0)
	query.Find("div").Each(func(i int, s *goquery.Selection) {
		if val, ok := s.Attr("data-url"); ok {

			ch := make(chan string)
			go func(val string) {
				url, err := GetDownloadURL(val)
				if err != nil {
					Log.Error(err)
					return
				}
				ch <- url
			}(val)

			GetDownloadURL(val)
			vers = append(vers, VersionData{
				Version:     strings.ReplaceAll(strings.TrimSpace(s.Contents().Not("span").Text()), "_", "."),
				URL:         val,
				Date:        s.Find("span").Text(),
				DownloadURL: <-ch,
			})
		}
	})
	return vers, nil
}
