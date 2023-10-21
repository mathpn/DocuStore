package search

import (
	"math"
	"sort"
)

func calculateIDF(counter *DocCounter) map[string]float64 {
	idf := make(map[string]float64)
	for token, count := range counter.DocCounts {
		idf[token] = math.Log(float64(counter.NumDocs)/(1+float64(count))) + 1
	}
	return idf
}

func TFIDFSimilarity(text string, c *DocCounter, docs ...*DocSummary) []*SearchResult {
	termFreqs, queryNorm := getTermFrequency(text)
	idf := calculateIDF(c)
	scores := make([]float64, len(docs))
	for i, docSummary := range docs {
		for token, queryCount := range termFreqs {
			refCount := docSummary.TermFreqs[token]
			factor, ok := idf[token]
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
