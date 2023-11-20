package main

import (
	"checkVersionApplication/app"
	"fmt"
	"os"
)

func main() {
	fmt.Println("CheckVersionApplication running")
	aplicacao := app.Gerar()
	if err := aplicacao.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
