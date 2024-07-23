package cmd

import (
	"github.com/rivo/tview"
)

type AppContext struct {
	tviewApp    *tview.Application
	menu        *tview.TreeView
	contentView *tview.TextView
	gitView     tview.Primitive
	dashboard   *tview.Flex
}
