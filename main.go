package main

import (
	"os"
	"sort"

	"github.com/otakakot/repository-dispatch/repository"
	"github.com/urfave/cli/v2"
)

func main() {
	tokenKeyFlag := &cli.StringFlag{
		Name:     "token",
		Usage:    "GitHub Apps Token",
		Required: false,
		Aliases:  []string{"t"},
	}

	idFlag := &cli.StringFlag{
		Name:     "id",
		Usage:    "GitHub Apps ID",
		Required: false,
		Aliases:  []string{"i"},
	}

	privateKeyFlag := &cli.StringFlag{
		Name:     "private-key",
		Usage:    "GitHub Apps Private Key",
		Required: false,
		Aliases:  []string{"k"},
	}

	repoOwnerFlag := &cli.StringFlag{
		Name:     "repository-owner",
		Usage:    "GitHub Repository Owner",
		Required: true,
		Aliases:  []string{"o"},
	}

	repoNameFlag := &cli.StringFlag{
		Name:     "repository-name",
		Usage:    "GitHub Repository Name",
		Required: true,
		Aliases:  []string{"n"},
	}

	eventTypeFlag := &cli.StringFlag{
		Name:     "event-type",
		Usage:    "GitHub Event Type",
		Required: true,
		Aliases:  []string{"e"},
	}

	clientPayloadFlag := &cli.StringFlag{
		Name:     "client-payload",
		Usage:    "GitHub Client Payload",
		Required: true,
		Aliases:  []string{"p"},
	}

	app := &cli.App{
		Name:        "repository-dispatch",
		Usage:       "Repository Dispatch a GitHub Actions workflow",
		Description: "Dispatch a GitHub Actions workflow for a repository. Please specify github apps token or github apps app id and github apps app private key.", //nolint:lll
		Flags: []cli.Flag{
			tokenKeyFlag,
			idFlag,
			privateKeyFlag,
			repoOwnerFlag,
			repoNameFlag,
			eventTypeFlag,
			clientPayloadFlag,
		},
		Action: func(ctx *cli.Context) error {
			token := ctx.String("token")

			if token == "" {
				res, _, err := repository.CreateGitHubAppsToken(
					ctx.Context,
					repository.CreateGitHubAppsTokenInput{
						GitHubAppID:         ctx.String("app-id"),
						GitHubAppPrivateKey: ctx.String("app-private-key"),
					},
				)
				if err != nil {
					return err
				}

				token = res.GetToken()
			}

			if _, _, err := repository.Dispatch(
				ctx.Context,
				repository.DispatchInput{
					GitHubAppsToken: token,
					RepositoryOwner: ctx.String("repository-owner"),
					RepositoryName:  ctx.String("repository-name"),
					EventType:       ctx.String("event-type"),
					ClientPayload:   ctx.String("client-payload"),
				},
			); err != nil {
				return err
			}

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
