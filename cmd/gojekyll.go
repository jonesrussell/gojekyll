package cmd

import (
	"fmt"
	"log"
	"os"

	"jonesrussell/jekyll-publisher/filehandler"
	"jonesrussell/jekyll-publisher/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	fileHandler filehandler.FileHandler
	ui          ui.UI
}

func NewApp(fileHandler filehandler.FileHandler, ui ui.UI) *App {
	return &App{
		fileHandler: fileHandler,
		ui:          ui,
	}
}

func (a *App) Run(args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := args[1]
	app := tview.NewApplication()

	// Add drafts and posts to the lists
	// Get drafts and posts
	drafts, err := a.fileHandler.GetFilenames(sitePath, "_drafts")
	if err != nil {
		log.Println(err)
		return
	}
	posts, err := a.fileHandler.GetFilenames(sitePath, "_posts")
	if err != nil {
		log.Println(err)
		return
	}

	// Create the dashboard with drafts and posts
	dashboard, menu, contentView := a.ui.CreateDashboard(sitePath, drafts, posts)

	// Set input capture to switch focus on Tab key press
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if app.GetFocus() == menu {
				app.SetFocus(contentView)
			} else {
				app.SetFocus(menu)
			}
		}
		return event
	})

	if err := app.SetRoot(dashboard, true).Run(); err != nil {
		log.Println("Could not set root")
		panic(err)
	}
}
