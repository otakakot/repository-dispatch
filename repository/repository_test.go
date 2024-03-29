package repository_test

import (
	"context"
	"os"
	"testing"

	"github.com/otakakot/repository-dispatch/repository"
)

func TestDispatch(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	res, _, err := repository.CreateGitHubAppsToken(
		ctx,
		repository.CreateGitHubAppsTokenInput{
			GitHubAppID:         os.Getenv("GITHUB_APP_ID"),
			GitHubAppPrivateKey: os.Getenv("GITHUB_APP_PRIVATE_KEY"),
		},
	)
	if err != nil {
		t.Fatalf("CreateGitHubAppsToken() error = %v", err)
	}

	input := repository.DispatchInput{
		RepositoryOwner: os.Getenv("GITHUB_REPOSITORY_OWNER"),
		RepositoryName:  os.Getenv("GITHUB_REPOSITORY_NAME"),
		EventType:       os.Getenv("EVENT_TYPE"),
		ClientPayload:   `{"ref":"` + os.Getenv("GITHUB_REF") + `", "sha":"` + os.Getenv("GITHUB_SHA") + `"}`,
		GitHubAppsToken: res.GetToken(),
	}

	if _, _, err := repository.Dispatch(ctx, input); err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}

	if _, _, err := repository.Dispatch(ctx, input); err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
}
