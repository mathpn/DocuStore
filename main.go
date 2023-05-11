package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

type RuntimeState struct {
	dataFolder    string
	summaryFolder string
	rawFolder     string
	index         *BinaryTree
}

func NewRuntimeState() *RuntimeState {
	gob.Register(BinaryTree{})
	gob.Register(DocSummary{})
	workDir, err := os.Getwd()
	check(err)
	dataFolder := filepath.Join(workDir, "data")
	summaryFolder := filepath.Join(dataFolder, "summary")
	rawFolder := filepath.Join(dataFolder, "raw")
	err = os.MkdirAll(rawFolder, 0755)
	check(err)
	err = os.MkdirAll(summaryFolder, 0755)
	check(err)

	// load or create index
	indexPath := filepath.Join(dataFolder, "index.gob")
	var index BinaryTree
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		fmt.Println("Index not found")
	} else {
		err := LoadStruct(indexPath, &index)
		if err != nil {
			fmt.Println("Error reading index:", err)
			index = BinaryTree{nil}
		}
	}
	return &RuntimeState{
		dataFolder,
		summaryFolder,
		rawFolder,
		&index,
	}
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
	summaryPath := filepath.Join(state.summaryFolder, docSummary.DocID+".gob")
	if fileExists(summaryPath) {
		fmt.Println("Document is already in the collection")
		return
	}
	state.index.InsertDoc(docSummary)
	SaveStruct(
		filepath.Join(state.summaryFolder, docSummary.DocID+".gob"),
		docSummary,
	)
	SaveStruct(
		filepath.Join(state.dataFolder, "index.gob"),
		state.index,
	)
	SaveText(
		filepath.Join(state.rawFolder, docSummary.DocID+".txt"),
		text,
	)
}

func queryDocument(text string, state *RuntimeState) {
	PrintTree(os.Stdout, state.index.Root, 0, 'M')
	tokens := Tokenize(text)
	docIDs := state.index.SearchDoc(tokens)
	docSummaries := make([]*DocSummary, len(docIDs))
	for i, docID := range docIDs {
		fpath := filepath.Join(state.summaryFolder, docID+".gob")
		LoadStruct(fpath, &docSummaries[i])
	}

	// files, err := os.ReadDir(state.summaryFolder)
	// docSummaries := make([]*DocSummary, len(files))
	// check(err)
	// for i, file := range files {
	// 	fpath := filepath.Join(state.summaryFolder, file.Name())
	// 	LoadStruct(fpath, &docSummaries[i])
	// }
	tfidf := NewTFIDF()
	tfidf.AddDocuments(docSummaries...)
	similarities := tfidf.Similarity(text)
	printSimilarities(similarities, state.rawFolder)
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
		queryDocument(query, state)
	default:
		fmt.Println("Valid commands: add, query")
	}
}
