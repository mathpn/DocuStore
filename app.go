package main

import (
	"context"
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

func (a *App) AddContent(content string) error {
	URLPrefix := URLRegex.FindString(content)
	var err error
	if URLPrefix != "" {
		err = addURL(content, a.state)
	} else {
		err = addText(content, a.state)
	}
	return err
}

// something here
func (a *App) Search(text string) ([]*SearchResult, error) {
	return queryDocument(text, a.state)
}
