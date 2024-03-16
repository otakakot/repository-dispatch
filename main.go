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
		Required: true,
		Aliases:  []string{"i"},
	}

	secretKeyFlag := &cli.StringFlag{
		Name:     "secret-key",
		Usage:    "GitHub App Secret Key",
		Required: true,
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
		Name:  "repository-dispatch",
		Usage: "Repository Dispatch a GitHub Actions workflow",
		Flags: []cli.Flag{
			appIDFlag,
			secretKeyFlag,
			repoOwnerFlag,
			repoNameFlag,
			eventTypeFlag,
			clientPayloadFlag,
		},
		Action: func(c *cli.Context) error {
			if _, _, err := repository.Dispatch(
				c.Context,
				repository.Input{
					GitHubAppID:     c.String("app-id"),
					GitHubSecretKey: c.String("secret-key"),
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
