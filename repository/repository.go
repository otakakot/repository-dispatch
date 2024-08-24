package repository

import (
	"cmp"
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v64/github"
)

var ErrExactryOneInstallation = errors.New("expected exactly one installation")

type CreateGitHubAppsTokenInput struct {
	GitHubAppID         string
	GitHubAppPrivateKey string
	Expires             int64
}

func CreateGitHubAppsToken(
	ctx context.Context,
	input CreateGitHubAppsTokenInput,
) (*github.InstallationToken, *github.Response, error) {
	now := time.Now()

	const def = 60

	expires := cmp.Or(input.Expires, def)

	payload := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Unix() + expires,
		"iss": input.GitHubAppID,
	}

	secret := []byte(input.GitHubAppPrivateKey)

	block, _ := pem.Decode(secret)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, payload)

	jwt, err := token.SignedString(key)
	if err != nil {
		return nil, nil, err
	}

	client := github.NewClient(&http.Client{
		Transport: &bearerTokenTransport{JWT: jwt},
	})

	installations, _, err := client.Apps.ListInstallations(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	if len(installations) != 1 {
		return nil, nil, ErrExactryOneInstallation
	}

	return client.Apps.CreateInstallationToken(ctx, installations[0].GetID(), nil)
}

type DispatchInput struct {
	GitHubAppsToken string
	RepositoryOwner string
	RepositoryName  string
	EventType       string
	ClientPayload   string
}

func Dispatch(
	ctx context.Context,
	input DispatchInput,
) (*github.Repository, *github.Response, error) {
	client := github.NewClient(&http.Client{
		Transport: &tokenTransport{Token: input.GitHubAppsToken},
	})

	payload := json.RawMessage(input.ClientPayload)

	return client.Repositories.Dispatch(
		ctx,
		input.RepositoryOwner,
		input.RepositoryName,
		github.DispatchRequestOptions{
			EventType:     input.EventType,
			ClientPayload: &payload,
		},
	)
}

type bearerTokenTransport struct {
	JWT string
}

func (btt *bearerTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+btt.JWT)

	return http.DefaultTransport.RoundTrip(req)
}

type tokenTransport struct {
	Token string
}

func (tt *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+tt.Token)

	return http.DefaultTransport.RoundTrip(req)
}
