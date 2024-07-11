package main

import (
	"os"

	"jonesrussell/jekyll-publisher/cmd"
	"jonesrussell/jekyll-publisher/filehandler"
	"jonesrussell/jekyll-publisher/ui"
)

func main() {
	app := cmd.NewApp(filehandler.FileHandler{}, ui.UI{})
	app.Run(os.Args)
}
