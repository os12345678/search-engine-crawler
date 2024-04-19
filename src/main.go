package main

import (
	"encoding/csv"
	// "fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

func ReadCsv(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer file.Close()

	return csv.NewReader(file).ReadAll()
}

func main() {
	c := colly.NewCollector(colly.Async(true))

	c.IgnoreRobotsTxt = false

	c.Limit(&colly.LimitRule{
		DomainGlob:   "*",
		Parallelism:  10,
		RandomDelay:  1 * time.Second,
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	records, err := ReadCsv("top-1m.csv")
	if err != nil {
		log.Fatal("Failed to read CSV: ", err)
	}
	for _, record := range records {
		url := "https://" + record[1] + "/"
		c.Visit(url)
	}
	c.Wait()
}
