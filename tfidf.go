package main

import (
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var asciiRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]`)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type DocSummary struct {
	DocID     string
	Title     string
	TermFreqs map[string]float64 // term frequency
	Norm      float64            // norm of TermFreqs vector
}

type SimResult struct {
	DocID string
	Title string
	Score float64
}

func NewDocSummary(text string, identifier string, title string) *DocSummary {
	termFreqs, norm := getTermFrequency(text)
	return &DocSummary{
		DocID:     hashDocument(identifier),
		Title:     title,
		TermFreqs: termFreqs,
		Norm:      norm,
	}
}

func Tokenize(text string) []string {
	text = asciiRegex.ReplaceAllString(text, "")
	text = strings.ToLower(text)
	text = unidecode.Unidecode(text)
	tokens := strings.Fields(text)
	return tokens
}

func getTermFrequency(text string) (map[string]float64, float64) {
	tokens := Tokenize(text)
	termCounts := make(map[string]int)
	nTokens := float64(len(tokens))
	for _, token := range tokens {
		termCounts[token]++
	}
	termFreqs := make(map[string]float64)
	norm := 0.0
	for token, count := range termCounts {
		freq := float64(count) / nTokens
		norm += freq * freq
		termFreqs[token] = freq
	}
	norm = math.Sqrt(norm)
	return termFreqs, norm
}

type TFIDF struct {
	DocSummaries []*DocSummary
	docCounts    map[string]int     // number of documents with word
	idf          map[string]float64 // log of inverse document frequency
	nDocs        int                // number of documents
}

func NewTFIDF() *TFIDF {
	var DocSummaries []*DocSummary
	return &TFIDF{
		DocSummaries,
		make(map[string]int),
		make(map[string]float64),
		0,
	}
}

func (tfidf *TFIDF) calculateIDF() {
	nDocs := float64(len(tfidf.DocSummaries))
	for token, count := range tfidf.docCounts {
		tfidf.idf[token] = math.Log(nDocs / float64(count))
	}
}

func (tfidf *TFIDF) AddDocuments(DocSummaries ...*DocSummary) {
	for i := 0; i < len(DocSummaries); i++ {
		tfidf.nDocs++
		tfidf.DocSummaries = append(tfidf.DocSummaries, DocSummaries[i])
		for token := range DocSummaries[i].TermFreqs {
			tfidf.docCounts[token]++
		}
	}
}

func (tfidf *TFIDF) Similarity(text string) []*SimResult {
	termFreqs, queryNorm := getTermFrequency(text)
	tfidf.calculateIDF()
	scores := make([]float64, tfidf.nDocs)
	for token, queryCount := range termFreqs {
		for i, DocSummary := range tfidf.DocSummaries {
			refCount := DocSummary.TermFreqs[token]
			scores[i] += queryCount * refCount
		}
	}
	result := make([]*SimResult, tfidf.nDocs)
	queryNorm = math.Sqrt(queryNorm)
	for i := 0; i < len(scores); i++ {
		result[i] = &SimResult{
			tfidf.DocSummaries[i].DocID,
			tfidf.DocSummaries[i].Title,
			scores[i] / (queryNorm*tfidf.DocSummaries[i].Norm + 1e-8),
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score // descending order
	})
	return result
}
