package main

import (
	com "github.com/hatchify/mod-common"
	gomu "github.com/hatchify/mod-utils"
)

func main() {
	// Parse command line values, check supported functions, set defaults
	mu := fromArgs()

	gomu.RunThen(mu, printOutput)
}

func printOutput(mu *gomu.MU) {
	if len(mu.Errors) > 0 {
		if mu.Options.Action != "list" {
			com.Println("")
		}
		com.Println(mu.Stats.Format())
		com.Println("Quitting with errors:\n", mu.Errors)
		com.Println("")
	} else {
		com.Println("All clean!\n ")
		com.Println(mu.Stats.Format())
	}
}
