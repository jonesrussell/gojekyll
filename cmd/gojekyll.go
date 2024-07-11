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

func Run(args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := args[1]
	app := tview.NewApplication()

	dashboard, draftsList, postsList := ui.CreateDashboard()

	drafts, err := filehandler.GetFilenames(sitePath, "_drafts")
	if err != nil {
		log.Println(err)
		return
	}
	for _, draft := range drafts {
		draftsList.AddItem(draft, "", 0, nil)
	}

	posts, err := filehandler.GetFilenames(sitePath, "_posts")
	if err != nil {
		log.Println(err)
		return
	}
	for _, post := range posts {
		postsList.AddItem(post, "", 0, nil)
	}

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
