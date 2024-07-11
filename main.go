package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := os.Args[1]
	app := tview.NewApplication()

	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	list.AddItem("_drafts", "", 0, nil)
	addFilenamesToList(sitePath, "_drafts", list)
	list.AddItem("_posts", "", 0, nil)
	addFilenamesToList(sitePath, "_posts", list)

	list.SetSelectedFunc(func(i int, mainText string, secondaryText string, shortcut rune) {
		app.Stop()
	})

	if err := app.SetRoot(list, true).Run(); err != nil {
		panic(err)
	}
}

func addFilenamesToList(sitePath string, dirName string, list *tview.List) {
	dirPath := filepath.Join(sitePath, dirName)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			list.AddItem(filepath.Base(path), "", 0, nil)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", dirName, err)
	}
}
