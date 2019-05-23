package auth0api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	misc "github.com/xinnige/asteraceae/calendula/astermisc"
	utils "github.com/xinnige/asteraceae/calendula/utils"
)

const (
	envProvider = "AUTH_PROVIDER"
	envConn     = "AUTH_CONNECTION"

	deProvider = "ad"
	deConn     = "ldap"
)

// Auth0Client defines the properties to access the auth endpoint
type Auth0Client struct {
	Endpoint   *Auth0Endpoint
	httpClient misc.AsterClient
	SerialAPI  utils.SerialInterface
	debug      bool
	log        misc.Ilogger
	token      string
}

// Auth0Endpoint wraps necessary info of auth0 endpoint
type Auth0Endpoint struct {
	URL        string `json:"ur"`
	Provider   string `json:"provider"`
	Connection string `json:"connection"`
	token      string
}

// NewAuth0Client returns a *AuthClient instance
func NewAuth0Client(rawtoken, endpoint string) *Auth0Client {
	return &Auth0Client{
		Endpoint: &Auth0Endpoint{
			URL:        endpoint,
			Provider:   utils.GetEnv(envProvider, deProvider),
			Connection: utils.GetEnv(envConn, deConn),
		},
		httpClient: &http.Client{},
		token:      rawtoken,
		SerialAPI:  &utils.JSONAPI{},
	}
}

// GetUserByName returns a user by unique name
func (client *Auth0Client) GetUserByName(name string) (*User, error) {
	userid := fmt.Sprintf("%s|%s|%s",
		client.Endpoint.Provider, client.Endpoint.Connection, name)
	user := &User{}
	endpoint := client.Endpoint.URL + path.Join("users", userid)
	values := url.Values{}
	if err := misc.GetJSON(context.Background(), client.httpClient, endpoint,
		client.token, values, user, client.ParseUser, client); err != nil {
		return nil, err
	}
	return user, nil
}

// Debugf print a formatted debug line.
func (client *Auth0Client) Debugf(format string, v ...interface{}) {
	if client.debug {
		client.log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debugln print a debug line.
func (client *Auth0Client) Debugln(v ...interface{}) {
	if client.debug {
		client.log.Output(2, fmt.Sprintln(v...))
	}
}

// Debug returns if debug is enabled.
func (client *Auth0Client) Debug() bool {
	return client.debug
}
