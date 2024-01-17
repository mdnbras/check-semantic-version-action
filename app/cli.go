package app

//
//import (
//	"github.com/urfave/cli"
//)
//
//func Gerar() *cli.App {
//	app := cli.NewApp()
//	app.Name = "GitOps Manager"
//	app.Usage = "The GitOps Manager Utility is a powerful tool designed to simplify and automate common tasks related to Git repository management and continuous integration operations."
//
//	flags := []cli.Flag{
//		cli.StringFlag{
//			Name:  "versionOld",
//			Value: "v0.0.1",
//		},
//		cli.StringFlag{
//			Name:  "versionNew",
//			Value: "v0.0.2",
//		},
//	}
//
//	app.Commands = []cli.Command{
//		{
//			Name:   "verify",
//			Usage:  "Verifica a versão semantica de um projeto",
//			Flags:  flags,
//			Action: verificarVersao,
//		},
//		{
//			Name:  "update-github-vars",
//			Usage: "Atualiza a versão da variavel em um projeto especifico no github",
//			Flags: []cli.Flag{
//				cli.StringFlag{
//					Name:  "owner",
//					Value: "my-profile",
//				},
//				cli.StringFlag{
//					Name:  "repository",
//					Value: "my-repository",
//				},
//				cli.StringFlag{
//					Name:  "varName",
//					Value: "VERSION",
//				},
//				cli.StringFlag{
//					Name:  "varValue",
//					Value: "v0.0.0",
//				},
//				cli.StringFlag{
//					Name:  "gbtoken",
//					Value: "personal_access_token",
//				},
//			},
//			Action: updateGithubVars,
//		},
//		{
//			Name:  "commits-verify",
//			Usage: "Verifica se o padrão de commits está de acordo com o padrão conventional commits",
//			Flags: []cli.Flag{
//				cli.StringFlag{
//					Name:  "owner",
//					Value: "my-profile",
//				},
//				cli.StringFlag{
//					Name:  "repository",
//					Value: "my-repository",
//				},
//				cli.StringFlag{
//					Name:  "merge_request_id",
//					Value: "my-merge-request-identifier",
//				},
//				cli.StringFlag{
//					Name:  "bypass",
//					Value: "bypass",
//				},
//				cli.StringFlag{
//					Name:  "gbtoken",
//					Value: "personal_access_token",
//				},
//			},
//			Action: CommitVerificationPatterns,
//		},
//	}
//
//	return app
//}
