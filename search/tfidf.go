package search

import (
	"math"
	"sort"
)

func TFIDFSimilarity(text string, c *DocCounter, docs ...*DocSummary) []*SearchResult {
	termFreqs, queryNorm := getTermFrequency(text)
	c.calculateIDF() // FIXME move here and cache
	scores := make([]float64, len(docs))
	for i, docSummary := range docs {
		for token, queryCount := range termFreqs {
			refCount := docSummary.TermFreqs[token]
			scores[i] += queryCount * refCount * c.idf[token]
			scores[i] *= c.idf[token]
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
