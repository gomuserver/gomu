package main

import (
	"flag"
	"os/exec"
	"path"
	"strings"

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
		libs[index] = path.Join(dir, libs[index])
	}

	return
}
