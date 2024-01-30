package main

import (
	"checkVersionApplication/app"
	"fmt"
	"github.com/actions-go/toolkit/core"
)

func main() {
	command, ok := core.GetInput("command")
	if !ok {
		fmt.Println("Invalid command input")
		return
	}

	switch command {
	case "commits-verify":
		app.CommitVerificationPatterns()
	case "version-verify":
		app.VersionVerify()
	case "update-github-vars":
		app.UpdateGithubVars()
	}
}
