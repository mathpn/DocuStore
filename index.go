package main

import "DocuStore/search"

type docToken struct {
	docID string
	token string
}

type HashmapIndex struct {
	Map       map[string][]string
	Timestamp int64 // timestamp of latest change
}

func (t *HashmapIndex) InsertDoc(doc *search.DocSummary, timestamp int64) {
	for token := range doc.TermFreqs {
		docToken := &docToken{doc.DocID, token}
		t.insert(docToken)
	}
	t.Timestamp = timestamp
}

func (t *HashmapIndex) insert(Data *docToken) {
	if t.Map == nil {
		t.Map = make(map[string][]string, 0)
	} else {
		val, ok := t.Map[Data.token]
		if ok {
			t.Map[Data.token] = append(val, Data.docID)
		} else {
			valArray := []string{Data.docID}
			t.Map[Data.token] = valArray
		}
	}
}

func (t *HashmapIndex) SearchTokens(tokens []string) []string {
	docMap := make(map[string]bool)
	out := make([]string, 0)
	for i := 0; i < len(tokens); i++ {
		docIDs, ok := t.Map[tokens[i]]
		if ok {
			for _, docID := range docIDs {
				if _, val := docMap[docID]; !val {
					docMap[docID] = true
					out = append(out, docID)
				}
			}
		}
	}
	return out
}
