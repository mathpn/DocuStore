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

type DocuEngine struct {
	dataFolder string
	db         *sql.DB
	index      *BinaryTree
	docCounter *DocCounter
}

func NewEngine() (*DocuEngine, error) {
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

	state := &DocuEngine{
		dataFolder,
		db,
		nil,
		nil,
	}
	return state, nil
}

// Load or create BTree index
func (e *DocuEngine) loadIndex() error {
	gob.Register(BinaryTree{})
	indexPath := filepath.Join(e.dataFolder, "index.gob")

	index := &BinaryTree{nil, 0}
	latestTs, err := GetLatestTimestamp(e.db)
	if err != nil {
		return err
	}
	if latestTs == 0 {
		// no documents
		e.index = index
		return nil
	}

	err = LoadStruct(indexPath, &index)
	if err != nil {
		fmt.Println("Error reading BTree index:", err, "- attempting to recover")
		index, err = e.recoverIndex()
		if err != nil {
			fmt.Println("Index recovery failed, documents may have been lost:", err)
			return err
		}
		fmt.Println("BTree index succesfully recovered")
	}

	if index.Timestamp != latestTs {
		fmt.Printf("%+v - %+v\n", index.Timestamp, latestTs)
		fmt.Println("BTree index is out of sync with latest changes, recovering")
		index, err = e.recoverIndex()
		if err != nil {
			fmt.Println("BTree index recovery failed, documents may have been lost:", err)
			return err
		}
	}

	e.index = index
	return nil
}

func (e *DocuEngine) recoverIndex() (*BinaryTree, error) {
	docIDs, err := ListDocuments(e.db)
	if err != nil {
		return nil, err
	}
	btree := &BinaryTree{nil, 0}
	for _, docID := range docIDs {
		doc, ts, err := LoadDocSummary(e.db, docID)
		if err != nil {
			return nil, err
		}
		btree.InsertDoc(doc, ts)
	}
	indexPath := filepath.Join(e.dataFolder, "index.gob")
	err = SaveStruct(indexPath, btree)
	if err != nil {
		return nil, err
	}
	return btree, nil
}

// Load or create DocCounter
func (e *DocuEngine) loadCounter() error {
	gob.Register(DocCounter{})
	dcPath := filepath.Join(e.dataFolder, "docCounter.gob")

	docCounter := NewDocCounter()
	latestTs, err := GetLatestTimestamp(e.db)
	if err != nil {
		return err
	}
	if latestTs == 0 {
		// no documents
		e.docCounter = docCounter
		return nil
	}

	err = LoadStruct(dcPath, docCounter)
	if err != nil {
		fmt.Println("Error reading docCounter:", err, "- attempting to recover")
		docCounter, err = e.recoverDocCounter()
		if err != nil {
			fmt.Println("DocCounter recovery failed, documents may have been lost:", err)
			return err
		}
		fmt.Println("docCounter succesfully recovered")
	}

	if docCounter.Timestamp != latestTs {
		fmt.Println("DocCounter is out of sync with latest changes, recovering")
		docCounter, err = e.recoverDocCounter()
		if err != nil {
			fmt.Println("DocCounter recovery failed, documents may have been lost:", err)
			return err
		}
	}

	e.docCounter = docCounter
	return nil
}

func (e *DocuEngine) recoverDocCounter() (*DocCounter, error) {
	docIDs, err := ListDocuments(e.db)
	if err != nil {
		return nil, err
	}

	docCounter := NewDocCounter()
	for _, docID := range docIDs {
		doc, ts, err := LoadDocSummary(e.db, docID)
		if err != nil {
			return nil, err
		}
		docCounter.AddDocument(doc, ts)
	}

	dcPath := filepath.Join(e.dataFolder, "docCounter.gob")
	err = SaveStruct(dcPath, docCounter)
	if err != nil {
		return nil, err
	}
	return docCounter, nil
}

func (e *DocuEngine) addFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	text := string(content)
	e.addDocument(text, text, filePath, DocType(Text))
	return nil
}

func (e *DocuEngine) addText(text string, title string) error {
	err := e.addDocument(text, text, title, DocType(Text))
	return err

}

func (e *DocuEngine) addURL(url string) error {
	title, text := ScrapeText(url)
	err := e.addDocument(text, url, title, DocType(URL))
	return err
}

func (e *DocuEngine) addDocument(text string, identifier string, title string, docType DocType) error {
	if title == "" {
		return errors.New("empty title is not allowed")
	}
	if text == "" {
		return errors.New("empty content")
	}
	ts := time.Now().Unix()
	docSummary := NewDocSummary(text, identifier, title, docType)
	rows, err := InsertDocument(e.db, docSummary, text, ts)
	if err != nil {
		return err
	}
	if rows == 0 {
		fmt.Println("Document is already in the collection")
		return nil
	}

	e.index.InsertDoc(docSummary, ts)
	e.docCounter.AddDocument(docSummary, ts)

	err = SaveStruct(
		filepath.Join(e.dataFolder, "index.gob"),
		e.index,
	)
	if err != nil {
		return err
	}

	err = SaveStruct(
		filepath.Join(e.dataFolder, "docCounter.gob"),
		e.docCounter,
	)
	return err
}

func (e *DocuEngine) queryDocument(text string) ([]*SearchResult, error) {
	tokens := Tokenize(text)
	docIDs := e.index.SearchTokens(tokens)
	docSummaries, err := LoadDocSummaries(context.Background(), e.db, docIDs...)
	if err != nil {
		return nil, err
	}

	similarities := TFIDFSimilarity(text, e.docCounter, docSummaries...)
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
