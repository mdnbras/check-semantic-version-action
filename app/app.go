package app

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/urfave/cli"
	"os"
	"strings"
)

func Gerar() *cli.App {
	app := cli.NewApp()
	app.Name = "Check Version Application"
	app.Usage = "Veificador de versão semantica de um projeto"

	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "versionOld",
			Value: "v0.0.1",
		},
		cli.StringFlag{
			Name:  "versionNew",
			Value: "v0.0.2",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "verify",
			Usage:  "Verifica a versão semantica de um projeto",
			Flags:  flags,
			Action: verificarVersao,
		},
	}

	return app
}

func checkVersion(versionOld string, versionNew string) (bool, error) {
	versionOld = strings.Replace(versionOld, "v", "", -1)
	versionNew = strings.Replace(versionNew, "v", "", -1)

	v1, err := version.NewVersion(versionOld)
	if err != nil {
		return false, err
	}

	v2, err := version.NewVersion(versionNew)
	if err != nil {
		return false, err
	}

	if v1.LessThan(v2) {
		return true, nil
	}

	return false, errors.New("versão atual é menor ou igual a versão anterior")
}

func verificarVersao(c *cli.Context) {
	versionOld := c.String("versionOld")
	versionNew := c.String("versionNew")

	_, erro := checkVersion(versionOld, versionNew)

	if erro != nil {
		fmt.Println("::error file=app.go,line=66::", erro)
		os.Exit(1)
	}
}
