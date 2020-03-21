package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	com "github.com/hatchify/mod-common"
	common "github.com/hatchify/mod-common"
	gomu "github.com/hatchify/mod-utils"
)

var version = "undefined"

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

func parseArgs() (options gomu.Options) {
	options.LogLevel = com.NORMAL

	var argV = os.Args
	var argC = len(argV)

	curFlag := ""
	var arg *string
	var gotTrailing = true
	for i := 1; i < argC; i++ {
		arg = &argV[i]

		if strings.HasPrefix(*arg, "-") {
			if !gotTrailing {
				exitWithError("Error: argument expected for flag <" + curFlag + ">")
			}

			// Parse flags
			curFlag = ""
			switch *arg {
			case "-name-only":
				// Set value if boolean flag
				options.LogLevel = com.NAMEONLY
				gotTrailing = true
			default:
				// Set flag if expecting trailing vailues
				curFlag = *arg

				// Waiting on following args
				gotTrailing = false
			}
		} else {
			if i == argC-1 {
				// End of args
				if gotTrailing {
					// Satisfied previous arg
					if len(options.Action) == 0 {
						// We need an action.. this one should do?
						options.Action = *arg
						break
					}
				}
			}

			gotTrailing = true

			// Parse args
			switch curFlag {
			case "-action", "-a":
				if options.Action == "" {
					options.Action = *arg
				} else if options.Action != *arg {
					// Action does not match parsed action
					exitWithError("Error: Unable to parse action <" + *arg + ">, already provided: " + options.Action)
				}
			case "-branch", "-b":
				options.Branch = *arg
				curFlag = ""

			case "-dep", "-depends", "-filter", "-f":
				options.FilterDependencies = append(options.FilterDependencies, *arg)

			case "-dir", "-directory", "-target", "-include", "-i":
				options.TargetDirectories = append(options.TargetDirectories, *arg)

			case "-log", "-level", "-log-level", "-l":
				if options.LogLevel != com.NAMEONLY {
					// Ignore log level if name-only is set
					options.LogLevel = com.LogLevelFrom(*arg)
				}
				curFlag = ""

			case "-tag", "-t":
				options.Tag = *arg

			case "":
				if len(options.Action) == 0 {
					// Comand
					options.Action = *arg
				} else {
					// Arg
					options.Args = append(options.Args, *arg)
				}
			}
		}
	}

	if len(options.TargetDirectories) == 0 {
		options.TargetDirectories = []string{"."}
	}
	options.Action = strings.ToLower(options.Action)

	if len(options.Action) == 0 {
		// Error parsing
		exitWithError("Error: Unable to parse action. No action provided.")
	}

	return
}

func gomuFromArgs() *gomu.MU {
	options := parseArgs()
	common.SetLogLevel(options.LogLevel)

	// TODO: Validate args/flags?

	// Check for supported actions
	switch options.Action {
	case "list", "pull", "reset", "replace-local":
		// Public commands

	case "sync", "deploy":
		// Supported actions. Fall through

	case "version":
		// Print version and exit without error
		fmt.Println(version)
		os.Exit(0)
	case "help":
		// Print help and exit without error
		showHelp()
		os.Exit(0)
	default:
		// Show usage and exit with error
		com.Errorln("\nError: Unsupported action: <" + options.Action + ">")
		showHelp()
		os.Exit(1)
	}

	return gomu.New(options)
}
