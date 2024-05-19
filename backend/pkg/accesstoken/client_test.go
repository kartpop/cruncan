package accesstoken

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

func TestClient(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

type ClientSuite struct {
	suite.Suite
}

func (c *ClientSuite) Test_GetTokenFromApigee_HappyPath() {
	var baseURL string
	clientID := uuid.NewString()
	clientSecret := "client-secret"
	mockedToken := Token{
		AccessToken: "q2h30F39wYDMJMITYWMp3Qkzkiv4",
		TokenType:   "BearerToken",
		ExpiresIn:   "899",
		ExpiresAt:   time.Now().Add(time.Second * time.Duration(899)),
	}

	handlerFunc := func(rw http.ResponseWriter, request *http.Request) {
		c.Require().Equal(authPath, request.URL.Path)

		username, password, ok := request.BasicAuth()
		c.Require().True(ok)
		c.Require().Equal(clientID, username)
		c.Require().Equal(clientSecret, password)

		contentType := request.Header.Get("Content-Type")
		c.Require().Equal("application/x-www-form-urlencoded", contentType)

		c.Require().NoError(request.ParseForm())
		c.Require().Equal("client_credentials", request.FormValue("grant_type"))
		c.Require().Equal("general", request.FormValue("scope"))

		rw.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(mockedToken)
		_, err := rw.Write(jsonResponse)
		c.Require().NoError(err)
	}

	testserver := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testserver.Close()
	baseURL = testserver.URL

	client, err := NewClient(clientID, clientSecret, baseURL)
	c.Require().NoError(err)

	token, err := client.GetToken(context.Background())
	c.Require().NoError(err)
	c.Equal(mockedToken.AccessToken, token.AccessToken)
}

func (c *ClientSuite) Test_GetTokenFromApigee_Failure() {
	var baseURL string
	clientID := uuid.NewString()
	clientSecret := "client-secret"

	handlerFunc := func(rw http.ResponseWriter, request *http.Request) {
		c.Require().Equal(authPath, request.URL.Path)

		responseBody := []byte(`{
				"ErrorCode": "invalid_client",
				"Error": "Client identifier is required"
			}`)

		rw.WriteHeader(http.StatusBadRequest)
		_, err := rw.Write(responseBody)
		c.Require().NoError(err)
	}

	testserver := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testserver.Close()
	baseURL = testserver.URL

	client, err := NewClient(clientID, clientSecret, baseURL)
	c.Require().NoError(err)

	token, err := client.GetToken(context.Background())
	c.Error(err)
	c.Nil(token)
}
