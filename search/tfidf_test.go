package search

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func loadWords() []string {
	file, err := os.Open("../assets/words.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	words := make([]string, 0)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words
}

func prepareBench(nDocs int, nWords int, words []string) (Searcher, []*DocSummary) {
	counts := make(map[string]int, 0)
	docs := make([]*DocSummary, 0)
	for _, w := range words {
		counts[w] = rand.Intn(nDocs)
	}
	var text string
	for i := 0; i < nDocs; i++ {
		text = ""
		for j := 0; j < nWords; j++ {
			id := rand.Intn(len(words))
			text += " " + words[id]
		}
		title := fmt.Sprintf("doc %d", i)
		doc := NewDocSummary(text, title, title, Text)
		docs = append(docs, doc)
	}
	counter := &DocCounter{
		DocCounts: counts,
		NumDocs:   100,
		Ts:        10000000,
	}
	searcher, err := NewTFIDFSearcher(counter)
	if err != nil {
		panic(err)
	}
	return searcher, docs
}

func BenchmarkTfidf(b *testing.B) {
	words := loadWords()
	var query string
	for _, nDocs := range []int{100, 1000, 10000} {
		for _, nWords := range []int{10, 100, 1000} {
			searcher, docs := prepareBench(nDocs, nWords, words)
			for _, lenQuery := range []int{1, 10, 100} {
				query = ""
				ids := rand.Perm(len(words))[:lenQuery]
				for _, id := range ids {
					query += " " + words[id]
				}
				b.Run(fmt.Sprintf("%d query words %d docs with %d words", lenQuery, nDocs, nWords), func(_ *testing.B) {
					_ = searcher.Search(query, docs...)
				})
			}
		}
	}
}
