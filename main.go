package main

import (
	"os"
	"sort"

	"github.com/urfave/cli/v2"

	"github.com/otakakot/repository-dispatch/repository"
)

func main() {
	appIDFlag := &cli.StringFlag{
		Name:     "app-id",
		Usage:    "GitHub App ID",
		Required: false,
		Aliases:  []string{"i"},
	}

	secretKeyFlag := &cli.StringFlag{
		Name:     "secret-key",
		Usage:    "GitHub App Secret Key",
		Required: false,
		Aliases:  []string{"k"},
	}

	tokenKeyFlag := &cli.StringFlag{
		Name:     "token",
		Usage:    "GitHub Apps Token",
		Required: false,
		Aliases:  []string{"t"},
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
		Description: "Dispatch a GitHub Actions workflow for a repository. Please specify token or app-id and secret-key.",
		Flags: []cli.Flag{
			tokenKeyFlag,
			appIDFlag,
			secretKeyFlag,
			repoOwnerFlag,
			repoNameFlag,
			eventTypeFlag,
			clientPayloadFlag,
		},
		Action: func(c *cli.Context) error {
			token := c.String("token")

			if token == "" {
				res, _, err := repository.CreateGitHubAppsToken(
					c.Context,
					repository.CreateGitHubAppsTokenInput{
						GitHubAppID:     c.String("app-id"),
						GitHubSecretKey: c.String("secret-key"),
					},
				)
				if err != nil {
					return err
				}

				token = res.GetToken()
			}

			if _, _, err := repository.Dispatch(
				c.Context,
				repository.DispatchInput{
					GitHubAppsToken: token,
					RepositoryOwner: c.String("repository-owner"),
					RepositoryName:  c.String("repository-name"),
					EventType:       c.String("event-type"),
					ClientPayload:   c.String("client-payload"),
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
