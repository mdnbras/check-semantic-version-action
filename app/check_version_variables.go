package app

import (
	"errors"
	"fmt"
	"github.com/actions-go/toolkit/core"
	"github.com/hashicorp/go-version"
	"os"
	"strings"
)

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

var versionOld, _ = core.GetInput("versionOld")
var versionNew, _ = core.GetInput("versionNew")

func VersionVerify() {
	_, erro := checkVersion(versionOld, versionNew)

	if erro != nil {
		fmt.Println("::error file=check_version_variables.go,line=40::", erro)
		os.Exit(1)
	}
}
