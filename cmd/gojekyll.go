package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"jonesrussell/gojekyll/filehandler"
	"jonesrussell/gojekyll/logger"
	"jonesrussell/gojekyll/ui"
)

type App struct {
	fileHandler *filehandler.FileHandler
	ui          *ui.UI
	logger      logger.LoggerInterface
	wm          *winman.Manager
	sitePath    string
}

// NewApp creates a new App instance
func NewApp(fileHandler *filehandler.FileHandler, ui *ui.UI, logger logger.LoggerInterface) *App {
	return &App{
		fileHandler: fileHandler,
		ui:          ui,
		logger:      logger,
		wm:          winman.NewWindowManager(),
	}
}

// Run starts the application
func (a *App) Run(args []string) {
	a.handleSitePathArg(args)
	tviewApp := tview.NewApplication()

	ctx, err := a.createDashboardContext(tviewApp)
	if err != nil {
		a.logger.Error("Could not create dashboard", err)
		return
	}

	a.setMenuSelectedFunc(ctx)
	a.setInputCapture(ctx)

	// Create a new text view for the status bar
	statusBar := a.ui.CreateStatusBar()

	// Create a new flex for the layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(a.wm, 0, 1, true).      // The WindowManager takes up the rest of the space
		AddItem(statusBar, 1, 1, false) // The status bar has a fixed height of 1

	if err := tviewApp.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		log.Println("Could not set root")
		panic(err)
	}
}

// handleSitePathArg handles the site path argument
func (a *App) handleSitePathArg(args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}
	a.sitePath = args[1]
}

// createDashboardContext creates the dashboard context
func (a *App) createDashboardContext(tviewApp *tview.Application) (*AppContext, error) {
	// Get drafts and posts
	drafts, err := a.fileHandler.GetFilenames(a.sitePath, "_drafts")
	if err != nil {
		return nil, err
	}
	posts, err := a.fileHandler.GetFilenames(a.sitePath, "_posts")
	if err != nil {
		return nil, err
	}

	// Create the dashboard with drafts and posts
	dashboard, menu, contentView, gitView, err := a.ui.CreateDashboard(drafts, posts)
	if err != nil {
		return nil, err
	}

	// Create a new resizable window with some text
	a.ui.CreateResizableWindow("Blog Posts", menu, a.wm)
	a.ui.CreateResizableWindow("Content View", contentView, a.wm)
	a.ui.CreateResizableWindow("Git View", gitView, a.wm)

	// Create a new text view for the status bar
	statusBar := a.ui.CreateStatusBar()

	return &AppContext{
		tviewApp:    tviewApp,
		menu:        menu,
		contentView: contentView,
		gitView:     gitView,
		dashboard:   dashboard,
		statusBar:   statusBar,
	}, nil
}

// Refactored publishSelectedDraft function
func (a *App) publishSelectedDraft(ctx *AppContext) {
	a.logger.Debug("Publish")

	node := a.getCurrentNode(ctx)
	filePath := a.getFilePath(node)
	newPath, newFilename := a.assembleNewPathAndFilename(node)

	a.logger.Debug(fmt.Sprintf("Moving file from '%s' to '%s'", filePath, newPath))

	if err := a.fileHandler.MoveFile(filePath, newPath, a.sitePath); err != nil {
		a.logger.Error("Could not move file", err, "path", filePath)
		return
	}

	a.logger.Debug(fmt.Sprintf("Successfully moved file from '%s' to '%s'", filePath, newPath))
	a.ui.UpdateUI(ctx.menu, node, newFilename)
}

// New function to get current node and path
func (a *App) getCurrentNode(ctx *AppContext) *tview.TreeNode {
	node := ctx.menu.GetCurrentNode()
	return node
}

func (a *App) setMenuSelectedFunc(ctx *AppContext) {
	ctx.menu.SetSelectedFunc(func(node *tview.TreeNode) {
		if node.GetText() == "Exit" {
			ctx.tviewApp.Stop()
		} else {
			// Get the path of the selected node
			pathNodes := ctx.menu.GetPath(node)

			// Check if a file is selected
			if len(pathNodes) > 1 && node.GetText() != "Drafts" && node.GetText() != "Posts" {
				// A file is selected, handle it
				a.handleFileSelection(node, ctx)
			}
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
		a.logger.Debug("Should be updating status")
		ctx.statusBar.SetText("Press 'p' to publish the draft") // Set the status bar text
	} else if pathNodes[1].GetText() == "Posts" {
		dir = "_posts"
		ctx.statusBar.SetText("") // Clear the status bar text
	} else {
		return
	}

	// Get the path of the selected file
	filePath := path.Join(a.sitePath, dir, node.GetText())

	// Read the content of the file
	content, err := a.fileHandler.ReadFile(filePath)
	if err != nil {
		a.logger.Error("Could not read file", err, "path", filePath)
		return
	}

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
	// Get the currently selected node
	node := ctx.menu.GetCurrentNode()

	modal := a.ui.CreatePublishModal(node, a.publishModalDoneFunc(ctx))
	ctx.tviewApp.SetRoot(modal, false).SetFocus(modal)
}

func (a *App) publishModalDoneFunc(ctx *AppContext) func(int, string) {
	return func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Publish" {
			a.publishSelectedDraft(ctx)
		}
		// Dismiss the modal and return to the previous view
		ctx.tviewApp.SetRoot(a.wm, true)
	}
}

func (a *App) getFilePath(node *tview.TreeNode) string {
	return path.Join("_drafts", node.GetText())
}

func (a *App) assembleNewPathAndFilename(node *tview.TreeNode) (string, string) {
	newFilename := time.Now().Format("2006-01-02") + "-" + node.GetText()
	newPath := path.Join("_posts", newFilename)
	return newPath, newFilename
}
