package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hatchify/mod-utils/com"
	flag "github.com/hatchify/parg"
)

var logLevel = "NORMAL"

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

func showHelp(cmd *flag.Command) {
	if cmd == nil {
		fmt.Println("# gomu :: Usage ::\n" + flag.Help(true))
	} else {
		fmt.Println("# gomu :: Usage ::\n" + cmd.Help(true))
	}
}

func exitWithError(message string) {
	com.Errorln(message)
	os.Exit(1)
}
