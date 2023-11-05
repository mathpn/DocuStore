package search

import (
	"crypto/sha256"
	"encoding/hex"
	"regexp"
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
}

type SearchResult struct {
	DocID      string
	Title      string
	Identifier string
	Type       string
	Score      float64
}

func NewDocSummary(text string, identifier string, title string, docType DocType) *DocSummary {
	termFreqs := getTermFrequency(text)
	return &DocSummary{
		DocID:      hashDocument(identifier),
		Title:      title,
		Identifier: identifier,
		Type:       docType,
		TermFreqs:  termFreqs,
	}
}

type Searcher interface {
	Search(text string, docs ...*DocSummary) []*SearchResult
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
	var t string
	for i := 0; i < len(tokens); i++ {
		t = tokens[i]
		if len(t) > maxTokenLength {
			t = t[:maxTokenLength]
			tokens[i] = t
		}
	}
	return tokens
}

func getTermFrequency(text string) map[string]float64 {
	tokens := Tokenize(text)
	termCounts := make(map[string]int)
	nTokens := float64(len(tokens))
	for _, token := range tokens {
		termCounts[token]++
	}
	termFreqs := make(map[string]float64, len(termCounts))
	for token, count := range termCounts {
		termFreqs[token] = float64(count) / nTokens
	}
	return termFreqs
}

type DocCounter struct {
	DocCounts map[string]int
	NumDocs   int
	Ts        int64
}

func NewDocCounter() *DocCounter {
	return &DocCounter{
		NumDocs:   0,
		DocCounts: make(map[string]int),
		Ts:        0,
	}
}

func (d *DocCounter) AddDocument(DocSummary *DocSummary, timestamp int64) {
	d.NumDocs++
	for token := range DocSummary.TermFreqs {
		d.DocCounts[token]++
	}
	if timestamp > d.Ts {
		d.Ts = timestamp
	}
}
