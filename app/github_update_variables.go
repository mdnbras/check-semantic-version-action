package app

import (
	"bytes"
	"fmt"
	"github.com/actions-go/toolkit/core"
	"io"
	"net/http"
	"os"
)

var repository, _ = core.GetInput("repository")
var varName, _ = core.GetInput("varName")
var varValue, _ = core.GetInput("varValue")
var gbtoken, _ = core.GetInput("gbtoken")

func UpdateGithubVars() {

	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/variables/%s", owner, repository, varName)
	userData := []byte(`{"name":"` + varName + `","value":"` + varValue + `"}`)

	request, erro := http.NewRequest("PATCH", apiUrl, bytes.NewBuffer(userData))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gbtoken))

	if erro != nil {
		fmt.Println("::error file=github_update_variables.go,line=28::", erro)
		os.Exit(1)
	}

	client := &http.Client{}
	response, erro := client.Do(request)

	if erro != nil {
		fmt.Println("::error file=github_update_variables.go,line=36::", erro)
		os.Exit(1)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("::error file=github_update_variables.go,line=43::", erro)
			os.Exit(1)
		}
	}(response.Body)

	if response.StatusCode == http.StatusNoContent {
		fmt.Println("::debug::Atualização realizada com sucesso!")
	}
}
