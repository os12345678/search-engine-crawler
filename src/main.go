package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"os21345678/search-engine-crawler/src/util"

	"github.com/gocolly/colly/v2"
	"github.com/oklog/ulid"
)

type Website struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	WordCount   int    `json:"wordCount"`
	Rank        int    `json:"rank"`
}

type Keyword struct {
	ID                      string `json:"id"`
	Word                    string `json:"word"`
	DocumentsContainingWord int    `json:"documentsContainingWord"`
}

type WebsiteKeyword struct {
	KeywordID   string `json:"keywordId"`
	WebsiteID   string `json:"websiteId"`
	Occurrences int    `json:"occurrences"`
	Position    int    `json:"position"`
}

type Database interface {
	Insert(website *Website) error
	InsertManyOrUpdate(keywords []*Keyword) error
	InsertMany(websiteKeywords []*WebsiteKeyword) error
}
type Crawler struct {
	collector *colly.Collector
	rank      int
	db        Database
}

func NewCrawler(rank int, db Database) *Crawler {
	c := colly.NewCollector()
	c.IgnoreRobotsTxt=false

	// Disallow crawling of unwanted resource types
	c.OnRequest(func(r *colly.Request) {
		switch r.Ctx.Get("resource_type") {
		case "image", "media", "stylesheet", "font", "script":
			r.Abort()
		}
	})

	return &Crawler{
		collector: c,
		rank:      rank,
		db:        db,
	}
}

func (c *Crawler) crawl(url string) {
	fmt.Println("Crawling:", url)

	c.collector.OnHTML("html", func(e *colly.HTMLElement) {
		lang := e.Attr("lang")
		if lang != "en" && lang != "en-gb" {
			return // Only index English websites 
		}

		title := e.ChildText("title")
		desc := e.ChildAttr("meta[name=description]", "content")

		// Extract and lemmatize words (you'll need lemmatization logic)
		text := e.Text
		words := util.Lemmatize(text)

		// Construct data objects
		website := &Website{
			ID:          ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String(),
			Title:       title,
			Description: desc,
			URL:         url,
			WordCount:   len(words),
			Rank:        c.rank,
		}

		keywordIDs := make(map[string]string)
		wordIndices := make(map[string]int)
		wordPositions := []int{}
		wordIDs := []string{}

		position := 1
		for _, word := range words {
			wordIndices[word]++

			if _, ok := keywordIDs[word]; !ok {
				keywordIDs[word] = ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
			}
			wordIDs = append(wordIDs, keywordIDs[word])
			wordPositions = append(wordPositions, position)
			position++
		}
	})
}