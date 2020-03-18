package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	common "github.com/hatchify/mod-common"
	sorter "github.com/hatchify/mod-sort"
)

// Parses arguments to load target directories
// Returns current lib if no args provided
func getTargetDirs() (targetLibs sorter.StringArray) {
	targetLibs = flag.Args()
	if len(targetLibs) == 0 {
		targetLibs = append(targetLibs, ".")
	}
	return
}

// Aggregates all libs within all target dirs
func getLibsInAny(targetDirs []string) (libs sorter.StringArray) {
	libs = make(sorter.StringArray, 0)
	for index := range targetDirs {
		libs = append(libs, getLibsInDirectory(targetDirs[index])...)
	}

	return
}

// Gets all libs in a given directory
func getLibsInDirectory(dir string) (libs sorter.StringArray) {
	cmd := exec.Command("ls")
	cmd.Dir = dir
	stdout, err := cmd.Output()

	if err != nil {
		return
	}

	// Parse files from exec "ls"
	libs = strings.Split(string(stdout), "\n")
	for index := range libs {
		libs[index] = path.Join(dir, libs[index])
	}

	return
}

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

func exit(status int) {
	fmt.Println("\nUsage: gomu <flags> <command: [list|sync|deploy|pull]> | gomu -action=<command [list|sync|deploy]> <other flags>")
	fmt.Println("\nNote: command must be a single token set by action, or trailing optional falgs")
	fmt.Println("\nView README.md @ https://github.com/hatchify/gomu")
	fmt.Println("")
	os.Exit(status)
}

func checkArgs(action, branch, tag *string, filterDeps, targetDirs *sorter.StringArray, debug *bool) {
	// Get optional args for forcing a tag number and filtering target deps
	flag.StringVar(action, "action", "", "function to perform [list|sync|deploy|pull]")
	flag.StringVar(branch, "branch", "master", "branch to user when pull (and eventually pull request) are used. Default to master (eventually default to current)")
	flag.StringVar(tag, "tag", "", "optional value to set for git tag")
	flag.BoolVar(debug, "debug", false, "optional value to get debug output")
	flag.Var(filterDeps, "dep", "optional dependency filter: accepts multiple -dep flags to only list/sort libs which depend on one of the provided filters")
	flag.Var(targetDirs, "dir", "optional directory aggregator: accepts multiple -dir flags to aggregate libs in multiple organizations")
	flag.Parse()

	// Set output level (TODO: Log level?)
	common.SetDebug(*debug)

	// Set default return
	command := *action

	// Check for conflict in action vs args to parse command
	if len(flag.Args()) == 1 {
		if len(command) != 0 {
			if command != flag.Arg(0) {
				// Conflict?
				fmt.Println("\nError: Unable to parse action: <" + command + "> from command: <" + flag.Arg(0) + ">")
				exit(1)
			}
		} else {
			command = flag.Arg(0)
		}
	}

	// Check for supported actions
	command = strings.ToLower(command)
	switch command {
	case "sync", "list", "deploy", "pull":
		// Supported actions
	case "help":
		// exit without error
		exit(0)
	default:
		fmt.Println("\nError: Unsupported action: <" + command + ">")
		exit(1)
	}

	// Set defaults if necessary
	if len(*targetDirs) == 0 {
		*targetDirs = append(*targetDirs, ".")
	}

	if len(*filterDeps) == 0 {
		*filterDeps = append(*filterDeps, "")
	}

	*action = command
	return
}
