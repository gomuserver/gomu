package main

import (
	"fmt"
	"strconv"
	"strings"

	common "github.com/hatchify/mod-common"
	sort "github.com/hatchify/mod-sort"
	sync "github.com/hatchify/mod-sync"
)

func main() {
	// Flags/Args
	var (
		action string
		tag    string
		branch string

		filterDeps sort.StringArray
		targetDirs sort.StringArray

		debug bool
	)

	// Parse command line values, check supported functions, set defaults
	checkArgs(&action, &branch, &tag, &filterDeps, &targetDirs, &debug)

	// Get all libs within target dirs
	libs := getLibsInAny(targetDirs)
	fmt.Println("Scanning", len(libs)+1, "file(s) in", targetDirs)

	// Clean working directory
	var f common.FileWrapper
	for i := range libs {
		f.Path = libs[i]
		// Hide local changes to prevent interference with searching/syncing
		f.Stash()
	}

	// Sort libs
	fileHead, depCount := libs.SortedDependingOnAny(filterDeps)
	if len(filterDeps) == 0 {
		fmt.Println("Performing", action, "on", depCount, "lib(s)")
	} else {
		fmt.Println("Performing", action, "on", depCount, "lib(s) depending on", filterDeps)
	}

	// Output Stats
	updateCount := 0
	tagCount := 0
	deployedCount := 0
	updatedOutput := ""
	taggedOutput := ""
	deployedOutput := ""

	// Perform action on sorted libs
	index := 0
	for itr := fileHead; itr != nil; itr = itr.Next {
		index++

		// If we're just listing files, we don't need to do anything else :)
		if action == "list" {
			fmt.Println(string(index) + ") " + itr.File.GetGoURL())
			continue
		}

		// Separate output
		fmt.Println("")
		fmt.Println("(", index, "/", depCount, ")", itr.File.Path)

		itr.File.Output("Checking out " + branch + "...")

		if action == "pull" {
			// Only git pull.
			if itr.File.CheckoutBranch(branch) != nil {
				itr.File.Output("Failed to checkout " + branch + " :(")
			}

			if itr.File.Pull() != nil {
				itr.File.Output("Failed to pull " + branch + " :(")
			}

			updateCount++
			updatedOutput += strconv.Itoa(updateCount) + ") " + itr.File.Path
			popOutput, err := itr.File.CmdOutput("git", "stash", "pop")
			if err == nil {
				updatedOutput += popOutput
			}
			updatedOutput += "\n"
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
				deployedCount++
				deployedOutput += strconv.Itoa(deployedCount) + ") " + itr.File.Path + "\n"
			}
		}

		// Aggregate updated versions of previously parsed deps
		lib.ModAddDeps(fileHead)

		// Update the dep if necessary
		if err := lib.ModUpdate(branch, "Update mod files. "+tag); err == nil {
			// Dep was updated
			lib.File.Updated = true
			updateCount++
			updatedOutput += strconv.Itoa(updateCount) + ") " + lib.File.Path + "\n"
		}

		if strings.HasSuffix(strings.Trim(itr.File.Path, "/"), "-plugin") {
			// Ignore tagging
			continue
		}

		// Tag if forced or if able to increment
		if len(tag) > 0 || lib.ShouldTag() {
			itr.File.Version = lib.TagLib(tag)
			itr.File.Tagged = true
		}

		if len(itr.File.Version) == 0 {
			if len(tag) == 0 {
				itr.File.Version = lib.GetCurrentTag()
			} else {
				itr.File.Version = tag
			}
		}
	}

	// Resume working directory
	for i := range libs {
		f.Path = libs[i]
		f.CheckoutBranch(branch)
		f.StashPop()
	}

	// Count files updated and prepare status output
	for fileItr := fileHead; fileItr != nil; fileItr = fileItr.Next {
		if fileItr.File.Tagged {
			tagCount++
			taggedOutput += strconv.Itoa(tagCount) + ") " + fileItr.File.Path + " " + fileItr.File.Version + "\n"
		}
	}

	// Separator
	fmt.Println("")

	if action == "list" {
		// If we're just listing files, we don't need to do anything else :)
		return
	}

	printStats(action, taggedOutput, updatedOutput, deployedOutput, tagCount, updateCount, deployedCount, depCount)
}
