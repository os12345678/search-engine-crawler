package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/gocolly/colly/v2"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

var (
	user = os.Getenv("DB_USER")
	host = os.Getenv("DB_HOST")
	dbname = os.Getenv("DB_NAME")
	password = os.Getenv("DB_PASSWORD")
	port = os.Getenv("DB_PORT")
)

var lemmatisedMap map[string]string

func init() {
	data, err := os.ReadFile("lemmatised.json")
	if err != nil {
		log.Fatal("Failed to read lemmatised.json: ", err)
	}
	err = json.Unmarshal(data, &lemmatisedMap)
	if err != nil {
		log.Fatal("Failed to unmarshal lemmatised.json: ", err)
	}
}

func ReadCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	return csv.NewReader(file).ReadAll()
}

func connectDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func extractMetaDescription(e *colly.HTMLElement) string {
	return e.ChildAttr("meta[name=description]", "content")
}

func tokenize(text string) []string {
    text = strings.ToLower(text)

    re := regexp.MustCompile(`[\W_]+`) // Matches non-word characters and underscores
    tokens := re.Split(text, -1)

    var filteredTokens []string
    for _, token := range tokens {
        if !isEmptyOrSpace(token) { 
            filteredTokens = append(filteredTokens, token)
        }
    }
    return filteredTokens
}

func isEmptyOrSpace(str string) bool {
    for _, r := range str {
        if !unicode.IsSpace(r) {
            return false
        }
    }
    return true
}

func extractKeywords(text string) map[string]int {
	words := tokenize(text)
	wordCounts := make(map[string]int)
	for _, word := range words {
		if lemma, found := lemmatisedMap[word]; found {
            wordCounts[lemma]++
		}
	}
	return wordCounts
}

func lemmatize(keywords map[string]int) []string {
	return []string{}
}

func storeData(db *pgx.Conn, title, description, text, url string, keywords []string) error {
    tx, err := db.Begin(context.Background()) // Start a transaction
    if err != nil {
        return err
    }
    defer tx.Rollback(context.Background())  // Rollback if an error occurs

    // 1. Insert into 'websites' table
    var websiteID uuid.UUID
    err = tx.QueryRow(context.Background(), `
        INSERT INTO websites (title, description, url, word_count, rank) 
        VALUES ($1, $2, $3, $4, $5) 
        RETURNING id`, 
        title, description, url, calculateWordCount(text), calculateRank(url),
    ).Scan(&websiteID)
    if err != nil {
        return err
    }

    // 2. Insert into 'keywords' and 'website_keywords'
    for _, word := range keywords {
        // Upsert logic for keywords
        var keywordID uuid.UUID
        err = tx.QueryRow(context.Background(), `
                INSERT INTO keywords (word, documents_containing_word) 
                VALUES ($1, 1)  
                ON CONFLICT (word) DO UPDATE SET documents_containing_word = keywords.documents_containing_word + 1
                RETURNING id`, word).Scan(&keywordID)
        if err != nil {
            return err
        }

        // Insert into website_keywords
        _, err = tx.Exec(context.Background(), `
                INSERT INTO website_keywords (keyword_id, website_id, occurrences, position)
                VALUES ($1, $2, $3, $4)`, 
                keywordID, websiteID, countKeywordOccurrences(text, word), calculatePosition(text, word),
        )
        if err != nil {
            return err
        }
    }

    err = tx.Commit(context.Background()) // Commit transaction
    return err
}




func main() {
	csvData, err := ReadCsv("top-1m.csv")
	if err != nil {
		log.Fatal("Failed to read CSV: ", err)
	}

	db, err := connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close(context.Background())

	c := colly.NewCollector(colly.Async(true))
	c.IgnoreRobotsTxt = false
	c.Limit(&colly.LimitRule{
		DomainGlob:   "*",
		Parallelism:  10,
		RandomDelay:  1 * time.Second,
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		title := e.ChildText("title")
		description := extractMetaDescription(e)
		url := e.Request.URL.String()
		text := e.Text

		// Keyword extraction and lemmatization
		keywords := extractKeywords(text)
		lemmatizedKeywords := lemmatize(keywords)

		// Database operations
		err := storeData(db, title, description, text, url, lemmatizedKeywords)
		if err != nil {
			log.Println("Error storing:", url, err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
	fmt.Println("Visiting:", r.URL)
	})

	// Start crawling from the URLs in the CSV
	for _, row := range csvData {
		c.Visit(row[1]) 
	}
	c.Wait()
}
