package main

import (
	"context"
	"encoding/base64"
)

// App struct
type App struct {
	ctx   context.Context
	state *RuntimeState
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.state = NewRuntimeState()
	err := a.state.loadIndex()
	check(err)
	err = a.state.loadCounter()
	check(err)

}

// Add content to collection. Argument is base64-encoded
func (a *App) AddContent(encodedContent string) error {
	var err error
	bytes, err := base64.StdEncoding.DecodeString(encodedContent)
	content := string(bytes)
	if err != nil {
		return err
	}
	URLPrefix := URLRegex.FindString(content)
	if URLPrefix != "" {
		err = addURL(content, a.state)
	} else {
		err = addText(content, a.state)
	}
	return err
}

// Search a given query in the collection
func (a *App) Search(text string) ([]*SearchResult, error) {
	return queryDocument(text, a.state)
}

// Read contents from a raw text file stored in the collection
func (a *App) ReadTextFile(docID string) (string, error) {
	return LoadText(filepath.Join(a.state.rawFolder, docID+".txt"))
}
