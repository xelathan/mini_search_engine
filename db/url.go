package db

import (
	"time"

	"gorm.io/gorm"
)

type CrawledUrl struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Url             string         `json:"url" gorm:"unique;not null"`
	Success         bool           `json:"success" gorm:"default:false"`
	CrawlDuration   time.Duration  `json:"crawlDuration"`
	ResponseCode    int            `json:"responseCode"`
	PageTitle       string         `json:"pageTitle"`
	PageDescription string         `json:"pageDescription"`
	Heading         string         `json:"heading"`
	LastTested      *time.Time     `json:"lastTested"`
	Indexed         bool           `json:"indexed" gorm:"default:false"`
	CreatedAt       *time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (crawled *CrawledUrl) UpdatedUrl(input CrawledUrl) error {
	tx := DBConn.Select("url", "success", "crawl_duration", "response_code", "page_title", "page_description", "heading", "last_tested", "indexed", "updated_at").
		Omit("created_at").Save(&input)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (crawled *CrawledUrl) GetNextCrawlUrls(limit int) ([]CrawledUrl, error) {
	var crawledUrls []CrawledUrl
	tx := DBConn.Where("last_tested IS NULL").Limit(limit).Find(&crawledUrls)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return crawledUrls, nil
}

func (crawled *CrawledUrl) Save() error {
	tx := DBConn.Save(&crawled)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (crawled *CrawledUrl) GetNotIndex() ([]CrawledUrl, error) {
	var urls []CrawledUrl
	tx := DBConn.Where("indexed = ? AND last_tested IS NOT NULL", false).Find(&urls)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return urls, nil
}

func (crawled *CrawledUrl) SetIndexedTrue(urls []CrawledUrl) error {
	for _, url := range urls {
		url.Indexed = true
		tx := DBConn.Save(&url)
		if tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
