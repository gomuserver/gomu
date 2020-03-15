package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	sorter "github.com/hatchify/dependency-sorter"
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
	paths := sorter.GetLibsInDirectory(".")
	itr := sorter.PopulateDeps(paths)
	// Print files
	for itr != nil {
		fmt.Println(itr.Path)
		itr = itr.Next
	}
}
