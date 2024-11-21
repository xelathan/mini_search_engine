package search

import (
	"fmt"
	"time"

	"github.com/xelathan/mini_search_engine/db"
)

func RunEngine() {
	fmt.Println("Running engine")
	defer fmt.Println("Engine stopped")

	settings := &db.SearchSetting{}
	err := settings.Get()

	if err != nil {
		fmt.Println("Error getting settings")
		return
	}

	if !settings.SearchOn {
		fmt.Println("Search is off")
		return
	}

	crawl := &db.CrawledUrl{}

	nextUrls, err := crawl.GetNextCrawlUrls(int(settings.Amount))
	if err != nil {
		fmt.Println("Error getting next urls")
		return
	}

	newUrls := []db.CrawledUrl{}
	testedTime := time.Now()

	for _, next := range nextUrls {
		result := runCrawl(next.Url)
		if !result.Success {
			err := next.UpdatedUrl(
				db.CrawledUrl{
					ID:              next.ID,
					Url:             next.Url,
					Success:         false,
					CrawlDuration:   result.CrawlData.CrawlTime,
					ResponseCode:    result.ResponseCode,
					PageTitle:       result.CrawlData.PageTitle,
					PageDescription: result.CrawlData.PageDescription,
					Heading:         result.CrawlData.Headings,
					LastTested:      &testedTime,
				},
			)

			if err != nil {
				fmt.Println("Error updating url")
			}
			continue
		}

		// Success

		err := next.UpdatedUrl(
			db.CrawledUrl{
				ID:              next.ID,
				Url:             next.Url,
				Success:         result.Success,
				CrawlDuration:   result.CrawlData.CrawlTime,
				ResponseCode:    result.ResponseCode,
				PageTitle:       result.CrawlData.PageTitle,
				PageDescription: result.CrawlData.PageDescription,
				Heading:         result.CrawlData.Headings,
				LastTested:      &testedTime,
			},
		)
		if err != nil {
			fmt.Println("Error updating url", next.Url)
		}

		for _, newUrl := range result.CrawlData.Links.External {
			newUrls = append(newUrls, db.CrawledUrl{Url: newUrl})
		}
	}

	if !settings.AddNewUrls {
		return
	}

	for _, newUrl := range newUrls {
		if err := newUrl.Save(); err != nil {
			fmt.Println("Error saving url")
		}
	}

	fmt.Println("Added new urls", len(newUrls))
}

func RunIndex() {
	fmt.Println("Running indexer")
	defer fmt.Println("Indexer finished")

	crawled := &db.CrawledUrl{}
	notIndexed, err := crawled.GetNotIndex()
	if err != nil {
		fmt.Println("Error getting not indexed urls")
		return
	}

	idx := make(Index)
	idx.Add(notIndexed)
	searchIndex := db.SearchIndex{}
	if err := searchIndex.Save(idx, notIndexed); err != nil {
		fmt.Println("Error saving index")
		return
	}

	if err := crawled.SetIndexedTrue(notIndexed); err != nil {
		fmt.Println("Error setting indexed true")
		return
	}
}
