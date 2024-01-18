package app

import (
	"encoding/json"
	"fmt"
	"github.com/actions-go/toolkit/core"
	"github.com/google/go-github/v39/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func CommitVerificationPatterns() {

	owner, _ := core.GetInput("owner")
	repo, _ := core.GetInput("repository")
	accessToken, _ := core.GetInput("gbtoken")
	prNumberStr, _ := core.GetInput("mergeRequestId")
	bypass, _ := core.GetInput("bypass")
	urlWebhook, _ := core.GetInput("urlWebhook")

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

	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		fmt.Printf("Erro ao obter Pull Request: %v\n", err)
		os.Exit(1)
	}

	// Obtém a lista de commits do Pull Request
	commits, _, err := client.PullRequests.ListCommits(ctx, owner, repo, prNumber, nil)
	if err != nil {
		fmt.Printf("Erro ao obter commits: %v\n", err)
		os.Exit(1)
	}

	var commitsArr []string

	commitsArr = append(commitsArr, "PR Aberta Por: "+pr.User.GetName())

	existsErro := false

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
			fmt.Println("::error file=github_commit_verification_patterns.go,line=74:: %s", fmt.Sprintf("Commit fora de padrão: %s\n", commitDetails.GetMessage()))
			commitsArr = append(commitsArr, fmt.Sprintf("Commit fora de padrão: **%s**", commitDetails.GetMessage()))
			//os.Exit(1)
			existsErro = true
		}

		//fmt.Println(ok)
		//fmt.Println("---------------------------")
	}

	justString := strings.Join(commitsArr, "\n")

	if urlWebhook != "" && existsErro {
		sendDiscordMessage(urlWebhook, justString)
	}
}

func sendDiscordMessage(urlWebhook string, commitMessage string) {
	method := "POST"

	payload := map[string]interface{}{
		"content": "*Commits Fora do Padrão*\n\n" + commitMessage,
	}

	client := &http.Client{}
	jsonStr, _ := json.Marshal(payload)
	req, err := http.NewRequest(method, urlWebhook, strings.NewReader(string(jsonStr)))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
