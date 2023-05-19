package main

import (
	"math"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

var asciiRegex = regexp.MustCompile(`[^a-zA-Z0-9\s]`)

var maxTokenLength int = 48

type DocType int

const (
	URL DocType = iota
	Text
)

func (t DocType) String() string {
	switch t {
	case URL:
		return "URL"
	case Text:
		return "Text"
	default:
		return "unknown"
	}
}

type DocSummary struct {
	DocID      string
	Title      string
	Type       DocType
	Identifier string
	TermFreqs  map[string]float64 // term frequency
	Norm       float64            // norm of TermFreqs vector
}

type SearchResult struct {
	DocID      string
	Title      string
	Type       DocType
	Identifier string
	Score      float64
	RawContent string
}

func NewDocSummary(text string, identifier string, title string, docType DocType) *DocSummary {
	termFreqs, norm := getTermFrequency(text)
	return &DocSummary{
		DocID:      hashDocument(identifier),
		Title:      title,
		Identifier: identifier,
		Type:       docType,
		TermFreqs:  termFreqs,
		Norm:       norm,
	}
}

func Tokenize(text string) []string {
	text = unidecode.Unidecode(text)
	text = strings.ToLower(text)
	text = asciiRegex.ReplaceAllString(text, "")
	tokens := strings.Fields(text)
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]
		if len(t) > maxTokenLength {
			t = t[:maxTokenLength]
			tokens[i] = t
		}
	}
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

type DocCounter struct {
	nDocs     int
	DocCounts map[string]int     // number of documents with word
	idf       map[string]float64 // log of inverse document frequency
}

func NewDocCounter() *DocCounter {
	return &DocCounter{
		0,
		make(map[string]int),
		make(map[string]float64),
	}
}

func (d *DocCounter) AddDocuments(DocSummaries ...*DocSummary) {
	for i := 0; i < len(DocSummaries); i++ {
		d.nDocs++
		// d.DocSummaries = append(d.DocSummaries, DocSummaries[i])
		for token := range DocSummaries[i].TermFreqs {
			d.DocCounts[token]++
		}
	}
}

func (d *DocCounter) calculateIDF() {
	for token, count := range d.DocCounts {
		d.idf[token] = math.Log(float64(d.nDocs) / float64(count))
	}
}

func TFIDFSimilarity(text string, rawFolder string, c *DocCounter, docs ...*DocSummary) []*SearchResult {
	termFreqs, queryNorm := getTermFrequency(text)
	c.calculateIDF()
	scores := make([]float64, len(docs))
	for i, docSummary := range docs {
		for token, queryCount := range termFreqs {
			refCount := docSummary.TermFreqs[token]
			scores[i] += queryCount * refCount
		}
	}
	result := make([]*SearchResult, len(docs))
	queryNorm = math.Sqrt(queryNorm)
	for i := 0; i < len(scores); i++ {
		result[i] = &SearchResult{
			docs[i].DocID,
			docs[i].Title,
			docs[i].Type,
			docs[i].Identifier,
			scores[i] / (queryNorm*docs[i].Norm + 1e-8),
			"",
		}
		if docs[i].Type == DocType(Text) {
			raw := LoadText(filepath.Join(rawFolder, docs[i].DocID+".txt"))
			result[i].RawContent = raw
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score // descending order
	})
	return result
}
