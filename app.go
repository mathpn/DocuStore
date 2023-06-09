package main

import (
	"context"
	"encoding/base64"
	"strings"
)

// App struct
type App struct {
	ctx    context.Context
	engine *DocuEngine
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	state, err := NewEngine()
	check(err)
	a.engine = state
	err = a.engine.loadIndex()
	check(err)
	err = a.engine.loadCounter()
	check(err)
}

// Decode base64-encoded input
func (a *App) decodeInput(encodedContent string) (string, error) {
	var err error
	bytes, err := base64.StdEncoding.DecodeString(encodedContent)
	content := string(bytes)
	if err != nil {
		return "", err
	}
	content = strings.TrimSpace(content)
	return content, nil
}

func (a *App) AddURL(encodedURL string) error {
	var err error
	content, err := a.decodeInput(encodedURL)
	if err != nil {
		return err
	}
	return a.engine.addURL(content)
}

func (a *App) AddText(encodedText string, encodedTitle string) error {
	var err error
	content, err := a.decodeInput(encodedText)
	if err != nil {
		return err
	}
	title, err := a.decodeInput(encodedTitle)
	if err != nil {
		return err
	}
	return a.engine.addText(content, title)
}

// Search a given query in the collection
func (a *App) Search(text string) ([]*SearchResult, error) {
	return a.engine.queryDocument(text)
}

// Read contents from a raw text file stored in the collection
func (a *App) ReadTextFile(docID string) (string, error) {
	return LoadText(a.engine.db, docID)
}
