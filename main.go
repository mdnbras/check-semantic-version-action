package main

import (
	"checkVersionApplication/app"
	"fmt"
	"github.com/actions-go/toolkit/core"
)

func main() {
	//aplicacao := app.Gerar()
	//if err := aplicacao.Run(os.Args); err != nil {
	//	fmt.Println(err)
	//}

	command, ok := core.GetInput("command")
	if !ok {
		fmt.Println("Invalid command input")
		return
	}

	switch command {
	case "commits-verify":
		app.CommitVerificationPatterns()
	}
}
