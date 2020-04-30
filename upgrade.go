package main

import (
	"os/user"
	"path"

	gomu "github.com/hatchify/mod-utils"
	com "github.com/hatchify/mod-utils/com"
	flag "github.com/hatchify/parg"
)

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
		if len(originalBranch) > 0 {
			file.CheckoutBranch(originalBranch)
		}
		return
	}
	tagCommit := string(output)

	output, err = lib.File.CmdOutput("git", "rev-parse", "HEAD")
	if err != nil {
		lib.File.Output("No revision head. Cannot checkout version.")
		if len(originalBranch) > 0 {
			file.CheckoutBranch(originalBranch)
		}
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
			if len(originalBranch) > 0 {
				file.CheckoutBranch(originalBranch)
			}
			return
		}
	}

	file.Output("Installing " + version + "...")

	if err := file.RunCmd("./install.sh", version); err != nil {
		// Try again with permissions
		err = nil
		if err = file.RunCmd("sudo", "./install.sh", version); err != nil {
			file.Output("Failed to install :(")
			if len(originalBranch) > 0 {
				file.CheckoutBranch(originalBranch)
			}
			return err
		}
		// Fix pkg permission issues
		if usr, err := user.Current(); err == nil {
			file.RunCmd("sudo", "chown", "-R", usr.Name, path.Join(usr.HomeDir, "go", "pkg"))
		}
	}

	file.Output("Installed Successfully!")

	if len(originalBranch) > 0 {
		file.CheckoutBranch(originalBranch)
	}

	return
}
