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

	"DocuStore/scraper"
	"DocuStore/search"

	"github.com/adrg/xdg"
)

type DocuEngine struct {
	db         *sql.DB
	index      *HashmapIndex
	docCounter *search.DocCounter
	searcher   search.Searcher
	dataFolder string
}

func NewEngine() (*DocuEngine, error) {
	gob.Register(search.DocSummary{})
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

	index, err := loadIndex(dataFolder, db)
	if err != nil {
		return nil, err
	}
	docCounter, err := loadCounter(dataFolder, db)
	if err != nil {
		return nil, err
	}

	searcher, err := search.NewTFIDFSearcher(docCounter)
	if err != nil {
		return nil, err
	}
	engine := &DocuEngine{
		db:         db,
		index:      index,
		docCounter: docCounter,
		searcher:   searcher,
		dataFolder: dataFolder,
	}
	return engine, nil
}

// Load or create Hashmap inverted index
func loadIndex(dataFolder string, db *sql.DB) (*HashmapIndex, error) {
	gob.Register(HashmapIndex{})
	indexPath := filepath.Join(dataFolder, "index.gob")
	// fmt.Println(indexPath)

	index := &HashmapIndex{nil, 0}
	latestTs, err := GetLatestTimestamp(db)
	if err != nil {
		return index, err
	}
	if latestTs == 0 {
		// no documents
		return index, nil
	}

	err = LoadStruct(indexPath, &index)
	if err != nil {
		fmt.Println("Error reading BTree index:", err, "- attempting to recover")
		index, err = recoverIndex(dataFolder, db)
		if err != nil {
			fmt.Println("Index recovery failed, documents may have been lost:", err)
			return nil, err
		}
		fmt.Println("BTree index succesfully recovered")
	}

	if index.Timestamp != latestTs {
		fmt.Printf("%+v - %+v\n", index.Timestamp, latestTs)
		fmt.Println("BTree index is out of sync with latest changes, recovering")
		index, err = recoverIndex(dataFolder, db)
		if err != nil {
			fmt.Println("BTree index recovery failed, documents may have been lost:", err)
			return nil, err
		}
	}
	return index, nil
}

func recoverIndex(dataFolder string, db *sql.DB) (*HashmapIndex, error) {
	docIDs, err := ListDocuments(db)
	if err != nil {
		return nil, err
	}
	btree := &HashmapIndex{nil, 0}
	var doc *search.DocSummary
	var ts int64
	for _, docID := range docIDs {
		doc, ts, err = LoadDocSummary(db, docID)
		if err != nil {
			return nil, err
		}
		btree.InsertDoc(doc, ts)
	}
	indexPath := filepath.Join(dataFolder, "index.gob")
	err = SaveStruct(indexPath, btree)
	if err != nil {
		return nil, err
	}
	return btree, nil
}

// Load or create DocCounter
func loadCounter(dataFolder string, db *sql.DB) (*search.DocCounter, error) {
	gob.Register(search.DocCounter{})
	dcPath := filepath.Join(dataFolder, "docCounter.gob")

	docCounter := search.NewDocCounter()
	latestTs, err := GetLatestTimestamp(db)
	if err != nil {
		return nil, err
	}
	if latestTs == 0 {
		// no documents
		return docCounter, nil
	}

	err = LoadStruct(dcPath, docCounter)
	if err != nil {
		fmt.Println("Error reading docCounter:", err, "- attempting to recover")
		docCounter, err = recoverDocCounter(dataFolder, db)
		if err != nil {
			fmt.Println("DocCounter recovery failed, documents may have been lost:", err)
			return nil, err
		}
		fmt.Println("docCounter succesfully recovered")
	}

	if docCounter.Ts != latestTs {
		fmt.Println("DocCounter is out of sync with latest changes, recovering")
		docCounter, err = recoverDocCounter(dataFolder, db)
		if err != nil {
			fmt.Println("DocCounter recovery failed, documents may have been lost:", err)
			return nil, err
		}
	}

	return docCounter, nil
}

func recoverDocCounter(dataFolder string, db *sql.DB) (*search.DocCounter, error) {
	docIDs, err := ListDocuments(db)
	if err != nil {
		return nil, err
	}

	docCounter := search.NewDocCounter()
	var doc *search.DocSummary
	var ts int64
	for _, docID := range docIDs {
		doc, ts, err = LoadDocSummary(db, docID)
		if err != nil {
			return nil, err
		}
		docCounter.AddDocument(doc, ts)
	}

	dcPath := filepath.Join(dataFolder, "docCounter.gob")
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
	e.addDocument(text, text, filePath, search.DocType(search.Text))
	return nil
}

func (e *DocuEngine) AddText(text string, title string) error {
	err := e.addDocument(text, text, title, search.DocType(search.Text))
	return err

}

func (e *DocuEngine) AddURL(url string) error {
	data, err := scraper.ScrapeText(url)
	if err != nil {
		return err
	}
	err = e.addDocument(data.Content, url, data.Title, search.DocType(search.URL))
	return err
}

func (e *DocuEngine) addDocument(text string, identifier string, title string, docType search.DocType) error {
	if title == "" {
		return errors.New("empty title is not allowed")
	}
	if text == "" {
		return errors.New("empty content")
	}
	ts := time.Now().Unix()
	docSummary := search.NewDocSummary(text, identifier, title, docType)
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

func (e *DocuEngine) QueryDocument(text string) ([]*search.SearchResult, error) {
	tokens := search.Tokenize(text)
	docIDs := e.index.SearchTokens(tokens)
	docSummaries, err := LoadDocSummaries(context.Background(), e.db, docIDs...)
	if err != nil {
		return nil, err
	}

	similarities := e.searcher.Search(text, docSummaries...)
	return similarities, nil
}

func (e *DocuEngine) LoadText(docID string) (string, error) {
	return LoadText(e.db, docID)
}

func printSearchResults(sims []*search.SearchResult) {
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
		if sim.Type == search.URL.String() {
			fmt.Println(sim.Identifier)
		}
		fmt.Println("--------------------------------------------")
	}
}
