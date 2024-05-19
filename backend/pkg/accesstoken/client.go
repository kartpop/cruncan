package accesstoken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	contentTypeHeader = "Content-Type"
	authPath          = "/v2/oauth/token"
)

// Token represents the access token response, may vary depending on the oauth provider
type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
	ExpiresAt   time.Time
}

func (t *Token) GetExpiration() string {
	return t.ExpiresIn
}

//go:generate mockgen -source=$GOFILE -package $GOPACKAGE -destination client_mock.go
type Client interface {
	GetToken(ctx context.Context) (*Token, error)
}

type client struct {
	BaseUrl      *url.URL
	ClientID     string
	ClientSecret string
}

func NewClient(clientId, clientSecret, baseUrl string) (*client, error) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}

	return &client{
		BaseUrl:      url,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	}, nil
}

func (c *client) GetToken(ctx context.Context) (*Token, error) {
	fullURL := *c.BaseUrl
	fullURL.Path = path.Join(fullURL.Path, authPath)

	urlValues := url.Values{}
	urlValues.Set("grant_type", "client_credentials")
	urlValues.Set("scope", "general")
	body := strings.NewReader(urlValues.Encode())

	request, err := http.NewRequest("POST", fullURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %v", err)
	}
	request.SetBasicAuth(c.ClientID, c.ClientSecret)
	request.Header.Set(contentTypeHeader, "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while sending request: %v", err)
	}

	defer response.Body.Close()
	responseCode := response.StatusCode
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading response body: %v", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("retriable error, reponse code: %d, response body: %v", responseCode, responseBody)
	}

	var token Token
	if err := json.Unmarshal(responseBody, &token); err != nil {
		return nil, fmt.Errorf("error unmarshaling json token: %v", err)
	}

	seconds, err := strconv.Atoi(token.GetExpiration())
	if err != nil {
		return nil, fmt.Errorf("could not convert expiration time to seconds: %v", err)
	}

	token.ExpiresAt = time.Now().Add(time.Second * time.Duration(seconds))

	return &token, nil
}
