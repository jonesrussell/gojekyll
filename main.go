package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := os.Args[1]
	printFilenames(sitePath, "_drafts")
	printFilenames(sitePath, "_posts")
}

func printFilenames(sitePath string, dirName string) {
	dirPath := filepath.Join(sitePath, dirName)
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			fmt.Printf("File in %s: %s\n", dirName, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", dirName, err)
	}
}
