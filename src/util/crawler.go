package util

import (
	"fmt"
	"net/url"

	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	collector *colly.Collector
}

func NewCrawler(url string, rank int) *Crawler {
	c := colly.NewCollector()
	c.IgnoreRobotsTxt=false

	c.OnRequest(func(r *colly.Request) {
		switch r.Ctx.Get("resource_type") {
		case "image", "media", "stylesheet", "font", "script":
			r.Abort()
		}
	})
	return &Crawler {
		collector: c,
	}
}

func (c *Crawler) Crawl(rawurl string) error {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		fmt.Println("Failed to parse URL:", err)
		return err
	}

	fmt.Println("Crawling:", parsedURL.String())

	c.collector.OnHTML("html", func(e *colly.HTMLElement) {
		lang := e.Attr("lang")
		if lang != "en" && lang != "en-gb" {
			return // Only index English websites 
		}
	})
	return nil
}








	// fmt.Println("Crawling:", url)

	// c.collector.OnHTML("html", func(e *colly.HTMLElement) {
	// 	lang := e.Attr("lang")
	// 	if lang != "en" && lang != "en-gb" {
	// 		return // Only index English websites 
	// 	}

	// 	title := e.ChildText("title")
	// 	desc := e.ChildAttr("meta[name=description]", "content")

	// 	text := e.Text
	// 	words := lemmatize(text)

	// 	keywordIDs := make(map[string]string)
	// 	wordIndices := make(map[string]int)
	// 	wordPositions := []int{}
	// 	wordIDs := []string{}

	// 	position := 1
	// 	for _, word := range words {
	// 		wordIndices[word]++

	// 		if _, ok := keywordIDs[word]; !ok {
	// 			keywordIDs[word] = ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader).String()
	// 		}
	// 		wordIDs = append(wordIDs, keywordIDs[word])
	// 		wordPositions = append(wordPositions, position)
	// 		position++
	// 	}
	// 	websiteId := uuid.New().String()

	// 	websiteIdsBatch := make([]string, len(words))
	// 	ranksBatch := make([]string, len(words))
	// 	wordIndicesBatch := make([]int, len(words))

	// 	// Populate batches - more optimized than map
	// 	for i := range websiteIdsBatch {
	// 		websiteIdsBatch[i] = websiteId
	// 		ranksBatch[i] = fmt.Sprintf("%d", c.rank)
	// 		wordIndicesBatch[i] = wordIndices[words[i]] // Assuming words[i] is the key
	// 	}

	// 	// Database Insertion - websites
	// 	err := c.Insert("website", map[string]interface{}{
	// 		"id":          websiteId,
	// 		"title":       title,
	// 		"description": desc,
	// 		"url":         url,
	// 		"word_count":  len(words),
	// 		"rank":        this.rank,
	// 	})

	// 	if err != nil {
	// 		fmt.Printf("[WARNING]: Failed to index: %s\n\n%v\n", url, err)
	// 		return
	// 	}

	// 	// Database Insertion - Keywords (with updates)
	// 	keywords := make([]*Keyword, 0, len(keywordIDs))
	// 	for word, id := range keywordIDs {
	// 		keywords = append(keywords, &Keyword{
	// 			ID:                      id,
	// 			Word:                    word,
	// 			DocumentsContainingWord: 1, // fix
	// 		})
	// 	}

	// 	//  _, err := c.db.InsertManyOrUpdate("keyword", []string{"id", "word", "documents_containing_word"}, [][]interface{}{}, []string{"id"}, "documents_containing_word = keyword.documents_containing_word + 1", []string{"id"})

	// 	// if err != nil {
	// 	// 	log.Printf("[WARNING]: Error updating keywords: %v\n", err)
	// 	// }

	// 	// Construct word IDs
	// 	updatedWordIds := make([]string, 0, len(keywords))
	// 	for _, keyword := range keywords {
	// 		updatedWordIds = append(updatedWordIds, keyword.ID)
	// 	}

	// 	// Database Insertion - website_keywords
	// 	websiteKeywords := make([]*WebsiteKeyword, 0, len(updatedWordIds)) 
	// 	for i := range updatedWordIds {
	// 		websiteKeywords = append(websiteKeywords, &WebsiteKeyword{
	// 			KeywordID:   updatedWordIds[i], 
	// 			WebsiteID:   websiteIdsBatch[i],
	// 			Occurrences: wordIndicesBatch[i],
	// 			Position:    wordPositions[i],
	// 		})
	// 	}

	// 	err = c.db.InsertMany("website_keyword", []string{"keyword_id", "website_id", "occurrences", "position"}, [][]interface{}{},)


	// 	if err != nil {
	// 		log.Printf("[WARNING]: Error inserting website keywords: %v\n", err)
	// 	}
		
	// 	fmt.Println("Successfully crawled:", url)
	// })

	// }
