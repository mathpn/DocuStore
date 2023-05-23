package main

import (
	"context"
	"encoding/base64"
	"strings"
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
	state, err := NewRuntimeState()
	check(err)
	a.state = state
	err = a.state.loadIndex()
	check(err)
	err = a.state.loadCounter()
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
	return addURL(content, a.state)
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
	return addText(content, title, a.state)
}

// Search a given query in the collection
func (a *App) Search(text string) ([]*SearchResult, error) {
	return queryDocument(text, a.state)
}

// Read contents from a raw text file stored in the collection
func (a *App) ReadTextFile(docID string) (string, error) {
	return LoadText(a.state.db, docID)
}
