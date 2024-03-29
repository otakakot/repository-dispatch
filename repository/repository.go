package repository

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v60/github"
)

type CreateGitHubAppsTokenInput struct {
	GitHubAppID         string
	GitHubAppPrivateKey string
}

func CreateGitHubAppsToken(
	ctx context.Context,
	input CreateGitHubAppsTokenInput,
) (*github.InstallationToken, *github.Response, error) {
	now := time.Now()

	payload := jwt.MapClaims{
		"exp": now.Unix() + 60,
		"iat": now.Unix(),
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
		return nil, nil, errors.New("expected exactly one installation")
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

	cp := json.RawMessage(input.ClientPayload)

	return client.Repositories.Dispatch(
		ctx,
		input.RepositoryOwner,
		input.RepositoryName,
		github.DispatchRequestOptions{
			EventType:     input.EventType,
			ClientPayload: &cp,
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
