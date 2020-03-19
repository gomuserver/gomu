package main

import "strconv"

type outputStats struct {
	depCount int

	updateCount   int
	updatedOutput string

	tagCount     int
	taggedOutput string

	deployedCount  int
	deployedOutput string

	installedCount  int
	installedOutput string
}

func (stats outputStats) format(action, branch string) (output string) {
	if action == "pull" {
		// Print pull status
		output += "Pulled latest version of " + branch + " in " + strconv.Itoa(stats.updateCount) + "/" + strconv.Itoa(stats.depCount) + " lib(s):"
		output += stats.updatedOutput
		return
	}

	if action == "replace-local" {
		// Print replacement status
		output += "Pulled latest version of " + branch + " in " + strconv.Itoa(stats.updateCount) + "/" + strconv.Itoa(stats.depCount) + " lib(s):"
		output += stats.updatedOutput
		return
	}

	// Print update status
	if stats.updateCount == 0 {
		output += "All lib dependencies already up to date!"
		output += ""
	} else {
		output += "Updated mod files in " + strconv.Itoa(stats.updateCount) + "/" + strconv.Itoa(stats.depCount) + " lib(s):"
		output += stats.updatedOutput
	}

	// Print tag status
	if stats.tagCount == 0 {
		output += "All lib tags already up to date!"
		output += ""
	} else {
		output += "Updated tag in " + strconv.Itoa(stats.tagCount) + "/" + strconv.Itoa(stats.depCount) + " lib(s):"
		output += stats.taggedOutput
	}

	if action == "deploy" {
		// Print deploy status
		if stats.deployedCount == 0 {
			output += "No local changes to deploy in " + strconv.Itoa(stats.depCount) + " lib(s)."
			output += ""
		} else {
			output += "Deployed new changes to <" + branch + "> in " + strconv.Itoa(stats.deployedCount) + "/" + strconv.Itoa(stats.depCount) + " lib(s):"
			output += stats.deployedOutput
		}
	} else if action == "install" {
		// Print install status
		if stats.installedCount == 0 {
			output += "No packages installed in " + strconv.Itoa(stats.depCount) + " libs."
			output += ""
		} else {
			output += "Installed " + strconv.Itoa(stats.deployedCount) + "/" + strconv.Itoa(stats.depCount) + " lib(s):"
			output += stats.deployedOutput
		}
	}

	return
}
