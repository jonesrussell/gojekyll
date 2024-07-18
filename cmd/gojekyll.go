package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/epiclabs-io/winman"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"jonesrussell/jekyll-publisher/filehandler"
	"jonesrussell/jekyll-publisher/logger"
	"jonesrussell/jekyll-publisher/ui"
)

type App struct {
	fileHandler *filehandler.FileHandler
	ui          *ui.UI
	logger      logger.LoggerInterface
	wm          *winman.Manager
}

type AppContext struct {
	sitePath    string
	tviewApp    *tview.Application
	menu        *tview.TreeView
	contentView *tview.TextView
	gitView     tview.Primitive
	dashboard   *tview.Flex
}

func NewApp(fileHandler *filehandler.FileHandler, ui *ui.UI, logger logger.LoggerInterface) *App {
	return &App{
		fileHandler: fileHandler,
		ui:          ui,
		logger:      logger,
		wm:          winman.NewWindowManager(),
	}
}

func (a *App) createDashboardContext(sitePath string, tviewApp *tview.Application) (*AppContext, error) {
	// Get drafts and posts
	drafts, err := a.fileHandler.GetFilenames(sitePath, "_drafts")
	if err != nil {
		return nil, err
	}
	posts, err := a.fileHandler.GetFilenames(sitePath, "_posts")
	if err != nil {
		return nil, err
	}

	// Create the dashboard with drafts and posts
	dashboard, menu, contentView, gitView, err := a.ui.CreateDashboard(drafts, posts)
	if err != nil {
		return nil, err
	}

	// Create a new resizable window with some text
	a.createResizableWindow(tviewApp, "Hello, world!")

	return &AppContext{
		sitePath:    sitePath,
		tviewApp:    tviewApp,
		menu:        menu,
		contentView: contentView,
		gitView:     gitView,
		dashboard:   dashboard,
	}, nil
}

func (a *App) createResizableWindow(tviewApp *tview.Application, contentText string) {
	content := tview.NewTextView().
		SetText(contentText).
		SetTextAlign(tview.AlignCenter)

	window := a.wm.NewWindow().
		Show().
		SetRoot(content).
		SetDraggable(true).
		SetResizable(true).
		SetTitle("Hi!").
		AddButton(&winman.Button{
			Symbol:  'X',
			OnClick: func() { tviewApp.Stop() },
		})

	window.SetRect(5, 5, 30, 10)
}

func (a *App) Run(args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := args[1]
	tviewApp := tview.NewApplication()

	ctx, err := a.createDashboardContext(sitePath, tviewApp)
	if err != nil {
		a.logger.Error("Could not create dashboard", err)
		return
	}

	// Set the selected function to handle "Exit" and preview content
	a.setMenuSelectedFunc(ctx)

	// Set input capture to switch focus on Tab key press
	a.setInputCapture(ctx)

	if err := tviewApp.SetRoot(ctx.dashboard, true).Run(); err != nil {
		log.Println("Could not set root")
		panic(err)
	}
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

	// TODO: logger that truncates
	// a.logger.Debug(string(content))
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
	modal := a.createPublishModal(ctx)
	ctx.tviewApp.SetRoot(modal, false).SetFocus(modal)
}

func (a *App) publishModalDoneFunc(ctx *AppContext) func(int, string) {
	return func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Publish" {
			a.publishSelectedDraft(ctx)
		}
		// Dismiss the modal and return to the previous view
		ctx.tviewApp.SetRoot(ctx.dashboard, true)
	}
}

func (a *App) createPublishModal(ctx *AppContext) *tview.Modal {
	// Get the currently selected node
	node := ctx.menu.GetCurrentNode()
	// Get the name of the draft
	draftName := node.GetText()

	return tview.NewModal().
		SetText(fmt.Sprintf("Do you want to publish the draft '%s'?", draftName)).
		AddButtons([]string{"Publish", "Cancel"}).
		SetDoneFunc(a.publishModalDoneFunc(ctx))
}

func (a *App) getFilePath(node *tview.TreeNode) string {
	return path.Join("_drafts", node.GetText())
}

func (a *App) assembleNewPathAndFilename(node *tview.TreeNode) (string, string) {
	newFilename := time.Now().Format("2006-01-02") + "-" + node.GetText()
	newPath := path.Join("_posts", newFilename)
	return newPath, newFilename
}

func (a *App) moveFile(ctx *AppContext, filePath string, newPath string) error {
	_, err := exec.Command("git", "-C", ctx.sitePath, "mv", filePath, newPath).Output()
	return err
}

func (a *App) updateUI(ctx *AppContext, node *tview.TreeNode, pathNodes []*tview.TreeNode, newFilename string) {
	pathNodes[1].RemoveChild(node)
	postsNode := pathNodes[0].GetChildren()[1]
	node.SetText(newFilename) // Update the node text with the new filename
	postsNode.AddChild(node)
	ctx.menu.SetCurrentNode(node)
}

func (a *App) publishSelectedDraft(ctx *AppContext) {
	a.logger.Debug("Publish")
	// Get the currently selected node
	node := ctx.menu.GetCurrentNode()
	// Get the path of the selected node
	pathNodes := ctx.menu.GetPath(node)
	// Get the path of the selected file
	filePath := a.getFilePath(node)
	// Create the new path with the date prepended to the filename
	newPath, newFilename := a.assembleNewPathAndFilename(node)

	// Debug logging before moving the file
	cwd, _ := os.Getwd()
	a.logger.Debug(fmt.Sprintf("Current working directory: '%s'", cwd))
	a.logger.Debug(fmt.Sprintf("Moving file from '%s' to '%s'", filePath, newPath))

	// Move the file
	err := a.moveFile(ctx, filePath, newPath)
	if err != nil {
		a.logger.Error("Could not move file", err, "path", filePath)
		return
	}

	// Debug logging after moving the file
	a.logger.Debug(fmt.Sprintf("Successfully moved file from '%s' to '%s'", filePath, newPath))

	// Update the UI
	a.updateUI(ctx, node, pathNodes, newFilename)
}
