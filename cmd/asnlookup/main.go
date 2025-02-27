package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/viktb/asnlookup"
	"github.com/viktb/asnlookup/pkg/database"
)

const (
	databaseFilenameEnvVar = "ASNLOOKUP_DB"
)

func main() {
	// Initialize flags.
	dbFilename := flag.String(
		"db",
		os.Getenv(databaseFilenameEnvVar),
		fmt.Sprintf("database file to use (env: %s)", databaseFilenameEnvVar),
	)
	batch := flag.Bool(
		"batch",
		false,
		"process IPs from stdin",
	)
	version := flag.Bool(
		"version",
		false,
		"print version information and exit",
	)
	flag.Usage = func() {
		fmt.Printf("Usage: asnlookup [OPTION]... [IP]\n")
		flag.PrintDefaults()
	}

	// Check provided arguments.
	flag.Parse()
	if *version {
		fmt.Printf("asnlookup %s\n", asnlookup.Version)
		os.Exit(0)
	}
	if !*batch && flag.NArg() != 1 {
		fmt.Println("Missing argument: IP")
		flag.Usage()
		os.Exit(1)
	}
	if *dbFilename == "" {
		fmt.Println("Missing required option: db")
		flag.Usage()
		os.Exit(1)
	}

	// Inflate database.
	dbFile, err := os.OpenFile(*dbFilename, os.O_RDONLY, 0)
	if err != nil {
		fmt.Println("Failed to open database file:", err)
		os.Exit(1)
	}
	defer dbFile.Close()
	db, err := database.NewFromDump(dbFile)
	if err != nil {
		fmt.Println("Failed to parse database file:", err)
		os.Exit(1)
	}

	// Do the lookup(s).
	if *batch {
		r := bufio.NewScanner(os.Stdin)
		for r.Scan() {
			ip := net.ParseIP(r.Text())
			lookup(db, &ip)
		}
	} else {
		ip := net.ParseIP(flag.Arg(0))
		lookup(db, &ip)
	}
}

func lookup(db *database.Database, ip *net.IP) {
	as, err := db.Lookup(*ip)
	if errors.Is(err, database.ErrASNotFound) {
		fmt.Println("not found")
		return
	}
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(as.Number)
}
