package main

import (
	"github.com/banviktor/asnlookup/pkg/database"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var convertCommand = &cli.Command{
	Name:  "convert",
	Usage: "converts an MRT file to an asnlookup database",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "input",
			Aliases:  []string{"i", "in"},
			Required: true,
			Usage:    "input MRT `file`",
		},
		&cli.StringFlag{
			Name:     "output",
			Aliases:  []string{"o", "out"},
			Required: true,
			Usage:    "output destination `file`",
		},
		&cli.IntFlag{
			Name:  "optimization",
			Value: 5,
			Usage: "set optimization `level` (1 - smallest, 8 - fastest)",
		},
	},
	Action: convertAction,
}

func convertAction(ctx *cli.Context) error {
	var err error

	optimization := ctx.Int("optimization")
	if optimization < 1 || optimization > 8 {
		log.Fatalf("Optimization level must be between 1 and 8")
	}

	// Initialize input.
	inFile := os.Stdin
	inFilePath := ctx.String("input")
	if inFilePath != "-" {
		inFile, err = os.OpenFile(inFilePath, os.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("Failed to open input file: %v", err)
		}
		defer inFile.Close()
	}

	// Initialize output.
	outFile := os.Stdout
	outFilePath := ctx.String("output")
	if outFilePath != "-" {
		outFile, err = os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer outFile.Close()
	}

	// Build database.
	builder := database.NewBuilder()
	err = builder.ImportMRT(inFile)
	if err != nil {
		log.Fatalf("Failed to import: %v", err)
	}
	builder.SetFillFactor(optimizationLevelToFillFactor(optimization))
	db, err := builder.Build()
	if err != nil {
		log.Fatalf("Failed to build database: %v", err)
	}

	// Dump optimized database.
	data, err := db.MarshalBinary()
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	_, err = outFile.Write(data)
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	return nil
}

func optimizationLevelToFillFactor(level int) float32 {
	return float32(9-level) * 0.125
}
