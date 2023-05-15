package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "asnlookup-utils",
		Usage: "utilities for asnlookup",
		Commands: []*cli.Command{
			convertCommand,
			versionCommand,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
