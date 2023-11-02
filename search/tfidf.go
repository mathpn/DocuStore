package search

import (
	"math"
	"sort"

	lru "github.com/hashicorp/golang-lru/v2"
)

type tfidfSearcher struct {
	counter *DocCounter
	idf     map[string]float64
	cache   *lru.Cache[string, float64]
	ts      int64
}

func NewTFIDFSearcher(c *DocCounter) (Searcher, error) {
	cache, err := lru.New[string, float64](1024)
	if err != nil {
		return nil, err
	}
	return &tfidfSearcher{
		counter: c,
		idf:     make(map[string]float64),
		cache:   cache,
	}, nil
}

func (s *tfidfSearcher) calculateIDF() {
	if s.ts != s.counter.Ts {
		for token, count := range s.counter.DocCounts {
			s.idf[token] = math.Log(float64(s.counter.NumDocs)/(1+float64(count))) + 1
		}
		s.ts = s.counter.Ts
	}
}

func (s *tfidfSearcher) computeNorm(doc *DocSummary) float64 {
	var norm float64
	for token, count := range doc.TermFreqs {
		factor, ok := s.idf[token]
		if !ok {
			factor = 1.0
		}
		norm += count * count * factor * factor
	}
	s.cache.Add(doc.DocID, norm)
	return norm
}

func (s *tfidfSearcher) Search(text string, docs ...*DocSummary) []*SearchResult {
	termFreqs := getTermFrequency(text)
	s.calculateIDF()
	scores := make([]float64, len(docs))
	var queryNorm float64
	docNorms := make([]float64, len(docs))

	var value float64
	for token, queryCount := range termFreqs {
		factor, ok := s.idf[token]
		if !ok {
			factor = 1.0
		}
		// NOTE multiply by factor squared to avoid a second access to idf
		value = queryCount * factor * factor
		termFreqs[token] = value
		queryNorm += queryCount * value
	}

	var norm, refCount float64
	var ok bool
	for i, doc := range docs {
		norm, ok = s.cache.Get(doc.DocID)
		if !ok {
			norm = s.computeNorm(doc)
		}
		docNorms[i] = norm
		for token, value := range termFreqs {
			refCount = doc.TermFreqs[token]
			scores[i] += value * refCount
		}
	}
	result := make([]*SearchResult, len(docs))
	for i := 0; i < len(scores); i++ {
		invNorm := 1 / math.Sqrt(queryNorm*docNorms[i]+1e-8)
		result[i] = &SearchResult{
			DocID:      docs[i].DocID,
			Title:      docs[i].Title,
			Type:       docs[i].Type.String(),
			Identifier: docs[i].Identifier,
			Score:      math.Sqrt(scores[i] * invNorm),
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score // descending order
	})
	return result
}
