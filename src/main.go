package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"os21345678/search-engine-crawler/src/util"
)

const SITE_LIMIT = 1000000

type WebsiteData struct {
	URL         string
	Title       string
	Description string
	Keywords    []string
	Rank        int
}


func main() {
    records := readCsvFile("top-1m.csv")
    var wg sync.WaitGroup

    for lineNumber, line := range records {
        if lineNumber == SITE_LIMIT {
            break
        }

        url := "https://" + strings.TrimSpace(line[1])

        wg.Add(1)
        go func(ln int, u string) {
            defer wg.Done()

            crawler:= util.NewCrawler(url, ln)

            if err := crawler.Crawl(url); err != nil {
                fmt.Println("Failed to crawl", url, ":", err)
                return
            }
        }(lineNumber, url)
    }

    wg.Wait()
}

func readCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}
