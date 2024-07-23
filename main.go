package main

import (
	"fmt"
	"os"

	"jonesrussell/gojekyll/cmd"
	"jonesrussell/gojekyll/filehandler"
	"jonesrussell/gojekyll/logger"
	"jonesrussell/gojekyll/ui"
)

func main() {
	logFilePath := "/tmp/gox.log"
	appLogger, err := logger.NewLogger(logFilePath)
	if err != nil {
		fmt.Println("Error creating logger:", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to the Jekyll site as an argument.")
		os.Exit(1)
	}

	sitePath := os.Args[1]
	app := cmd.NewApp(filehandler.NewFileHandler(), ui.NewUI(sitePath), appLogger)
	app.Run(os.Args)
}
