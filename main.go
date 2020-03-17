package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	sort "github.com/hatchify/mod-sort"
	sync "github.com/hatchify/mod-sync"
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
	var (
		tag        string
		filterDeps sort.StringArray
	)

	// Get optional args for forcing a tag number and filtering target deps
	flag.StringVar(&tag, "tag", "", "optional value to set for git tag")
	flag.Var(&filterDeps, "filter", "optional value to set for git tag")
	flag.Parse()

	// Get directories to search in
	targetDirs := getTargetDirs()

	// Get all libs within target dirs
	libs := getLibsInAny(targetDirs)
	fmt.Println("Scanning", len(libs)+1, "file(s) in", targetDirs)

	// Sort libs
	fileHead, depCount := libs.SortedLibsDependingOn(filterDeps)
	if len(filterDeps) == 0 {
		fmt.Println("Found", depCount, "lib(s)")
	} else {
		fmt.Println("Found", depCount, "lib(s) depending on", filterDeps)
	}

	// Sort libs, filter if deps provided, list all if no arguments are given
	index := 0
	for fileItr := fileHead; fileItr != nil; fileItr = fileItr.Next {
		index++

		// Separate output
		fmt.Println("")
		fmt.Println("(", index, "/", depCount, ")", fileItr.Path)

		// Update the dep if necessary
		if err := sync.Update(fileItr.Path, "Update mod files. "+tag); err == nil {
			// Dep was updated
			fileItr.Updated = true
		}

		// Tag if forced or if able to increment
		if len(tag) > 0 || sync.ShouldTag(fileItr.Path) {
			// Ignore plugins even when forcing a tag
			if !strings.HasSuffix(strings.Trim(fileItr.Path, "/"), "-plugin") {
				fileItr.Version = sync.TagLib(fileItr.Path, tag)
			}
		}
	}

	// Count files updated and prepare status output
	updateCount := 0
	output := "\n"
	for fileItr := fileHead; fileItr != nil; fileItr = fileItr.Next {
		if fileItr.Updated {
			updateCount++
			output += fileItr.Path + " " + fileItr.Version + "\n"
		}
	}

	// Separator
	fmt.Println("")

	// Print status
	if updateCount == 0 {
		fmt.Println("All libs already up to date!")
		fmt.Println("")
	} else {
		fmt.Println("Updated", updateCount, "/", depCount, "lib(s):")
		fmt.Println(output)
	}
}
