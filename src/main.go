package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly/v2"
)

// func processTopMil(rc io.Reader) (ch chan []string) {
// 	ch = make(chan []string, 10)
// 	go func() {
// 		r := csv.NewReader(rc)
// 		if _, err := r.Read(); err != nil { //read header
// 			log.Fatal(err)
// 		}
// 		defer close(ch)
// 		for {
// 			rec, err := r.Read()
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 				log.Fatal(err)
// 			}
// 			ch <- rec
// 		}
// 	}()
// 	return
// }

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

	records, err := ReadCsv("/Users/oliver/Documents/Personal/Programming/search-engine-crawler/top-1m.csv")
	if err != nil {
		log.Fatal("Failed to read CSV: ", err)
	}
	for _, record := range records {
		url := "https://" + record[1] + "/"
		c.Visit(url)
		fmt.Println(url)
	}
	c.Wait()
}
