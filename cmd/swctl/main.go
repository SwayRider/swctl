package main

import (
	// Keep this on top !!
	_ "github.com/swayrider/swctl/internal"
	"github.com/swayrider/swctl/internal/cmd"

	"context"
	"os"

	log "github.com/swayrider/swlib/logger"
)

func main() {
	app := cmd.App

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}
