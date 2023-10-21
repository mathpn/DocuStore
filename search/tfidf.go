package search

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
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
	TermFreqs  map[string]float64
	DocID      string
	Title      string
	Identifier string
	Type       DocType
	SquareNorm float64
}

type SearchResult struct {
	DocID      string
	Title      string
	Identifier string
	Type       DocType
	Score      float64
}

func NewDocSummary(text string, identifier string, title string, docType DocType) *DocSummary {
	termFreqs, norm := getTermFrequency(text)
	return &DocSummary{
		DocID:      hashDocument(identifier),
		Title:      title,
		Identifier: identifier,
		Type:       docType,
		TermFreqs:  termFreqs,
		SquareNorm: norm,
	}
}

func hashDocument(text string) string {
	hash := sha256.Sum256([]byte(text))
	hashString := hex.EncodeToString(hash[:])
	return hashString
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
	var squareNorm float64
	for token, count := range termCounts {
		freq := float64(count) / nTokens
		squareNorm += freq * freq
		termFreqs[token] = freq
	}
	return termFreqs, squareNorm
}

type DocCounter struct {
	DocCounts map[string]int
	idf       map[string]float64
	NumDocs   int
	Timestamp int64
}

func NewDocCounter() *DocCounter {
	return &DocCounter{
		NumDocs:   0,
		DocCounts: make(map[string]int),
		idf:       make(map[string]float64),
		Timestamp: 0,
	}
}

func (d *DocCounter) AddDocument(DocSummary *DocSummary, timestamp int64) {
	d.NumDocs++
	for token := range DocSummary.TermFreqs {
		d.DocCounts[token]++
	}
	d.Timestamp = timestamp
}

func (d *DocCounter) calculateIDF() {
	for token, count := range d.DocCounts {
		d.idf[token] = math.Log(float64(d.NumDocs) / float64(count))
	}
}

func TFIDFSimilarity(text string, c *DocCounter, docs ...*DocSummary) []*SearchResult {
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
	for i := 0; i < len(scores); i++ {
		invNorm := math.Sqrt(1 / (queryNorm*docs[i].SquareNorm + 1e-8))
		result[i] = &SearchResult{
			DocID:      docs[i].DocID,
			Title:      docs[i].Title,
			Type:       docs[i].Type,
			Identifier: docs[i].Identifier,
			Score:      math.Sqrt(scores[i] * invNorm),
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score // descending order
	})
	return result
}
