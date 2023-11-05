package main

import (
	"embed"
	"flag"
	"fmt"

	"DocuStore/scraper"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func runApp() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:              "DocuStore",
		Width:              1024,
		Height:             768,
		WindowStartState:   options.Normal,
		Frameless:          false,
		MinWidth:           300,
		LogLevel:           logger.DEBUG,
		LogLevelProduction: logger.WARNING,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		panic(err)
	}
}

func cliInterface() {
	var err error
	engine, err := NewEngine()
	if err != nil {
		panic(err)
	}

	cmd := flag.Arg(0)
	switch cmd {
	case "add":
		fmt.Println("adding document")
		arg := flag.Arg(1)
		if arg == "" {
			fmt.Println("You must provide a valid file path or URL.")
			return
		}
		found := scraper.URLRegex.FindString(arg)
		if found != "" {
			err = engine.AddURL(arg)
			if err != nil {
				panic(err)
			}
		} else {
			err = engine.addFile(arg)
			if err != nil {
				panic(err)
			}
		}
	case "query":
		fmt.Println("querying documents")
		query := flag.Arg(1)
		if query == "" {
			fmt.Println("You must provide a query string.")
			return
		}
		result, err := engine.QueryDocument(query)
		if err != nil {
			panic(err)
		}
		printSearchResults(result)
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
