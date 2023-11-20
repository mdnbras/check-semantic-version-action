package main

import (
	"checkVersionApplication/app"
	"fmt"
	"os"
)

func main() {
	aplicacao := app.Gerar()
	if err := aplicacao.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
