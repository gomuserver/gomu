package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func readInput() {
	var (
		err  error
		text string
	)

	files := make([]string, 0)
	reader := bufio.NewReader(os.Stdin)

	// Get files from stdin (piped from another program's output)
	for err == nil {
		if text = strings.TrimSpace(text); len(text) > 0 {
			files = append(files, text)
		}

		text, err = reader.ReadString('\n')
	}

	// Print files
	for i := range files {
		fmt.Println(files[i])
	}
}

func main() {
	flag.Parse()

	// Parse libs
	libs := getLibsInDirectory(".")
	filterDeps := flag.Args()

	// Sort libs, filter if deps provided
	for fileItr := libs.SortedLibsDependingOn(filterDeps); fileItr != nil; fileItr = fileItr.Next {
		// Print files
		fmt.Println(fileItr.Path)
	}
}
