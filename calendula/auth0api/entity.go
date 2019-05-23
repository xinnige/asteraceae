package auth0api

import (
	"encoding/json"
	"fmt"
)

// User defines properties of an auth0 user
type User struct {
	AppMeta     interface{}
	CreatedAt   string          `json:"created_at"`
	DN          string          `json:"dn"`
	Identities  []Identity      `json:"identities"`
	LastIP      string          `json:"last_ip"`
	LastLogin   string          `json:"last_login"`
	LoginsCount int             `json:"logs_count"`
	Name        string          `json:"name"`
	Nickname    string          `json:"nickname"`
	OrgUnits    string          `json:"organizationUnits"`
	RawAppMeta  json.RawMessage `json:"app_metadata"`
	RawUserMeta json.RawMessage `json:"user_metadata"`
	UpdatedAt   string          `json:"updated_at"`
	UserID      string          `json:"user_id"`
	UserMeta    interface{}
}

// Identity defines user identity related info
type Identity struct {
	IsSocial bool   `json:"isSocial"`
	Conn     string `json:"connection"`
	Provider string `json:"provider"`
	UserID   string `json:"user_id"`
}

// SimpleUserMeta defines properties of a basic user metadata
type SimpleUserMeta struct {
	Surname   string `json:"surname,omitempty"`
	Givenname string `json:"givenname,omitempty"`
}

// AuthAppMeta defines properties of an auth app metadata
type AuthAppMeta struct {
	LambdaAuthorizer bool     `json:"lambda_authorizer"`
	Apps             []string `json:"apps"`
}

// ParseUser unmarshals a raw json to a User instance
func (client *Auth0Client) ParseUser(raw []byte, inf interface{}) error {
	user, ok := inf.(*User)
	if !ok {
		return fmt.Errorf("unsupport type %T (expected %T)", inf, &User{})
	}
	if err := client.SerialAPI.Unmarshal(raw, user); err != nil {
		return err
	}
	appmeta := &AuthAppMeta{}
	if err := client.SerialAPI.Unmarshal(user.RawAppMeta, appmeta); err != nil {
		return err
	}
	user.AppMeta = appmeta
	usermeta := &SimpleUserMeta{}
	if err := client.SerialAPI.Unmarshal(user.RawUserMeta, usermeta); err != nil {
		return err
	}
	user.UserMeta = usermeta
	return nil
}
