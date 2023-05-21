package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type RuntimeState struct {
	dataFolder string
	rawFolder  string
	db         *sql.DB
	index      *BinaryTree
	docCounter *DocCounter
}

func NewRuntimeState() *RuntimeState {
	gob.Register(DocSummary{})
	workDir, err := os.Getwd()
	check(err)
	dataFolder := filepath.Join(workDir, "data")
	rawFolder := filepath.Join(dataFolder, "raw")
	err = os.MkdirAll(rawFolder, 0755)
	check(err)

	db, err := sql.Open("sqlite3", filepath.Join(dataFolder, "storage.db"))
	check(err)
	err = createTables(db)
	check(err)

	return &RuntimeState{
		dataFolder,
		rawFolder,
		db,
		nil,
		nil,
	}
}

// Load or create BTree index
func (s *RuntimeState) loadIndex() error {
	gob.Register(BinaryTree{})
	indexPath := filepath.Join(s.dataFolder, "index.gob")

	index := &BinaryTree{nil}
	err := LoadStruct(indexPath, &index)
	if err != nil {
		fmt.Println("Error reading index:", err, " - attempting to recover index")
		index, err = s.recoverIndex()
		if err != nil {
			fmt.Println("Index recovery failed, all documents are lost!", err)
			return err
		}
		fmt.Println("BTree index succesfully recovered")
	}
	s.index = index
	return nil
}

func (s *RuntimeState) recoverIndex() (*BinaryTree, error) {
	docIDs, err := ListDocuments(s.db)
	if err != nil {
		return nil, err
	}
	btree := &BinaryTree{nil}
	for _, docID := range docIDs {
		doc, err := LoadDocSummary(s.db, docID)
		if err != nil {
			return nil, err
		}
		btree.InsertDoc(doc)
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
	err := LoadStruct(dcPath, docCounter)
	if err != nil {
		fmt.Println("Error reading docCounter:", err, " - attempting to recover")
		docCounter, err = s.recoverDocCounter()
		if err != nil {
			fmt.Println("DocCounter recovery failed, all documents are lost!")
			return err
		}
		fmt.Println("docCounter succesfully recovered")
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
		doc, err := LoadDocSummary(s.db, docID)
		if err != nil {
			return nil, err
		}
		docCounter.AddDocuments(doc)
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
	docSummary := NewDocSummary(text, identifier, title, docType)
	rows, err := InsertDocument(state.db, docSummary, text)
	if err != nil {
		return err
	}
	if rows == 0 {
		fmt.Println("Document is already in the collection")
		return nil
	}

	state.index.InsertDoc(docSummary)
	state.docCounter.AddDocuments(docSummary)

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

	similarities := TFIDFSimilarity(text, state.rawFolder, state.docCounter, docSummaries...)
	return similarities, nil
}

func printSearchResults(sims []*SearchResult, rawFolder string) {
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
