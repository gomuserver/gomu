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
		switch libs[index] {
		case ".", "..", dir:
			// Ignore non-repositories
		default:
			libs[index] = path.Join(dir, libs[index])
		}
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

func checkArgs(action, branch, tag *string, filterDeps, targetDirs *sorter.StringArray, debug, verbose, nameOnly *bool) {
	// Get optional args for forcing a tag number, setting branches, and passing actions
	flag.StringVar(action, "action", "", "function to perform [list|sync|deploy|pull]")
	flag.StringVar(branch, "branch", "master", "branch to user when pull (and eventually pull request) are used. Default to master (eventually default to current)")
	flag.StringVar(tag, "tag", "", "optional value to set for git tag")

	// Filter/Aggregator
	flag.Var(filterDeps, "dep", "optional dependency filter: accepts multiple -dep flags to only list/sort libs which depend on one of the provided filters")
	flag.Var(targetDirs, "dir", "optional directory aggregator: accepts multiple -dir flags to aggregate libs in multiple organizations")

	// Output flags
	flag.BoolVar(debug, "debug", false, "optional value to get debug output")
	flag.BoolVar(verbose, "verbose", true, "optional value to print progress output")
	flag.BoolVar(nameOnly, "name-only", false, "optional value to minimize output to just the sorted paths of libs that were changed by gomu. Note, this overrides debug and verbose flags")

	// Load flags
	flag.Parse()

	// Set output level (TODO: Log level?)
	common.SetDebug(*debug && !*nameOnly)
	common.SetVerbose(*verbose && !*nameOnly)

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
		// Supported actions. Fall through

	case "version":
		// Print version and exit without error
		fmt.Println(version)
		os.Exit(0)
	case "help":
		// Print help and exit without error
		exit(0)
	default:
		// Show usage and exit with error
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

func performPull(branch string, itr *sorter.FileNode) (success bool) {
	success = true

	if itr.File.CheckoutBranch(branch) != nil {
		itr.File.Output("Failed to checkout " + branch + " :(")
		success = false
	}

	if itr.File.Pull() == nil {
		itr.File.Output("Pull successful!")
	} else {
		itr.File.Output("Failed to pull " + branch + " :(")
		success = false
	}

	return
}

func printStats(action, branch, taggedOutput, updatedOutput, deployedOutput, installedOutput string, tagCount, updateCount, deployedCount, installedCount, depCount int) {
	if action == "pull" {
		// Print pull status
		fmt.Println("Pulled latest version of", branch, "in", updateCount, "/", depCount, "lib(s):")
		fmt.Println(updatedOutput)
		return
	}

	// Print update status
	if updateCount == 0 {
		fmt.Println("All lib dependencies already up to date!")
		fmt.Println("")
	} else {
		fmt.Println("Updated mod files in", updateCount, "/", depCount, "lib(s):")
		fmt.Println(updatedOutput)
	}

	// Print tag status
	if tagCount == 0 {
		fmt.Println("All lib tags already up to date!")
		fmt.Println("")
	} else {
		fmt.Println("Updated tag in", tagCount, "/", depCount, "lib(s):")
		fmt.Println(taggedOutput)
	}

	if action == "deploy" {
		// Print deploy status
		if deployedCount == 0 {
			fmt.Println("No local changes to deploy in", depCount, "libs.")
			fmt.Println("")
		} else {
			fmt.Println("Deployed new changes to", branch, "in", deployedCount, "/", depCount, "lib(s):")
			fmt.Println(deployedOutput)
		}
	} else if action == "install" {
		// Print install status
		if installedCount == 0 {
			fmt.Println("No packages installed in", depCount, "libs.")
			fmt.Println("")
		} else {
			fmt.Println("Installed", deployedCount, "/", depCount, "lib(s):")
			fmt.Println(deployedOutput)
		}
	}
}
