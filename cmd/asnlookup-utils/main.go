package main

import (
	"fmt"
	"os"

	"github.com/viktb/asnlookup"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "convert":
		executeConvertCommand(args)
	case "version":
		printVersion()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`asnlookup-utils - utilities for asnlookup

USAGE:
   asnlookup-utils command [command flags] [arguments...]

COMMANDS:
	convert		converts an MRT file to an asnlookup database
	version		prints the version
	help		prints this help message
`)
}

func printVersion() {
	fmt.Println(fmt.Sprintf("asnlookup-utils %s", asnlookup.Version))
}
