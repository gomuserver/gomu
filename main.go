package main

import (
	"fmt"

	com "github.com/hatchify/mod-common"
	gomu "github.com/hatchify/mod-utils"
)

func main() {
	// Parse command line values, check supported functions, set defaults
	mu := checkArgs()

	switch mu.Options.Action {
	case "deploy", "sync":
		// Clear mod cache before updating mod files
		gomu.CleanModCache()
	}

	fmt.Println("Options:", mu.Options)
	gomu.RunThen(mu, printOutput)
}

func printOutput(mu *gomu.MU) {
	if len(mu.Errors) > 0 {
		com.Println(mu.Stats.Format(mu.Options.Action, mu.Options.Branch))
		com.Println("Quitting with errors:\n", mu.Errors)
		com.Println("")
	} else {
		com.Println("All clean!\n")
		com.Println(mu.Stats.Format(mu.Options.Action, mu.Options.Branch))
	}
}
