package app

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/urfave/cli"
	"io"
	"net/http"
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
		{
			Name:  "update-github-vars",
			Usage: "Atualiza a versão da variavel em um projeto especifico no github",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "owner",
					Value: "my-profile",
				},
				cli.StringFlag{
					Name:  "repository",
					Value: "my-repository",
				},
				cli.StringFlag{
					Name:  "varName",
					Value: "VERSION",
				},
				cli.StringFlag{
					Name:  "varValue",
					Value: "v0.0.0",
				},
				cli.StringFlag{
					Name:  "gbtoken",
					Value: "personal_access_token",
				},
			},
			Action: updateGithubVars,
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

func updateGithubVars(c *cli.Context) {
	owner := c.String("owner")
	repository := c.String("repository")
	varName := c.String("varName")
	varValue := c.String("varValue")
	gbtoken := c.String("gbtoken")

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/variables/%s", owner, repository, varName)
	userData := []byte(`{"name":"` + varName + `","value":"` + varValue + `"}`)

	request, erro := http.NewRequest("PATCH", apiUrl, bytes.NewBuffer(userData))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gbtoken))

	if erro != nil {
		fmt.Println("::error file=app.go,line=120::", erro)
		os.Exit(1)
	}

	client := &http.Client{}
	response, erro := client.Do(request)

	if erro != nil {
		fmt.Println("::error file=app.go,line=128::", erro)
		os.Exit(1)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("::error file=app.go,line=135::", erro)
			os.Exit(1)
		}
	}(response.Body)

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("Update realizado com sucesso!")
	}
}
