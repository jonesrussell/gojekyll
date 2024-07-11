package ui

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

type UIInterface interface {
	CreateDashboard(repoPath string, drafts []string, posts []string) (*tview.Flex, *tview.TreeView, *tview.TextView, *tview.TextView)
	CreateGitView(repoPath string) *tview.TextView
}

type UI struct {
}

// Ensure Menu implements MenuInterface
var _ UIInterface = &UI{}

// CreateDashboard creates a new tview.Flex that contains two lists titled "Drafts" and "Posts".
func (ui UI) CreateDashboard(repoPath string, drafts []string, posts []string) (*tview.Flex, *tview.TreeView, *tview.TextView, *tview.TextView) {
	gitView := ui.CreateGitView(repoPath)

	// Create a tree for the menu
	menu := tview.NewTreeView()

	// Create root for the tree
	root := tview.NewTreeNode("")

	// Add drafts and posts to the tree
	draftsNode := tview.NewTreeNode("Drafts")
	for _, draft := range drafts {
		draftsNode.AddChild(tview.NewTreeNode(draft))
	}
	root.AddChild(draftsNode)

	postsNode := tview.NewTreeNode("Posts")
	for _, post := range posts {
		postsNode.AddChild(tview.NewTreeNode(post))
	}
	root.AddChild(postsNode)

	// Add an "Exit" option to the tree
	exitNode := tview.NewTreeNode("Exit")
	root.AddChild(exitNode)

	menu.SetRoot(root).SetCurrentNode(root)

	// Create a text view for the content of the selected draft or post
	contentView := tview.NewTextView()

	dashboard := tview.NewFlex().
		AddItem(gitView, 0, 1, false).
		AddItem(menu, 0, 1, true).
		AddItem(contentView, 0, 1, false)

	return dashboard, menu, contentView, gitView
}

func (ui UI) CreateGitView(repoPath string) *tview.TextView {
	gitView := tview.NewTextView()

	// Check if the repoPath is a Git repository
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		gitView.SetText(fmt.Sprintf("The directory %s is not a Git repository. Consider running 'git init'.\n", repoPath))
	} else {
		gitView.SetText(fmt.Sprintf("The directory %s is a Git repository.\n", repoPath))
	}

	return gitView
}
