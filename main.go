package main

import (
	mod "github.com/hatchify/mod-utils"
	com "github.com/hatchify/mod-utils/com"
	"github.com/hatchify/scribe"
)

var out *scribe.Scribe

func main() {
	outW := scribe.NewStdout()
	outW.SetTypePrefix(scribe.TypeNotification, ":: gomu :: ")
	out = scribe.NewWithWriter(outW, "")

	// Parse command line values, check supported functions, set defaults
	gomu := fromArgs()

	gomu.RunThen(printOutput)
}

func printOutput(mu *mod.MU) {
	if len(mu.Errors) > 0 {
		if mu.Options.Action != "list" {
			com.Println("")
		}
		com.Println(mu.Stats.Format())
		com.Println("Quitting with errors:\n", mu.Errors)
		com.Println("")
	} else {
		com.Println("\nAll clean!\n ")
		com.Println(mu.Stats.Format())
	}
}
