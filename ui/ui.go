package ui

import (
	"github.com/rivo/tview"
)

type UIInterface interface {
	CreateDashboard() (*tview.Flex, *tview.List, *tview.List)
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
func (ui UI) CreateDashboard() (*tview.Flex, *tview.List, *tview.List) {
	drafts, draftsList := createList("Drafts")
	posts, postsList := createList("Posts")

	dashboard := tview.NewFlex().
		AddItem(drafts, 0, 1, true).
		AddItem(posts, 0, 1, false)

	return dashboard, draftsList, postsList
}
