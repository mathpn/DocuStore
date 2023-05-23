package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

type RuntimeState struct {
	dataFolder string
	db         *sql.DB
	index      *BinaryTree
	docCounter *DocCounter
}

func NewRuntimeState() (*RuntimeState, error) {
	gob.Register(DocSummary{})
	stateDir := xdg.StateHome
	dataFolder := filepath.Join(stateDir, "DocuStore")
	err := os.MkdirAll(dataFolder, 0755)
	if err != nil {
		return nil, err
	}

	db, err := NewDBConnection(filepath.Join(dataFolder, "storage.db"))
	if err != nil {
		return nil, err
	}

	state := &RuntimeState{
		dataFolder,
		db,
		nil,
		nil,
	}
	return state, nil
}

// Load or create BTree index
func (s *RuntimeState) loadIndex() error {
	gob.Register(BinaryTree{})
	indexPath := filepath.Join(s.dataFolder, "index.gob")

	index := &BinaryTree{nil, 0}
	latestTs, err := GetLatestTimestamp(s.db)
	if err != nil {
		return err
	}
	if latestTs == 0 {
		// no documents
		s.index = index
		return nil
	}

	err = LoadStruct(indexPath, &index)
	if err != nil {
		fmt.Println("Error reading BTree index:", err, "- attempting to recover")
		index, err = s.recoverIndex()
		if err != nil {
			fmt.Println("Index recovery failed, documents may have been lost:", err)
			return err
		}
		fmt.Println("BTree index succesfully recovered")
	}

	if index.Timestamp != latestTs {
		fmt.Printf("%+v - %+v\n", index.Timestamp, latestTs)
		fmt.Println("BTree index is out of sync with latest changes, recovering")
		index, err = s.recoverIndex()
		if err != nil {
			fmt.Println("BTree index recovery failed, documents may have been lost:", err)
			return err
		}
	}

	s.index = index
	return nil
}

func (s *RuntimeState) recoverIndex() (*BinaryTree, error) {
	docIDs, err := ListDocuments(s.db)
	if err != nil {
		return nil, err
	}
	btree := &BinaryTree{nil, 0}
	for _, docID := range docIDs {
		doc, ts, err := LoadDocSummary(s.db, docID)
		if err != nil {
			return nil, err
		}
		btree.InsertDoc(doc, ts)
	}
	indexPath := filepath.Join(s.dataFolder, "index.gob")
	err = SaveStruct(indexPath, btree)
	if err != nil {
		return nil, err
	}
	return btree, nil
}

// Load or create DocCounter
func (s *RuntimeState) loadCounter() error {
	gob.Register(DocCounter{})
	dcPath := filepath.Join(s.dataFolder, "docCounter.gob")

	docCounter := NewDocCounter()
	latestTs, err := GetLatestTimestamp(s.db)
	if err != nil {
		return err
	}
	if latestTs == 0 {
		// no documents
		s.docCounter = docCounter
		return nil
	}

	err = LoadStruct(dcPath, docCounter)
	if err != nil {
		fmt.Println("Error reading docCounter:", err, "- attempting to recover")
		docCounter, err = s.recoverDocCounter()
		if err != nil {
			fmt.Println("DocCounter recovery failed, documents may have been lost:", err)
			return err
		}
		fmt.Println("docCounter succesfully recovered")
	}

	if docCounter.Timestamp != latestTs {
		fmt.Println("DocCounter is out of sync with latest changes, recovering")
		docCounter, err = s.recoverDocCounter()
		if err != nil {
			fmt.Println("DocCounter recovery failed, documents may have been lost:", err)
			return err
		}
	}

	s.docCounter = docCounter
	return nil
}

func (s *RuntimeState) recoverDocCounter() (*DocCounter, error) {
	docIDs, err := ListDocuments(s.db)
	if err != nil {
		return nil, err
	}

	docCounter := NewDocCounter()
	for _, docID := range docIDs {
		doc, ts, err := LoadDocSummary(s.db, docID)
		if err != nil {
			return nil, err
		}
		docCounter.AddDocument(doc, ts)
	}

	dcPath := filepath.Join(s.dataFolder, "docCounter.gob")
	err = SaveStruct(dcPath, docCounter)
	if err != nil {
		return nil, err
	}
	return docCounter, nil
}

func addFile(filePath string, state *RuntimeState) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	text := string(content)
	addDocument(text, text, filePath, DocType(Text), state)
	return nil
}

func addText(text string, title string, state *RuntimeState) error {
	err := addDocument(text, text, title, DocType(Text), state)
	return err

}

func addURL(url string, state *RuntimeState) error {
	title, text := ScrapeText(url)
	err := addDocument(text, url, title, DocType(URL), state)
	return err
}

func addDocument(text string, identifier string, title string, docType DocType, state *RuntimeState) error {
	if title == "" {
		return errors.New("empty title is not allowed")
	}
	if text == "" {
		return errors.New("empty content")
	}
	ts := time.Now().Unix()
	docSummary := NewDocSummary(text, identifier, title, docType)
	rows, err := InsertDocument(state.db, docSummary, text, ts)
	if err != nil {
		return err
	}
	if rows == 0 {
		fmt.Println("Document is already in the collection")
		return nil
	}

	state.index.InsertDoc(docSummary, ts)
	state.docCounter.AddDocument(docSummary, ts)

	err = SaveStruct(
		filepath.Join(state.dataFolder, "index.gob"),
		state.index,
	)
	if err != nil {
		return err
	}

	err = SaveStruct(
		filepath.Join(state.dataFolder, "docCounter.gob"),
		state.docCounter,
	)
	return err
}

func queryDocument(text string, state *RuntimeState) ([]*SearchResult, error) {
	tokens := Tokenize(text)
	docIDs := state.index.SearchTokens(tokens)
	docSummaries, err := LoadDocSummaries(context.Background(), state.db, docIDs...)
	if err != nil {
		return nil, err
	}

	similarities := TFIDFSimilarity(text, state.docCounter, docSummaries...)
	return similarities, nil
}

func printSearchResults(sims []*SearchResult) {
	fmt.Println("Here are the top 5 matches:")
	for i, sim := range sims {
		if sim.Score == 0.0 {
			break
		}
		if i == 5 {
			break
		}
		fmt.Printf("Match: %d | Score: %.2f\n", i+1, sim.Score)
		fmt.Println(sim.Title)
		fmt.Println("--------------------------------------------")
	}
}
