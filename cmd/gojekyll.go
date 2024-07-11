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

	dashboard, draftsList, postsList, _ := a.ui.CreateDashboard(sitePath)

	// Add drafts and posts to the lists
	drafts, err := a.fileHandler.GetFilenames(sitePath, "_drafts")
	addItemsToList(draftsList, drafts, err)
	posts, err := a.fileHandler.GetFilenames(sitePath, "_posts")
	addItemsToList(postsList, posts, err)

	// Add a special "Exit" item to the list
	draftsList.AddItem("Exit", "", 0, func() {
		app.Stop()
	})
	postsList.AddItem("Exit", "", 0, func() {
		app.Stop()
	})

	// Create a slice of the lists to switch focus between
	lists := []*tview.List{draftsList, postsList}
	focusIndex := 0

	// Set input capture to switch focus on Tab key press
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			focusIndex = (focusIndex + 1) % len(lists)
			app.SetFocus(lists[focusIndex])
		}
		return event
	})

	if err := app.SetRoot(dashboard, true).Run(); err != nil {
		log.Println("Could not set root")
		panic(err)
	}
}

func addItemsToList(list *tview.List, items []string, err error) {
	if err != nil {
		log.Println(err)
		return
	}
	for _, item := range items {
		list.AddItem(item, "", 0, nil)
	}
}
