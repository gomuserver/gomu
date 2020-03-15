package main

import (
	"bytes"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Hatch1fy/errors"
	sorter "github.com/hatchify/dependency-sorter"
)

func getLibsInDirectory(dir string) (libs sorter.StringArray) {
	cmd := exec.Command("ls")
	cmd.Dir = dir
	stdout, err := cmd.Output()

	if err != nil {
		return
	}

	// Parse files from exec "ls"
	libs = strings.Split(string(stdout), "\n")

	return
}

func gitPull(gitURL string) (resp string, err error) {
	gitpull := exec.Command("git", "pull", "origin")
	gitpull.Dir = getGitDir(gitURL)
	gitpull.Stdin = os.Stdin

	outBuf := bytes.NewBuffer(nil)
	gitpull.Stdout = outBuf

	errBuf := bytes.NewBuffer(nil)
	gitpull.Stderr = errBuf

	if err = gitpull.Run(); err != nil {
		if errBuf.Len() > 0 {
			err = errors.Error(errBuf.String())
		}

		return
	}

	outStr := outBuf.String()
	if strings.Index(outStr, "up to date") > -1 {
		return
	}

	resp = outStr
	return
}

func gitCheckout(gitURL, branch string) (resp string, err error) {
	gitcheckout := exec.Command("git", "checkout", branch)
	gitcheckout.Dir = getGitDir(gitURL)
	gitcheckout.Stdin = os.Stdin

	outBuf := bytes.NewBuffer(nil)
	gitcheckout.Stdout = outBuf

	errBuf := bytes.NewBuffer(nil)
	gitcheckout.Stderr = errBuf

	if err = gitcheckout.Run(); err == nil && errBuf.Len() == 0 {
		resp = outBuf.String()
		return
	}

	errStr := errBuf.String()
	switch {
	case errStr == "":
		return
	case strings.Index(errStr, "Already on") > -1:
		return

	case strings.Index(errStr, "Switched to") > -1:
		resp = errBuf.String()
		return

	default:
		err = errors.Error(errBuf.String())
		return
	}
}

func getGoDir(gitURL string) (goDir string) {
	homeDir := os.Getenv("HOME")
	return path.Join(homeDir, "go", "src", gitURL)
}

func getGitDir(gitURL string) (goDir string) {
	homeDir := os.Getenv("HOME")
	spl := strings.Split(gitURL, "/")

	var parts []string
	parts = append(parts, homeDir)
	parts = append(parts, "go")
	parts = append(parts, "src")

	if len(spl) > 0 {
		// Append host
		parts = append(parts, spl[0])
	}

	if len(spl) > 1 {
		// Append git user
		parts = append(parts, spl[1])
	}

	if len(spl) > 2 {
		// Append repo name
		parts = append(parts, spl[2])
	}

	return path.Join(parts...)
}
