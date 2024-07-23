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
	_, err := git.PlainOpen(ui.sitePath)
	if err != nil {
		return nil, fmt.Errorf((fmt.Sprintf("The directory %s is not a Git repository. Consider running 'git init'.\n", ui.sitePath)))
	} else {
		gitView.SetText(fmt.Sprintf("The directory %s is a Git repository.\n", ui.sitePath))
	}

	return gitView, nil
}
