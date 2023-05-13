package main

import (
	"context"
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
		doc, err := LoadDocument(s.db, docID)
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
		doc, err := LoadDocument(s.db, docID)
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

var newLineRegex = regexp.MustCompile(`\s`)

func titleFromText(text string) string {
	title := text[0:50]
	title = newLineRegex.ReplaceAllString(title, " ")
	return title
}

func addFile(filePath string, state *RuntimeState) {
	content, err := os.ReadFile(filePath)
	check(err)
	text := string(content)
	title := titleFromText(text)
	addDocument(text, text, title, state)

}

func addURL(url string, state *RuntimeState) {
	title, text := ScrapeText(url)
	addDocument(text, url, title, state)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // file exists
	}
	if os.IsNotExist(err) {
		return false // file does not exist
	}
	// if an error other than "file does not exist" occurs, we assume the file exists
	return true
}

func addDocument(text string, identifier string, title string, state *RuntimeState) {
	docSummary := NewDocSummary(text, identifier, title)
	rows, err := InsertDocument(state.db, docSummary)
	if err != nil {
		panic(err)
	}
	if rows == 0 {
		fmt.Println("Document is already in the collection")
		return
	}

	state.index.InsertDoc(docSummary)
	state.docCounter.AddDocuments(docSummary)

	err = SaveStruct(
		filepath.Join(state.dataFolder, "index.gob"),
		state.index,
	)
	check(err)

	err = SaveText(
		filepath.Join(state.rawFolder, docSummary.DocID+".txt"),
		text,
	)
	check(err)

	err = SaveStruct(
		filepath.Join(state.dataFolder, "docCounter.gob"),
		state.docCounter,
	)
	check(err)
}

func queryDocument(text string, state *RuntimeState) []*SimResult {
	// PrintTree(os.Stdout, state.index.Root, 0, 'M')
	tokens := Tokenize(text)
	docIDs := state.index.SearchDoc(tokens)
	docSummaries, err := LoadDocuments(context.Background(), state.db, docIDs...)
	check(err)

	similarities := TFIDFSimilarity(text, state.docCounter, docSummaries...)
	printSimilarities(similarities, state.rawFolder)
	return similarities
}

func printSimilarities(sims []*SimResult, rawFolder string) {
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

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Please provide a command to execute.")
	}

	state := NewRuntimeState()
	err := state.loadIndex()
	check(err)
	err = state.loadCounter()
	check(err)

	cmd := flag.Arg(0)
	switch cmd {
	case "add":
		fmt.Println("adding document")
		arg := flag.Arg(1)
		if arg == "" {
			fmt.Println("You must provide a valid file path or URL.")
			return
		}
		found := URLRegex.FindString(arg)
		if found != "" {
			addURL(arg, state)
		} else {
			addFile(arg, state)
		}
	case "query":
		fmt.Println("querying documents")
		query := flag.Arg(1)
		if query == "" {
			fmt.Println("You must provide a query string.")
			return
		}
		// XXX loop for profiling only
		for i := 0; i < 10; i++ {
			queryDocument(query, state)
		}
	default:
		fmt.Println("Valid commands: add, query")
	}
}
