package main

import (
	"embed"
	"flag"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runApp() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "DocuStore",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func cliInterface() {
	var err error
	state, err := NewRuntimeState()
	check(err)
	err = state.loadIndex()
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
			err = addURL(arg, state)
			check(err)
		} else {
			err = addFile(arg, state)
			check(err)
		}
	case "query":
		fmt.Println("querying documents")
		query := flag.Arg(1)
		if query == "" {
			fmt.Println("You must provide a query string.")
			return
		}
		// FIXME loop for profiling only, remove later
		for i := 0; i < 100; i++ {
			result, err := queryDocument(query, state)
			if err != nil {
				panic(err)
			}
			printSearchResults(result)
		}
	default:
		fmt.Println("Valid commands: add, query")
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		runApp()
	} else {
		cliInterface()
	}
}
