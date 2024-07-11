package main

import (
	"fmt"
	"os"

	"jonesrussell/jekyll-publisher/cmd"
	"jonesrussell/jekyll-publisher/filehandler"
	"jonesrussell/jekyll-publisher/logger"
	"jonesrussell/jekyll-publisher/ui"
)

func main() {
	logFilePath := "/tmp/gox.log"
	logger, err := logger.NewLogger(logFilePath)
	if err != nil {
		fmt.Println("Error creating logger:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := os.Args[1]
	app := cmd.NewApp(filehandler.FileHandler{}, ui.NewUI(sitePath), logger)
	app.Run(os.Args)
}
