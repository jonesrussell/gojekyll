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

	app := cmd.NewApp(filehandler.FileHandler{}, ui.UI{}, logger)
	app.Run(os.Args)
}
