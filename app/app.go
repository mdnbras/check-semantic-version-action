package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/go-version"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
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
		{
			Name:  "commits-verify",
			Usage: "Verifica se o padrão de commits está de acordo com o padrão conventional commits",
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
					Name:  "merge_request_id",
					Value: "my-merge-request-identifier",
				},
				cli.StringFlag{
					Name:  "gbtoken",
					Value: "personal_access_token",
				},
			},
			Action: commitVerificationPatterns,
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
		fmt.Println("::error file=app.go,line=124::", erro)
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
		fmt.Println("::error file=app.go,line=145::", erro)
		os.Exit(1)
	}

	client := &http.Client{}
	response, erro := client.Do(request)

	if erro != nil {
		fmt.Println("::error file=app.go,line=153::", erro)
		os.Exit(1)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("::error file=app.go,line=160::", erro)
			os.Exit(1)
		}
	}(response.Body)

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("Update realizado com sucesso!")
	}
}

func commitVerificationPatterns(c *cli.Context) {

	owner := c.String("owner")
	repo := c.String("repository")
	accessToken := c.String("gbtoken")
	prNumber := c.Int("merge_request_id")

	ctx := context.Background()

	// Configuração da autenticação usando token de acesso
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Criação do cliente GitHub
	client := github.NewClient(tc)

	// Obtém a lista de commits do Pull Request
	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, prNumber, nil)
	if err != nil {
		fmt.Printf("Erro ao obter commits: %v\n", err)
		os.Exit(1)
	}

	// Itera sobre os commits e imprime as mensagens
	for _, commit := range commits {
		commitDetails, _, err := client.Git.GetCommit(ctx, owner, repo, commit.GetSHA())
		if err != nil {
			fmt.Printf("Erro ao obter detalhes do commit %s: %v\n", commit.GetSHA(), err)
			continue
		}

		//fmt.Printf("Commit: %s\n", commit.GetSHA())
		//fmt.Printf("Message: %s\n", commitDetails.GetMessage())

		regexPattern := "(feat|chore|refactor|style|fix|docs|build|perf|ci|revert)([\\(])([\\#0-9]+)([\\)\\: ]+)(\\W|\\w)+"

		ok, erro := Check(commitDetails.GetMessage(), regexPattern)
		if erro != nil {
			fmt.Printf("Erro match string: %v\n", erro)
		}

		if !ok {
			fmt.Println("::error file=app.go,line=214::", fmt.Sprintf("Commit fora de padrão: %s\n", commitDetails.GetMessage()))
			os.Exit(1)
		}

		//fmt.Println(ok)
		//fmt.Println("---------------------------")
	}
}
