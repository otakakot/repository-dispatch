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

type Input struct {
	GitHubAppID     string
	GitHubSecretKey string
	RepositoryOwner string
	RepositoryName  string
	EventType       string
	ClientPayload   string
}

func Dispatch(
	ctx context.Context,
	input Input,
) (*github.Repository, *github.Response, error) {
	now := time.Now()

	payload := jwt.MapClaims{
		"exp": now.Unix() + 60,
		"iat": now.Unix(),
		"iss": input.GitHubAppID,
	}

	secret := []byte(input.GitHubSecretKey)

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
		Transport: &BearerTokenTransport{JWT: jwt},
	})

	installations, _, err := client.Apps.ListInstallations(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	if len(installations) != 1 {
		return nil, nil, errors.New("expected exactly one installation")
	}

	accessToken, _, err := client.Apps.CreateInstallationToken(ctx, installations[0].GetID(), nil)
	if err != nil {
		return nil, nil, err
	}

	client = github.NewClient(&http.Client{
		Transport: &TokenTransport{Token: accessToken.GetToken()},
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

type BearerTokenTransport struct {
	JWT string
}

func (btt *BearerTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+btt.JWT)

	return http.DefaultTransport.RoundTrip(req)
}

type TokenTransport struct {
	Token string
}

func (tt *TokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+tt.Token)

	return http.DefaultTransport.RoundTrip(req)
}
