package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"jonesrussell/jekyll-publisher/filehandler"
	"jonesrussell/jekyll-publisher/logger"
	"jonesrussell/jekyll-publisher/ui"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type App struct {
	fileHandler filehandler.FileHandler
	ui          ui.UI
	logger      logger.LoggerInterface
}

func NewApp(fileHandler filehandler.FileHandler, ui ui.UI, logger logger.LoggerInterface) *App {
	return &App{
		fileHandler: fileHandler,
		ui:          ui,
		logger:      logger,
	}
}

func (a *App) Run(args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := args[1]
	tviewApp := tview.NewApplication()

	// Get drafts and posts
	drafts, err := a.getFilenames(sitePath, "_drafts")
	if err != nil {
		log.Println(err)
		return
	}
	posts, err := a.getFilenames(sitePath, "_posts")
	if err != nil {
		log.Println(err)
		return
	}

	// Create the dashboard with drafts and posts
	dashboard, menu, contentView, gitView := a.ui.CreateDashboard(sitePath, drafts, posts)

	// Set the selected function to handle "Exit" and preview content
	a.setMenuSelectedFunc(menu, tviewApp, sitePath, contentView)

	// Set input capture to switch focus on Tab key press
	a.setInputCapture(tviewApp, gitView, menu, contentView)

	if err := tviewApp.SetRoot(dashboard, true).Run(); err != nil {
		log.Println("Could not set root")
		panic(err)
	}
}

func (a *App) getFilenames(sitePath, dir string) ([]string, error) {
	return a.fileHandler.GetFilenames(sitePath, dir)
}

func (a *App) setMenuSelectedFunc(menu *tview.TreeView, tviewApp *tview.Application, sitePath string, contentView *tview.TextView) {
	menu.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.GetText() == "Exit" {
			tviewApp.Stop()
		} else {
			a.handleFileSelection(node, sitePath, contentView, menu)
		}
	})
}

func (a *App) handleFileSelection(node *tview.TreeNode, sitePath string, contentView *tview.TextView, menu *tview.TreeView) {
	// Get the path of the selected node
	pathNodes := menu.GetPath(node)

	// Determine the directory of the selected file
	var dir string
	if pathNodes[1].GetText() == "Drafts" {
		dir = "_drafts"
	} else if pathNodes[1].GetText() == "Posts" {
		dir = "_posts"
	} else {
		return
	}

	// Get the path of the selected file
	filePath := path.Join(sitePath, dir, node.GetText())

	// Read the content of the file
	content, err := a.fileHandler.ReadFile(filePath)
	if err != nil {
		a.logger.Error("Could not read file", err, "path", filePath)
		return
	}

	a.logger.Debug(string(content))
	// Display the content in contentView
	contentView.SetText(string(content))
}

func (a *App) setInputCapture(tviewApp *tview.Application, gitView, menu, contentView tview.Primitive) {
	tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if tviewApp.GetFocus() == gitView {
				tviewApp.SetFocus(menu)
			} else if tviewApp.GetFocus() == menu {
				tviewApp.SetFocus(contentView)
			} else {
				tviewApp.SetFocus(gitView)
			}
		}
		return event
	})
}
