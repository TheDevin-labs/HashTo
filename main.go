package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: ./hashto js=<variable_name> @file=<input.json> <output.js>")
		os.Exit(1)
	}

	var jsVar, inputFile, outputFile string

	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "js=") {
			jsVar = strings.Split(arg, "=")[1]
		} else if strings.HasPrefix(arg, "@file=") {
			inputFile = strings.Split(arg, "=")[1]
		} else if !strings.Contains(arg, "=") {
			outputFile = arg
		}
	}

	if jsVar == "" || inputFile == "" || outputFile == "" {
		fmt.Println("Error: Invalid arguments.")
		fmt.Println("Required format: ./hashto js=hash @file=work.json work.js")
		os.Exit(1)
	}

	jsonBytes, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	compiler := NewCompiler(jsVar)
	jsCode, err := compiler.Compile(string(jsonBytes))
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, []byte(jsCode), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully compiled %s to %s with variable '%s'\n", inputFile, outputFile, jsVar)
}
