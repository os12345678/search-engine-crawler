package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func loadLemmatizedWords() map[string]string {
    data, err := os.ReadFile("./lemmatizedMap.json")
    if err != nil {
        fmt.Println("Error reading lemmatizedMap.json:", err)
    }
	LemmatizedWords := make(map[string]string)
	err = json.Unmarshal(data, &LemmatizedWords)
	if err != nil {
		fmt.Println("Error unmarshalling lemmatizedMap.json:", err)
	}
	return LemmatizedWords
}

func lemmatize(text string) []string {
	
	lemmatizedWords := loadLemmatizedWords()
    // Pre-allocate a reasonable capacity based on an estimated average word count
    words := make([]string, 0, len(text)/5) 

    // Lowercase and filter punctuation
    var sb strings.Builder
    for _, r := range text {
        if unicode.IsLetter(r) || unicode.IsNumber(r) {
            sb.WriteRune(unicode.ToLower(r))
        } else if unicode.IsSpace(r) {
            if sb.Len() > 0 {
                words = append(words, sb.String())
                sb.Reset()
            }
        }
    }
    // Handle the last word
    if sb.Len() > 0 {
        words = append(words, sb.String())
    }

    // Lemmatization
    for i, word := range words {
        if lemma, found := lemmatizedWords[word]; found {
            words[i] = lemma
        }
    }

    return words
}