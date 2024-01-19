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

type PullRequest struct {
	Url    string `json:"url"`
	Id     int    `json:"id"`
	Number int    `json:"number"`
	State  string `json:"state"`
	Locked bool   `json:"locked"`
	Title  string `json:"title"`
	User   struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"user"`
}

var owner, _ = core.GetInput("owner")
var repo, _ = core.GetInput("repository")
var accessToken, _ = core.GetInput("gbtoken")
var prNumberStr, _ = core.GetInput("mergeRequestId")
var bypass, _ = core.GetInput("bypass")
var urlWebhook, _ = core.GetInput("urlWebhook")
var regexPattern, _ = core.GetInput("regexPattern")
var prNumber, _ = strconv.Atoi(prNumberStr)

func CommitVerificationPatterns() {

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

	var commitsArr []string

	userLogin, err := GetUserOfPR(owner, repo, prNumber)
	if err != nil {
		return
	}
	commitsArr = append(commitsArr, fmt.Sprintf("PR Aberta Por: %s", userLogin))

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

		ok, erro := Check(commitDetails.GetMessage(), regexPattern)
		if erro != nil {
			fmt.Printf("Erro match string: %v\n", erro)
		}

		if !ok {
			fmt.Printf("::error file=github_commit_verification_patterns.go,line=82:: %s\n", fmt.Sprintf("Commit fora de padrão: %s\n", commitDetails.GetMessage()))
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

func GetUserOfPR(owner string, repository string, prNumber int) (string, error) {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%v", owner, repository, prNumber)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Authorization", "token "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	//fmt.Println(string(body))

	var pullRequest PullRequest

	if err = json.Unmarshal(body, &pullRequest); err != nil {
		fmt.Println(err)
		return "", err
	}

	return pullRequest.User.Login, nil
}
