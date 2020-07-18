package main

import (
	"fmt"
	"os"

	gomu "github.com/gomuserver/mod-utils"
	"github.com/gomuserver/mod-utils/com"
	flag "github.com/hatchify/parg"
)

// Parg will parse your args
func configureCommand() (cmd *flag.Command, err error) {
	// Command/Arg/Flag parser
	parg := flag.New()

	// Configure commands
	parg.AddAction("", "Designed to make working with mod files easier.\n  To learn more, run `gomu help` or `gomu help <command>`\n  (Flags can be added to either help command)")
	parg.AddAction("help", "Prints available commands and flags.\n  Use `gomu help <command> <flags>` to get more specific info.")
	parg.AddAction("version", "Prints current version.\n  Install using `gomu upgrade` to get version support.")

	parg.AddAction("list", "Prints each file in dependency chain.")
	parg.AddAction("pull", "Updates branch for file in dependency chain.\n  Providing a -branch will checkout given branch.\n  Creates branch if provided none exists.")

	parg.AddAction("replace", "Replaces each versioned file in the dependency chain.\n  Uses the current checked out local copy.")
	parg.AddAction("reset", "Reverts go.mod and go.sum back to last committed version.\n  Usage: `gomu reset mod-common parg`")
	parg.AddAction("test", "Runs `go test` on each library in the dependency chain.\n  Prints names of failing libraries.\n  Usage: `gomu test mod-common`")

	parg.AddAction("sync", "Updates modfiles.\n  Conditionally performs extra tasks depending on flags.\n  Usage: `gomu <flags> sync mod-common parg simply <flags>`")

	parg.AddAction("workflow", "Adds a github workflow to a repo.\n  Requires -source <template path>.\n  Usage: `gomu workflow mod-utils -c -b new-workflow -source workflows/templates/autotag.yml`")
	//parg.AddAction("secret", "Adds a secret to a repo's github actions.\n  Requires -source <file containing secret>.\n  Usage: `gomu secret mod-utils -source ~/.ssh/server_key.crt`")

	parg.AddAction("upgrade", "Updates gomu itself!\n  Optionally accepts a version number.\n  Without argument, updates to latest tag.\n  Otherwise updates to latest branch/tag provided by first arg or -b.\n  Usage: `gomu upgrade` or `gomu upgrade -b master` or `gomu upgrade v0.5.1`")

	// Configure flags
	parg.AddGlobalFlag(flag.Flag{ // Directories to search in
		Name:        "-include",
		Identifiers: []string{"-i", "-in", "-include"},
		Type:        flag.STRINGS,
		Help:        "Will aggregate files in 1 or more directories.\n  Usage: `gomu list -i hatchify -i vroomy`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Branch to checkout/create
		Name:        "-branch",
		Identifiers: []string{"-b", "-branch"},
		Help:        "Will checkout or create said branch.\n  Updating or creating a pull request.\n  Depending on command and other flags.\n  Usage: `gomu pull -b feature/Jira-Ticket`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Minimal output for | chains
		Name:        "-direct-import",
		Identifiers: []string{"-direct", "-direct-import"},
		Type:        flag.BOOL,
		Help:        "Will avoid recursion in dependency sorting.\n  Only includes deps in go.mod (not go.sum).\n  Usage: `gomu list mod-utils -direct`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Minimal output for | chains
		Name:        "-name-only",
		Identifiers: []string{"-name", "-name-only"},
		Type:        flag.BOOL,
		Help:        "Will reduce output to just the filenames changed.\n  (ls-styled output for | chaining)\n  Usage: `gomu list -name`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Commits local changes
		Name:        "-commit",
		Identifiers: []string{"-c", "-commit"},
		Type:        flag.BOOL,
		Help:        "Will commit local changes if present.\n  Includes all changed files in repository.\n  Usage: `gomu sync -c`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Creates pull request if possible
		Name:        "-pull-request",
		Identifiers: []string{"-pr", "-pull-request"},
		Type:        flag.BOOL,
		Help:        "Will create a pull request if possible.\n  Fails if on master, or if no changes.\n  Usage: `gomu sync -pr`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Branch to checkout/create
		Name:        "-message",
		Identifiers: []string{"-m", "-msg", "-message"},
		Help:        "Will set a custom commit message.\n  Applies to -c and -pr flags.\n  Usage: `gomu sync -c -m \"Update all the things!\"`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Update tag/version for changed libs or subdeps
		Name:        "-tag",
		Identifiers: []string{"-t", "-tag"},
		Type:        flag.BOOL,
		Help:        "Will increment tag if new commits since last tag.\n  Requires tag previously set.\n  Usage: `gomu sync -t`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Update tag/version for changed libs or subdeps
		Name:        "-set-version",
		Identifiers: []string{"-set", "-set-version"},
		Help:        "Can be used with -tag to update sem-ver.\n  Will force tag version for all deps in chain.\n  Usage: `gomu sync -t -set v0.5.0`",
	})
	parg.AddGlobalFlag(flag.Flag{ // Update tag/version for changed libs or subdeps
		Name:        "-source-path",
		Identifiers: []string{"-s", "-source", "-source-path"},
		Help:        "Required for workflow and secret commands.\n  Will provide a source template or secret file.\n  Usage: `gomu workflow mod-utils -source path/to/template.yml`",
	})

	return flag.Validate()
}

func gomuOptions() (options gomu.Options) {
	// Get command from args
	cmd, err := configureCommand()

	if err != nil {
		// Show usage and exit with error
		showHelp(nil)
		com.Errorln("Error parsing arguments: ", err)
		os.Exit(1)
	}
	if cmd == nil {
		showHelp(cmd)
		com.Errorln("Error parsing command: ", err)
		os.Exit(1)
	}

	switch cmd.Action {
	case "version":
		// Print version and exit without error
		fmt.Println(version)
		os.Exit(0)
	case "help", "":
		// Print help and exit without error
		showHelp(cmd)
		os.Exit(0)
	case "upgrade":
		upgrade(cmd)
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
	options.SetVersion = cmd.StringFrom("-set-version")

	options.SourcePath = cmd.StringFrom("-source-path")

	options.DirectImport = cmd.BoolFrom("-direct-import")
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
	com.SetLogLevel(options.LogLevel)

	if len(options.TargetDirectories) == 0 {
		options.TargetDirectories = []string{"."}
	}

	return gomu.New(options)
}
