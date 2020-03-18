package main

import (
	"fmt"
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
	updatedOutput := "\n"
	taggedOutput := "\n"
	deployedOutput := "\n"

	// Perform action on sorted libs
	index := 0
	for itr := fileHead; itr != nil; itr = itr.Next {
		index++

		// If we're just listing files, we don't need to do anything else :)
		if action == "list" {
			fmt.Println("(", index, "/", depCount, ")", itr.File.Path)
			continue
		}

		// Separate output
		fmt.Println("")
		fmt.Println("(", index, "/", depCount, ")", itr.File.Path)
		if action == "pull" {
			// Only git pull.
			itr.File.CheckoutBranch(branch)
			itr.File.Pull()
			updateCount++
			updatedOutput += itr.File.Path
			popOutput, err := itr.File.CmdOutput("git", "stash", "pop")
			if err == nil {
				updatedOutput += popOutput
			}
			updatedOutput += "\n"
			continue
		}

		// Create sync lib ref from dep file
		var lib sync.Library
		lib.File = itr.File

		if action == "deploy" {
			// TODO: Branch and PR? Diff?
			lib.File.Deployed = lib.ModDeploy(tag)
			deployedCount++
			deployedOutput += itr.File.Path + "\n"
		}

		// Aggregate updated versions of previously parsed deps
		lib.ModAddDeps(fileHead)

		// Update the dep if necessary
		if err := lib.ModUpdate("Update mod files. " + tag); err == nil {
			// Dep was updated
			lib.File.Updated = true
			updateCount++
			updatedOutput += lib.File.Path
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

	if action != "deploy" {
		// Resume working directory
		for i := range libs {
			f.Path = libs[i]
			f.CheckoutBranch(branch)
			f.StashPop()
		}
	}

	// Count files updated and prepare status output
	for fileItr := fileHead; fileItr != nil; fileItr = fileItr.Next {
		if fileItr.File.Tagged {
			tagCount++
			taggedOutput += fileItr.File.Path + " " + fileItr.File.Version + "\n"
		}
	}

	// Separator
	fmt.Println("")

	if action == "list" {
		// If we're just listing files, we don't need to do anything else :)
		return
	}

	// Print update status
	if updateCount == 0 {
		fmt.Println("All libs already up to date!")
		fmt.Println("")
	} else {
		fmt.Println("Updated", updateCount, "/", depCount, "lib(s):")
		fmt.Println(updatedOutput)
	}

	// Print tag status
	if tagCount == 0 {
		fmt.Println("All tags already up to date!")
		fmt.Println("")
	} else {
		fmt.Println("Tagged", tagCount, "/", depCount, "lib(s):")
		fmt.Println(taggedOutput)
	}

	if action != "deploy" {
		return
	}

	// Print deploys status
	if deployedCount == 0 {
		fmt.Println("No local changes to deploy in", depCount, "libs.")
		fmt.Println("")
	} else {
		fmt.Println("Tagged", tagCount, "/", depCount, "lib(s):")
		fmt.Println(taggedOutput)
	}
}
