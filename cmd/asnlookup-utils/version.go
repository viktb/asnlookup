package main

import (
	"fmt"
	"github.com/banviktor/asnlookup"
	"github.com/urfave/cli/v2"
)

var versionCommand = &cli.Command{
	Name:   "version",
	Usage:  "",
	Action: versionAction,
}

func versionAction(_ *cli.Context) error {
	fmt.Printf("asnlookup-utils v%s\n", asnlookup.Version)
	return nil
}
