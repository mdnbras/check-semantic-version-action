package app

import (
	"fmt"
	"github.com/actions-go/toolkit/core"
	"github.com/google/go-github/v39/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"os"
	"strconv"
)

func CommitVerificationPatterns() {

	owner, _ := core.GetInput("owner")
	repo, _ := core.GetInput("repository")
	accessToken, _ := core.GetInput("gbtoken")
	prNumberStr, _ := core.GetInput("merge_request_id")
	bypass, _ := core.GetInput("bypass")

	prNumber, _ := strconv.Atoi(prNumberStr)

	if bypass == "YES" {
		return
	}

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

		regexPattern := "(feat|chore|refactor|style|fix|docs|doc|build|perf|ci|revert)([\\(])([\\#0-9]+)([\\)\\: ]+)(\\W|\\w)+"

		ok, erro := Check(commitDetails.GetMessage(), regexPattern)
		if erro != nil {
			fmt.Printf("Erro match string: %v\n", erro)
		}

		if !ok {
			fmt.Println("::error file=github_commit_verification_patterns.go,line=61::", fmt.Sprintf("Commit fora de padrão: %s\n", commitDetails.GetMessage()))
			//os.Exit(1)
		}

		//fmt.Println(ok)
		//fmt.Println("---------------------------")
	}
}
