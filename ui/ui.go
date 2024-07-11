package ui

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/rivo/tview"
)

type UIInterface interface {
	CreateDashboard(repoPath string) (*tview.Flex, *tview.List, *tview.List, *tview.TextView)
	CreateGitView(repoPath string) *tview.TextView
}

type UI struct {
}

// Ensure Menu implements MenuInterface
var _ UIInterface = &UI{}

// createList creates a new tview.List with the given title.
func createList(title string) (*tview.Flex, *tview.List) {
	list := tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText(title), 2, 1, false).
		AddItem(list, 0, 1, true)
	return flex, list
}

// CreateDashboard creates a new tview.Flex that contains two lists titled "Drafts" and "Posts".
func (ui UI) CreateDashboard(repoPath string) (*tview.Flex, *tview.List, *tview.List, *tview.TextView) {
	drafts, draftsList := createList("Drafts")
	posts, postsList := createList("Posts")
	gitView := ui.CreateGitView(repoPath)

	dashboard := tview.NewFlex().
		AddItem(gitView, 0, 1, false).
		AddItem(drafts, 0, 1, true).
		AddItem(posts, 0, 1, false)

	return dashboard, draftsList, postsList, gitView
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
