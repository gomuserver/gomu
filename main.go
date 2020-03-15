package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var (
		err  error
		text string
	)

	files := make([]string, 0)
	reader := bufio.NewReader(os.Stdin)

	// Get files from stdin (piped from another program's output)
	for err == nil {
		text = strings.TrimSpace(text)
		if len(text) > 0 {
			files = append(files, text)
		}
		text, err = reader.ReadString('\n')
	}

	// Print files
	for i := range files {
		fmt.Println(files[i])
	}
}
