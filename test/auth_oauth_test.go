package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/TrogonStack/TrogonEventStore-Client-Go/trogoneventstore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const (
	oauthTokenEndpoint  = "http://localhost:18080/realms/trogon-eventstore/protocol/openid-connect/token"
	oauthNodeConnString = "esdb://localhost:2116?tls=true&tlscafile=../certs/ca/ca.crt"
)

func fetchOAuthToken(username, password string) (string, error) {
	response, err := (&http.Client{Timeout: 5 * time.Second}).PostForm(oauthTokenEndpoint, url.Values{
		"grant_type": {"password"},
		"client_id":  {"trogon-eventstore-client"},
		"username":   {username},
		"password":   {password},
		"scope":      {"openid"},
	})
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token endpoint returned %d: %s", response.StatusCode, body)
	}

	var token struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &token); err != nil {
		return "", err
	}
	if token.AccessToken == "" {
		return "", fmt.Errorf("token endpoint returned no access token")
	}

	return token.AccessToken, nil
}

func bearerProvider(username, password string) trogoneventstore.CredentialsProvider {
	return func(context.Context) (*trogoneventstore.Credentials, error) {
		token, err := fetchOAuthToken(username, password)
		if err != nil {
			return nil, err
		}
		return &trogoneventstore.Credentials{BearerToken: token}, nil
	}
}

func staticBearer(token string) trogoneventstore.CredentialsProvider {
	return func(context.Context) (*trogoneventstore.Credentials, error) {
		return &trogoneventstore.Credentials{BearerToken: token}, nil
	}
}

func TestOAuthAuthenticationSuite(t *testing.T) {
	if os.Getenv("TROGON_EVENTSTORE_OAUTH_TESTS") != "1" {
		t.Skip("set TROGON_EVENTSTORE_OAUTH_TESTS=1 to run the OAuth integration suite")
	}
	suite.Run(t, new(OAuthAuthenticationSuite))
}

type OAuthAuthenticationSuite struct {
	suite.Suite
}

func (s *OAuthAuthenticationSuite) SetupSuite() {
	token, err := fetchOAuthToken("admin", "changeit")
	s.Require().NoError(err, "OAuth stack must be available")

	client := s.newClient(staticBearer(token))
	defer client.Close()
	if err := s.createSubscription(client, trogoneventstore.PersistentStreamSubscriptionOptions{}); err != nil {
		s.T().Fatalf("OAuth server rejected a valid administrator token: %v", err)
	}
}

func (s *OAuthAuthenticationSuite) newClient(provider trogoneventstore.CredentialsProvider) *trogoneventstore.Client {
	configuration, err := trogoneventstore.ParseConnectionString(oauthNodeConnString)
	s.Require().NoError(err)
	configuration.CredentialsProvider = provider

	client, err := trogoneventstore.NewClient(configuration)
	s.Require().NoError(err)
	return client
}

func (s *OAuthAuthenticationSuite) createSubscription(client *trogoneventstore.Client, options trogoneventstore.PersistentStreamSubscriptionOptions) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.CreatePersistentSubscription(ctx, uuid.NewString(), uuid.NewString(), options)
}

func (s *OAuthAuthenticationSuite) TestProviderTokenSucceeds() {
	client := s.newClient(bearerProvider("admin", "changeit"))
	defer client.Close()
	s.NoError(s.createSubscription(client, trogoneventstore.PersistentStreamSubscriptionOptions{}))
}

func (s *OAuthAuthenticationSuite) TestPerCallTokenSucceeds() {
	token, err := fetchOAuthToken("admin", "changeit")
	s.Require().NoError(err)
	client := s.newClient(nil)
	defer client.Close()

	s.NoError(s.createSubscription(client, trogoneventstore.PersistentStreamSubscriptionOptions{
		Authenticated: &trogoneventstore.Credentials{BearerToken: token},
	}))
}

func (s *OAuthAuthenticationSuite) TestTokenWithoutRequiredRoleIsDenied() {
	client := s.newClient(bearerProvider("noroles", "changeit"))
	defer client.Close()

	err := s.createSubscription(client, trogoneventstore.PersistentStreamSubscriptionOptions{})
	clientError, _ := trogoneventstore.FromError(err)
	s.Require().NotNil(clientError)
	s.True(clientError.IsErrorCode(trogoneventstore.ErrorCodeAccessDenied), "expected access denied, got %v", err)
}

func (s *OAuthAuthenticationSuite) TestProviderErrorIsReturned() {
	client := s.newClient(func(context.Context) (*trogoneventstore.Credentials, error) {
		return nil, fmt.Errorf("token source unavailable")
	})
	defer client.Close()

	err := s.createSubscription(client, trogoneventstore.PersistentStreamSubscriptionOptions{})
	s.ErrorContains(err, "token source unavailable")
}
