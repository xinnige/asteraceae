package auth0api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	misc "github.com/xinnige/asteraceae/calendula/astermisc"
	utils "github.com/xinnige/asteraceae/calendula/utils"
)

// UserPagination allows for paginating over the users
type UserPagination struct {
	Users  []User
	size   int
	cursor int
	last   int
	remain int
	client *Auth0Client
	err    error
}

func newUserPagination(start, limit int, aclient *Auth0Client) *UserPagination {
	return &UserPagination{
		cursor: start,
		size:   min(limit, max),
		remain: limit,
		last:   -1,
		client: aclient,
	}
}

// Done checks if the pagination has completed
func (p *UserPagination) Done() bool {
	if p == nil || p.err != nil {
		return true
	}
	if p.last == -1 {
		return false
	}
	if p.remain <= 0 {
		return true
	}
	return len(p.Users) == 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Next gets the next page
func (p *UserPagination) Next(ctx context.Context) (*UserPagination, error) {
	var (
		resp []User
	)
	if p.Done() {
		return p, nil
	}

	values := url.Values{
		"page":     {strconv.Itoa(p.cursor)},
		"per_page": {strconv.Itoa(min(p.size, p.remain))},
	}

	endpoint := p.client.Endpoint.URL + "users"
	if err := misc.GetJSON(ctx, p.client.httpClient, endpoint,
		p.client.token, values, &resp, p.client.ParseUsers, p.client); err != nil {
		return nil, err
	}
	p.client.Debugf("ListUsers: %d users (%v)\n", len(resp), p)
	p.Users = resp
	p.remain = p.remain - len(resp)
	p.cursor++
	return p, nil
}

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
	return user.parseUser(client.SerialAPI)
}

func (user *User) parseUser(siface utils.SerialInterface) error {
	if len(user.RawAppMeta) != 0 {
		appmeta := &AuthAppMeta{}
		if err := siface.Unmarshal(user.RawAppMeta, appmeta); err != nil {
			return err
		}
		user.AppMeta = appmeta
	}
	if len(user.RawUserMeta) != 0 {
		usermeta := &SimpleUserMeta{}
		if err := siface.Unmarshal(user.RawUserMeta, usermeta); err != nil {
			return err
		}
		user.UserMeta = usermeta
	}
	return nil
}

// ParseUsers unmarshals a raw json to an array of Users
func (client *Auth0Client) ParseUsers(raw []byte, inf interface{}) error {
	users, ok := inf.(*[]User)
	if !ok {
		return fmt.Errorf("unsupport type %T (expected %T)", inf, &[]User{})
	}
	if err := client.SerialAPI.Unmarshal(raw, users); err != nil {
		return err
	}
	for idx := range *users {
		if err := (*users)[idx].parseUser(client.SerialAPI); err != nil {
			return err
		}
	}
	return nil
}
