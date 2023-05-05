package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/mozillazg/go-unidecode"
)

var asciiRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]`)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type DocVector struct {
	docId     uuid.UUID
	termFreqs map[string]float64 // term frequency
}

func NewDocVector(text string) *DocVector {
	termFreqs := getTermFrequency(text)
	return &DocVector{
		docId:     uuid.New(),
		termFreqs: termFreqs,
	}
}

func tokenize(text string) []string {
	text = asciiRegex.ReplaceAllString(text, "")
	text = strings.ToLower(text)
	text = unidecode.Unidecode(text)
	tokens := strings.Fields(text)
	return tokens
}

func getTermFrequency(text string) map[string]float64 {
	tokens := tokenize(text)
	termCounts := make(map[string]int)
	nTokens := float64(len(tokens))
	for _, token := range tokens {
		termCounts[token]++
	}
	termFreqs := make(map[string]float64)
	for token, count := range termCounts {
		termFreqs[token] = float64(count) / nTokens
	}
	return termFreqs
}

type TFIDF struct {
	docVectors []*DocVector
	idf        map[string]float64 // log of inverse document frequency
}

func NewTFIDF(docVectors ...*DocVector) *TFIDF {
	docFreqs := make(map[string]int)
	for _, docVector := range docVectors {
		for token := range docVector.termFreqs {
			docFreqs[token]++
		}
	}
	idf := make(map[string]float64)
	nDocs := float64(len(docVectors))
	for token, count := range docFreqs {
		idf[token] = math.Log(nDocs / float64(count))
	}

	return &TFIDF{
		docVectors,
		idf,
	}
}

func main() {
	text, err := os.ReadFile("./text.txt")
	check(err)
	text2, err := os.ReadFile("./text2.txt")
	check(err)

	vector := NewDocVector(string(text))
	vector2 := NewDocVector(string(text2))
	fmt.Printf("%+v\n", vector)
	fmt.Printf("%+v\n", vector2)

	tfidf := NewTFIDF(vector, vector2)
	fmt.Printf("%+v\n", tfidf)
}
