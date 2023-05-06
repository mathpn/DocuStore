package main

import (
	"fmt"
	"math"
	"os"
	"regexp"
	"sort"
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

type SimResult struct {
	docId uuid.UUID
	score float64
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
	docCounts  map[string]int     // number of documents with word
	idf        map[string]float64 // log of inverse document frequency
	nDocs      int                // number of documents
}

func NewTFIDF() *TFIDF {
	var docVectors []*DocVector
	return &TFIDF{
		docVectors,
		make(map[string]int),
		make(map[string]float64),
		0,
	}
}

func (tfidf *TFIDF) calculateIDF() {
	nDocs := float64(len(tfidf.docVectors))
	for token, count := range tfidf.docCounts {
		tfidf.idf[token] = math.Log(nDocs / float64(count))
	}
}

func (tfidf *TFIDF) AddDocuments(docVectors ...*DocVector) {
	for i := 0; i < len(docVectors); i++ {
		tfidf.nDocs++
		tfidf.docVectors = append(tfidf.docVectors, docVectors[i])
		for token := range docVectors[i].termFreqs {
			tfidf.docCounts[token]++
		}
	}
}

func (tfidf *TFIDF) Similarity(text string) []*SimResult {
	termFreqs := getTermFrequency(text)
	tfidf.calculateIDF()
	scores := make([]float64, tfidf.nDocs)
	refNorms := make([]float64, tfidf.nDocs)
	var queryNorm float64 = 0.0
	for token, queryCount := range termFreqs {
		for i, docVector := range tfidf.docVectors {
			refCount := docVector.termFreqs[token]
			scores[i] += queryCount * refCount
			refNorms[i] += refCount * refCount
			queryNorm += queryCount * queryCount
		}
	}
	result := make([]*SimResult, tfidf.nDocs)
	queryNorm = math.Sqrt(queryNorm)
	for i := 0; i < len(scores); i++ {
		result[i] = &SimResult{
			tfidf.docVectors[i].docId,
			scores[i] / (queryNorm * math.Sqrt(refNorms[i])),
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].score > result[j].score // descending order
	})
	return result
}

func main() {
	text, err := os.ReadFile("./text.txt")
	check(err)
	text2, err := os.ReadFile("./text2.txt")
	check(err)

	vector := NewDocVector(string(text))
	vector2 := NewDocVector(string(text2))

	tfidf := NewTFIDF()
	tfidf.AddDocuments(vector, vector2)
	scores := tfidf.Similarity("lorem ipsum")
	fmt.Printf("%+v\n", scores[0])
}
