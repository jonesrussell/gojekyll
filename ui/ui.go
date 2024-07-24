package ui

import (
	"fmt"

	"github.com/epiclabs-io/winman"
	"github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

type UIInterface interface {
	CreateDashboard(drafts []string, posts []string) (*winman.Manager, *tview.TreeView, *tview.TextView, *tview.TextView, error)
	CreateGitView() (*tview.TextView, error)
	CreateResizableWindow(title string, content tview.Primitive, wm *winman.Manager)
	CreateStatusBar() *tview.TextView
	CreatePublishModal(node *tview.TreeNode, doneFunc func(int, string)) *tview.Modal
	UpdateUI(menu *tview.TreeView, node *tview.TreeNode, newFilename string)
}

type UI struct {
	sitePath string
}

// Ensure UI implements UIInterface
var _ UIInterface = &UI{}

// NewUI creates a new UI instance
func NewUI(sitePath string) *UI {
	return &UI{
		sitePath: sitePath,
	}
}

// createNode creates a new node with the given title and items
func (ui *UI) createNode(title string, items []string) *tview.TreeNode {
	node := tview.NewTreeNode(title)
	for _, item := range items {
		node.AddChild(tview.NewTreeNode(item))
	}
	return node
}

// CreateDashboard creates a new tview.Flex that contains two lists titled "Drafts" and "Posts".
func (ui *UI) CreateDashboard(drafts []string, posts []string) (*winman.Manager, *tview.TreeView, *tview.TextView, *tview.TextView, error) {
	gitView, err := ui.CreateGitView()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Create a tree for the menu
	menu := tview.NewTreeView()

	// Create root for the tree
	root := tview.NewTreeNode("")

	// Add drafts and posts to the tree
	root.AddChild(ui.createNode("Drafts", drafts))
	root.AddChild(ui.createNode("Posts", posts))

	menu.SetRoot(root).SetCurrentNode(root)

	// Create a text view for the content of the selected draft or post
	contentView := tview.NewTextView()

	// Create a WindowManager
	wm := winman.NewWindowManager()

	// Create windows for the gitView, menu, and contentView
	gitWindow := wm.NewWindow().SetRoot(gitView).SetTitle("Git View")
	menuWindow := wm.NewWindow().SetRoot(menu).SetTitle("Menu")
	contentWindow := wm.NewWindow().SetRoot(contentView).SetTitle("Content View")

	// Add the windows to the WindowManager
	wm.AddWindow(gitWindow)
	wm.AddWindow(menuWindow)
	wm.AddWindow(contentWindow)

	return wm, menu, contentView, gitView, nil
}

// CreateGitView creates a new Git view
func (ui *UI) CreateGitView() (*tview.TextView, error) {
	gitView := tview.NewTextView()

	// Check if the sitePath is a Git repository
	repo, err := git.PlainOpen(ui.sitePath)
	if err != nil {
		return nil, fmt.Errorf((fmt.Sprintf("The directory %s is not a Git repository. Consider running 'git init'.\n", ui.sitePath)))
	}

	// Get the worktree of the repository
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get the worktree of the repository: %v", err)
	}

	// Get the status of the worktree
	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get the status of the worktree: %v", err)
	}

	// Write the status of the repository to the gitView
	gitView.SetText(fmt.Sprintf("The directory %s is a Git repository.\nStatus:\n%s\n", ui.sitePath, status))

	return gitView, nil
}

// CreateResizableWindow creates a resizable window
func (ui *UI) CreateResizableWindow(title string, content tview.Primitive, wm *winman.Manager) { // Add this function
	// Define the window variable
	window := wm.NewWindow()

	// Set the window properties
	window.Show().
		SetRoot(content).
		SetDraggable(true).
		SetResizable(true).
		SetTitle(title)

	// Set the position and size of the window based on its title
	switch title {
	case "Blog Posts":
		window.SetRect(0, 0, 40, 20) // Larger window at the top left corner
	case "Content View":
		window.SetRect(40, 0, 40, 20) // Larger window at the top right corner
	default:
		window.SetRect(0, 20, 80, 10) // Smaller window at the bottom
	}
}

// CreateStatusBar creates a new text view for the status bar
func (ui *UI) CreateStatusBar() *tview.TextView {
	statusBar := tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Status Bar")
	return statusBar
}

// CreatePublishModal creates a new modal for publishing drafts
func (ui *UI) CreatePublishModal(node *tview.TreeNode, doneFunc func(int, string)) *tview.Modal {
	// Get the name of the draft
	draftName := node.GetText()

	return tview.NewModal().
		SetText(fmt.Sprintf("Do you want to publish the draft '%s'?", draftName)).
		AddButtons([]string{"Publish", "Cancel"}).
		SetDoneFunc(doneFunc)
}

// UpdateUI updates the UI after a draft is published
func (ui *UI) UpdateUI(menu *tview.TreeView, node *tview.TreeNode, newFilename string) {
	pathNodes := menu.GetPath(node)
	pathNodes[1].RemoveChild(node)
	postsNode := pathNodes[0].GetChildren()[1]
	node.SetText(newFilename) // Update the node text with the new filename
	postsNode.AddChild(node)
	menu.SetCurrentNode(node)
}
