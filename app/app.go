package app

import (
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"strconv"
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

func transformVersionToInt(version string) (int64, error) {
	version = strings.Replace(version, "v", "", -1)
	version = strings.Replace(version, ".", "", -1)
	v, erro := strconv.ParseInt(version, 10, 64)

	if erro != nil {
		err := errors.New("falha ao converter a versão para inteiro")
		if err != nil {
			return 0, err
		}
		return 0, erro
	}

	return v, nil
}

func verificarVersao(c *cli.Context) {
	versionOld := c.String("versionOld")
	versionNew := c.String("versionNew")

	var v1, v2 int64
	var erro error

	v1, erro = transformVersionToInt(versionOld)
	v2, erro = transformVersionToInt(versionNew)

	if erro != nil {
		fmt.Println("::error file=app.go,line=66::", erro)
		return
	}

	if v1 > v2 {
		fmt.Println("::error file=app.go,line=74::Versão atual deve ser maior do que a anterior")
	}

	if v1 == v2 {
		fmt.Println("::error file=app.go,line=78::Versão atual não pode ser igual a anterior")
	}

}
