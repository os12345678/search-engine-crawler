package main

import (
	"log"
)

func loadLemmatisedMap() map[string]string {
	lemmatisedMap := make(map[string]string)
	records, err := ReadCsv("lemmatised.csv")
	if err != nil {
		log.Fatal("Failed to read CSV: ", err)
	}
	for _, record := range records {
		lemmatisedMap[record[0]] = record[1]
	}
	return lemmatisedMap
}