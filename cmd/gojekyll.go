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
	ui          *ui.UI
	logger      logger.LoggerInterface
}

type AppContext struct {
	sitePath    string
	tviewApp    *tview.Application
	menu        *tview.TreeView
	contentView *tview.TextView
	gitView     tview.Primitive
}

func NewApp(fileHandler filehandler.FileHandler, ui *ui.UI, logger logger.LoggerInterface) *App {
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
	dashboard, menu, contentView, gitView := a.ui.CreateDashboard(drafts, posts)

	ctx := &AppContext{
		sitePath:    sitePath,
		tviewApp:    tviewApp,
		menu:        menu,
		contentView: contentView,
		gitView:     gitView,
	}

	// Set the selected function to handle "Exit" and preview content
	a.setMenuSelectedFunc(ctx)

	// Set input capture to switch focus on Tab key press
	a.setInputCapture(ctx)

	if err := tviewApp.SetRoot(dashboard, true).Run(); err != nil {
		log.Println("Could not set root")
		panic(err)
	}
}

func (a *App) getFilenames(sitePath, dir string) ([]string, error) {
	return a.fileHandler.GetFilenames(sitePath, dir)
}

func (a *App) setMenuSelectedFunc(ctx *AppContext) {
	ctx.menu.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.GetText() == "Exit" {
			ctx.tviewApp.Stop()
		} else {
			a.handleFileSelection(node, ctx)
		}
	})
}

func (a *App) handleFileSelection(node *tview.TreeNode, ctx *AppContext) {
	// Get the path of the selected node
	pathNodes := ctx.menu.GetPath(node)

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
	filePath := path.Join(ctx.sitePath, dir, node.GetText())

	// Read the content of the file
	content, err := a.fileHandler.ReadFile(filePath)
	if err != nil {
		a.logger.Error("Could not read file", err, "path", filePath)
		return
	}

	a.logger.Debug(string(content))
	// Display the content in contentView
	ctx.contentView.SetText(string(content))
}

func (a *App) setInputCapture(ctx *AppContext) {
	ctx.tviewApp.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if ctx.tviewApp.GetFocus() == ctx.gitView {
				ctx.tviewApp.SetFocus(ctx.menu)
			} else if ctx.tviewApp.GetFocus() == ctx.menu {
				ctx.tviewApp.SetFocus(ctx.contentView)
			} else {
				ctx.tviewApp.SetFocus(ctx.gitView)
			}
		} else if event.Rune() == 'p' {
			a.showPublishModal(ctx)
		}
		return event
	})
}

func (a *App) showPublishModal(ctx *AppContext) {
	modal := tview.NewModal().
		SetText("Do you want to publish the selected draft?").
		AddButtons([]string{"Publish", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Publish" {
				a.logger.Debug("Publish")
				// Add your publishing logic here
			}
		})
	ctx.tviewApp.SetRoot(modal, false).SetFocus(modal)
}
