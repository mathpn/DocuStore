package search

import (
	"math"
	"sort"
)

type tfidfSearcher struct {
	counter *DocCounter
	idf     map[string]float64
	ts      int64
}

func NewTFIDFSearcher(c *DocCounter) Searcher {
	return &tfidfSearcher{
		counter: c,
		idf:     make(map[string]float64),
	}
}

func (s *tfidfSearcher) calculateIDF() {
	if s.ts != s.counter.Ts {
		for token, count := range s.counter.DocCounts {
			s.idf[token] = math.Log(float64(s.counter.NumDocs)/(1+float64(count))) + 1
		}
		s.ts = s.counter.Ts
	}
}

func (s *tfidfSearcher) Search(text string, docs ...*DocSummary) []*SearchResult {
	termFreqs, queryNorm := getTermFrequency(text)
	s.calculateIDF()
	scores := make([]float64, len(docs))
	for i, docSummary := range docs {
		for token, queryCount := range termFreqs {
			refCount := docSummary.TermFreqs[token]
			factor, ok := s.idf[token]
			if !ok {
				factor = 1.0
			}
			scores[i] += queryCount * refCount * factor
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
