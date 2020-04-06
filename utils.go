package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"

	gomu "github.com/hatchify/mod-utils"
	"github.com/hatchify/mod-utils/com"
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

func showHelp(cmd *flag.Command) {
	if cmd == nil {
		fmt.Println(flag.Help())
	} else {
		fmt.Println(cmd.ShowHelp())
	}
}

func exitWithError(message string) {
	com.Errorln(message)
	os.Exit(1)
}

func upgradeGomu(cmd *flag.Command) (err error) {
	var lib gomu.Library
	var file com.FileWrapper
	usr, err := user.Current()
	if err != nil {
		com.Println("gomu :: Unable to get current user dir :(")
		return
	}
	file.Path = path.Join(usr.HomeDir, "go", "src", "github.com", "hatchify", "gomu")
	lib.File = &file

	version := ""
	originalBranch, _ := lib.File.CurrentBranch()
	if len(cmd.Arguments) > 0 {
		// Set version from args
		if val, ok := cmd.Arguments[0].Value.(string); ok {
			version = val
		} else {
			version = cmd.Arguments[0].Name
		}
	} else {
		version = cmd.StringFrom("-branch")
	}

	file.Output("Checking Installation...")
	currentVersion, _ := file.CmdOutput("gomu", "version")

	if len(version) > 0 {
		// Attempt to checkout this version of source
	} else {
		// TODO: Check current repo tag, not latest repo tag
		version = lib.GetCurrentTag()
		if len(currentVersion) > 0 && currentVersion == version {
			var output = ""
			output, err = lib.File.CmdOutput("git", "rev-list", "-n", "1", version)
			if err != nil {
				// No tag set. skip tag
				lib.File.Output("No revision history. Skipping tag.")
				return
			}
			tagCommit := string(output)

			output, err = lib.File.CmdOutput("git", "rev-parse", "HEAD")
			if err != nil {
				// No tag set. skip tag
				lib.File.Output("No revision head. Skipping tag.")
				return
			}
			headCommit := string(output)

			if tagCommit == headCommit {
				file.Output("Version is up to date!")
				return
			}
		}
	}

	msg := version
	if len(msg) == 0 {
		msg = "latest"
	}
	file.Output("Upgrading Installation from " + currentVersion + " to " + version + "...")
	if file.Fetch() != nil {
		file.Output("Failed to update refs :(")
	}

	if len(version) > 0 {
		file.Output("Setting local gomu repo to: " + version + "...")

		if err = file.CheckoutBranch(version); err != nil {
			file.Output("Failed to checkout " + version + " :(")
			return
		}
		file.Pull()

	} else {
		file.Output("Updating source...")
		if file.Pull() != nil {
			file.Output("Failed to update source :(")
		}
	}

	var output = ""
	output, err = lib.File.CmdOutput("git", "rev-list", "-n", "1", version)
	if err != nil {
		// No tag set. skip tag
		lib.File.Output("No revision history. Skipping tag.")
		return
	}
	tagCommit := string(output)

	output, err = lib.File.CmdOutput("git", "rev-parse", "HEAD")
	if err != nil {
		// No tag set. skip tag
		lib.File.Output("No revision head. Skipping tag.")
		return
	}
	headCommit := string(output)

	if file.HasChanges() {
		headCommit = "local"
	}

	if file.HasChanges() || version != lib.GetCurrentTag() {
		version += "-(" + headCommit + ")"
	}

	if currentVersion == version && tagCommit == headCommit {
		if !file.HasChanges() {
			file.Output("Version is up to date!")
			return
		}
	}

	file.Output("Installing " + version + "...")

	if err := file.RunCmd("./install.sh", version); err != nil {
		// Try again with permissions
		err = nil
		if err = file.RunCmd("sudo", "./install.sh", version); err != nil {
			file.Output("Failed to install :(")
			return err
		}
	}

	file.Output("Installed Successfully!")

	if len(originalBranch) > 0 {
		file.CheckoutBranch(originalBranch)
	}

	return
}

// Parg will parse your args
func getCommand() (cmd *flag.Command, err error) {
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

	parg.AddAction("sync", "Updates modfiles\n  Conditionally performs extra tasks depending on flags.\n  Usage: `gomu <flags> sync mod-common parg simply <flags>`")

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

	return flag.Validate()
}

func gomuOptions() (options gomu.Options) {
	// Get command from args
	cmd, err := getCommand()

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
		upgradeGomu(cmd)
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
