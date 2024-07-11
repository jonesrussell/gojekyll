package ui

import (
	"github.com/rivo/tview"
)

func CreateDashboard() (*tview.Flex, *tview.List, *tview.List) {
	draftsList := tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true)
	postsList := tview.NewList().ShowSecondaryText(false).SetHighlightFullLine(true)

	drafts := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Drafts"), 2, 1, false).
		AddItem(draftsList, 0, 1, true)

	posts := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Posts"), 2, 1, false).
		AddItem(postsList, 0, 1, false)

	dashboard := tview.NewFlex().
		AddItem(drafts, 0, 1, true).
		AddItem(posts, 0, 1, false)

	return dashboard, draftsList, postsList
}
