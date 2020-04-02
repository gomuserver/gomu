package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	com "github.com/hatchify/mod-common"
	common "github.com/hatchify/mod-common"
	gomu "github.com/hatchify/mod-utils"
	flag "github.com/hatchify/parg"
)

var version = "undefined"
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

func showHelp() {
	fmt.Println("\nUsage: gomu <flags> <command: [list|pull|replace-local]> | gomu <command: [list|pull|replace-local]> <other flags>")
	fmt.Println("\nAction Note: command must be a single token before or after all flags/args")
	fmt.Println("\nFilter/Target Note: accepts multiple -f/-filter/-dep or -t/-target/-dir arguments")
	fmt.Println("\nView README.md @ https://github.com/hatchify/gomu")
	fmt.Println("")
}

func exitWithError(message string) {
	showHelp()
	com.Errorln(message)
	os.Exit(1)
}

// Parg will parse your args
func getCommand() (cmd *flag.Command, err error) {
	// Command/Arg/Flag parser
	parg := flag.New()

	// Configure commands
	parg.AddAction("")        // Prints help
	parg.AddAction("help")    // Prints help
	parg.AddAction("version") // Prints version (if available)

	parg.AddAction("list") // Prints each file in chain
	parg.AddAction("pull") // Pulls latest changes for each file in chain

	parg.AddAction("replace") // Replaces local for each dep in chain
	parg.AddAction("reset")   // Resets mod files for each dep in chain

	parg.AddAction("sync") // Updates mod files for each dep in chain

	// Configure flags
	parg.AddGlobalFlag(flag.Flag{ // Directories to search in
		Name:        "-include",
		Identifiers: []string{"-i", "-in", "-include"},
		Type:        flag.STRINGS,
	})
	parg.AddGlobalFlag(flag.Flag{ // Branch to checkout/create
		Name:        "-branch",
		Identifiers: []string{"-b", "-branch"},
	})
	parg.AddGlobalFlag(flag.Flag{ // Branch to checkout/create
		Name:        "-message",
		Identifiers: []string{"-m", "-msg", "-message"},
	})
	parg.AddGlobalFlag(flag.Flag{ // Minimal output for | chains
		Name:        "-name-only",
		Identifiers: []string{"-name", "-name-only"},
		Type:        flag.BOOL,
	})
	parg.AddGlobalFlag(flag.Flag{ // Commits local changes
		Name:        "-commit",
		Identifiers: []string{"-c", "-commit"},
		Type:        flag.BOOL,
	})
	parg.AddGlobalFlag(flag.Flag{ // Creates pull request if possible
		Name:        "-pull-request",
		Identifiers: []string{"-pr", "-pull-request"},
		Type:        flag.BOOL,
	})
	parg.AddGlobalFlag(flag.Flag{ // Update tag/version for changed libs or subdeps
		Name:        "-tag",
		Identifiers: []string{"-t", "-tag"},
		Type:        flag.BOOL,
	})

	return flag.Validate()
}

func gomuOptions() (options gomu.Options) {
	// Get command from args
	cmd, err := getCommand()
	if err != nil {
		// Show usage and exit with error
		showHelp()
		com.Errorln("\nError parsing arguments: ", err)
		os.Exit(1)
	}
	if cmd == nil {
		showHelp()
		com.Errorln("\nError parsing command: ", err)
		os.Exit(1)
	}

	switch cmd.Action {
	case "version":
		// Print version and exit without error
		fmt.Println(version)
		os.Exit(0)
	case "help", "", " ":
		// Print help and exit without error
		showHelp()
		os.Exit(0)
	}

	// Parse options from cmd
	options.Action = cmd.Action

	// Args
	options.FilterDependencies = make([]string, len(cmd.Arguments))
	for i, argument := range cmd.Arguments {
		options.FilterDependencies[i] = argument.Name
	}

	// Flags
	options.TargetDirectories = cmd.StringsFrom("-include")

	options.Branch = cmd.StringFrom("-branch")
	options.CommitMessage = cmd.StringFrom("-message")

	options.Commit = cmd.BoolFrom("-commit")
	options.PullRequest = cmd.BoolFrom("-pull-request")
	options.Tag = cmd.BoolFrom("-tag")
	nameOnly := cmd.BoolFrom("-name-only")
	if nameOnly {
		options.LogLevel = com.NAMEONLY
	} else {
		options.LogLevel = com.NORMAL
	}

	return
}

func fromArgs() *gomu.MU {
	options := gomuOptions()
	common.SetLogLevel(options.LogLevel)

	if len(options.TargetDirectories) == 0 {
		options.TargetDirectories = []string{"."}
	}

	return gomu.New(options)
}
