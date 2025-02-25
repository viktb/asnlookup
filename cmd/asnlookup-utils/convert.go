package main

import (
	"flag"
	"log"
	"os"

	"github.com/viktb/asnlookup/pkg/database"
)

func executeConvertCommand(args []string) {
	convertCmd := flag.NewFlagSet("convert", flag.ExitOnError)

	inputFile := convertCmd.String("input", "", "input MRT `file` (required)")
	inShort := convertCmd.String("i", "", "input MRT `file` (shorthand)")
	inAlt := convertCmd.String("in", "", "input MRT `file` (shorthand)")

	outputFile := convertCmd.String("output", "", "output destination `file` (required)")
	outShort := convertCmd.String("o", "", "output destination `file` (shorthand)")
	outAlt := convertCmd.String("out", "", "output destination `file` (shorthand)")

	optimization := convertCmd.Int("optimization", 5, "set optimization `level` (1 - smallest, 8 - fastest)")

	if err := convertCmd.Parse(args); err != nil {
		log.Fatal(err)
	}

	// Handle flag aliases - use the first non-empty value
	inputPath := firstNonEmpty(*inputFile, *inShort, *inAlt)
	outputPath := firstNonEmpty(*outputFile, *outShort, *outAlt)

	// Check required flags
	if inputPath == "" {
		log.Fatal("Required flag --input not provided")
	}
	if outputPath == "" {
		log.Fatal("Required flag --output not provided")
	}

	// Execute the convert action
	err := doConvert(inputPath, outputPath, *optimization)
	if err != nil {
		log.Fatal(err)
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func doConvert(inputPath, outputPath string, optimization int) error {
	var err error

	if optimization < 1 || optimization > 8 {
		log.Fatalf("Optimization level must be between 1 and 8")
	}

	// Initialize input.
	inputFile := os.Stdin
	if inputPath != "-" {
		inputFile, err = os.OpenFile(inputPath, os.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("Failed to open input file: %v", err)
		}
		defer inputFile.Close()
	}

	// Initialize output.
	outputFile := os.Stdout
	if outputPath != "-" {
		outputFile, err = os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Failed to create output file: %v", err)
		}
		defer outputFile.Close()
	}

	// Build database.
	builder := database.NewBuilder()
	err = builder.ImportMRT(inputFile)
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
		log.Fatalf("Failed to marshal database: %v", err)
	}
	_, err = outputFile.Write(data)
	if err != nil {
		log.Fatalf("Failed to write database to file: %v", err)
	}

	return nil
}

func optimizationLevelToFillFactor(level int) float32 {
	return float32(9-level) * 0.125
}
