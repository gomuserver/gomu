package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	common "github.com/hatchify/mod-common"
	sort "github.com/hatchify/mod-sort"
	sync "github.com/hatchify/mod-sync"
)

func cleanupStash(libs sort.StringArray) {
	// Resume working directory
	var f common.FileWrapper
	for i := range libs {
		f.Path = libs[i]
		f.StashPop()
	}
}

func main() {
	// Flags/Args
	var (
		action string
		tag    string
		branch string

		filterDeps sort.StringArray
		targetDirs sort.StringArray

		debug   bool
		verbose bool
	)

	// Parse command line values, check supported functions, set defaults
	checkArgs(&action, &branch, &tag, &filterDeps, &targetDirs, &debug, &verbose, &nameOnly)

	// Output Stats
	var stats outputStats

	switch action {
	case "deploy", "sync":
		// Clear mod cache before updating mod files
		sync.CleanModCache()
	}

	// Get all libs within target dirs
	libs := getLibsInAny(targetDirs)
	Println("Scanning", len(libs)+1, "file(s) in", targetDirs)

	// Clean working directory
	var f common.FileWrapper
	for i := range libs {
		f.Path = libs[i]
		// Hide local changes to prevent interference with searching/syncing
		f.Stash()
	}

	// Sort libs
	var fileHead *sort.FileNode
	fileHead, stats.depCount = libs.SortedDependingOnAny(filterDeps)
	if len(filterDeps) == 0 || len(filterDeps[0]) == 0 {
		Println("Performing", action, "on", stats.depCount, "lib(s)")
	} else {
		Println("Performing", action, "on", stats.depCount, "lib(s) depending on", filterDeps)
	}

	switch action {
	case "sync", "deploy":
		if !showWarning("\nIs this ok?") {
			cleanupStash(libs)
			os.Exit(-1)
		}
	default:
		// No worries
	}

	// Perform action on sorted libs
	index := 0
	for itr := fileHead; itr != nil; itr = itr.Next {
		index++

		// If we're just listing files, we don't need to do anything else :)
		if action == "list" {
			Println(strconv.Itoa(index) + ") " + itr.File.GetGoURL())
			continue
		}

		// Separate output
		Println("")
		Println("(", index, "/", stats.depCount, ")", itr.File.Path)

		itr.File.Output("Checking out " + branch + "...")

		if action == "pull" {
			// Only git pull.
			if performPull(branch, itr) {
				itr.File.Updated = true
				stats.updateCount++
				stats.updatedOutput += strconv.Itoa(stats.updateCount) + ") " + itr.File.Path

				if debug { // Verbose?
					popOutput, err := itr.File.CmdOutput("git", "stash", "pop")
					if err == nil {
						stats.updatedOutput += popOutput
					}
				}
				stats.updatedOutput += "\n"
			}

			continue
		}

		if itr.File.CheckoutOrCreateBranch(branch) != nil {
			itr.File.Output("Failed to checkout " + branch + " :(")
		}

		if itr.File.Pull() != nil {
			itr.File.Output("Failed to pull " + branch + " :(")
		}

		// Create sync lib ref from dep file
		var lib sync.Library
		lib.File = itr.File

		if action == "deploy" {
			// TODO: Branch and PR? Diff?
			lib.File.Output("Checking for local changes...")
			lib.File.Deployed = lib.ModDeploy(tag)
			if lib.File.Deployed {
				stats.deployedCount++
				stats.deployedOutput += strconv.Itoa(stats.deployedCount) + ") " + itr.File.Path + "\n"
			}
		}

		// Aggregate updated versions of previously parsed deps
		lib.ModAddDeps(fileHead)

		if action == "replace-local" {
			// Append local replacements for all libs in lib.updatedDeps
			lib.File.Output("Replacing deps with local directories...")
			if lib.SetLocalDeps() {
				lib.File.Updated = true
				stats.updateCount++
				stats.updatedOutput += strconv.Itoa(stats.updateCount) + ") " + lib.File.Path + "\n"
				lib.File.Output("Local replacements set!")
			} else {
				lib.File.Output("Failed to set local deps :(")
			}
			continue
		}

		// Update the dep if necessary
		if err := lib.ModUpdate(branch, "Update mod files. "+tag); err == nil {
			// Dep was updated
			lib.File.Updated = true
			stats.updateCount++
			stats.updatedOutput += strconv.Itoa(stats.updateCount) + ") " + lib.File.Path + "\n"
		}

		if strings.HasSuffix(strings.Trim(itr.File.Path, "/"), "-plugin") {
			// Ignore tagging
			continue
		}

		// Tag if forced or if able to increment
		if len(tag) > 0 || lib.ShouldTag() {
			itr.File.Version = lib.TagLib(tag)
			itr.File.Tagged = true
			stats.tagCount++
			stats.taggedOutput += strconv.Itoa(stats.tagCount) + ") " + lib.File.Path + " " + lib.File.Version + "\n"
		}

		if len(itr.File.Version) == 0 {
			if len(tag) == 0 {
				itr.File.Version = lib.GetCurrentTag()
			} else {
				itr.File.Version = tag
			}
		}

		if action == "install" {
			// Attempt installation
			itr.File.Output("Building...")
			if itr.File.RunCmd("go", "build", "-o", "testartifact.test") != nil {
				itr.File.Output("Nothing to build")
				continue
			}

			if itr.File.RunCmd("rm", "testartifact.test") != nil {
				itr.File.Output("Nothing to install")
				continue
			}

			if itr.File.RunCmd("go", "install", "-trimpath") == nil {
				itr.File.Installed = true
				stats.installedCount++
				stats.installedOutput += strconv.Itoa(stats.installedCount) + ") " + itr.File.GetGoURL()
				itr.File.Output("Install successful!")
			} else {
				itr.File.Output("Install failed :(")
			}
		}
	}

	// Cleanup
	for i := range libs {
		f.Path = libs[i]
		f.CheckoutBranch(branch)
		f.StashPop()
	}

	if nameOnly {
		// Print names and quit
		for fileItr := fileHead; fileItr != nil; fileItr = fileItr.Next {
			if fileItr.File.Tagged || fileItr.File.Deployed || fileItr.File.Updated || fileItr.File.Installed || action == "list" {
				fmt.Println(fileItr.File.GetGoURL())
			}
		}

		return
	}

	// Separator
	Println("")

	if action == "list" {
		// If we're just listing files, we don't need to do anything else :)
		return
	}

	Println(stats.format(action, branch))
}
