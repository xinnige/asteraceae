package slackapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	misc "github.com/xinnige/asteraceae/calendula/astermisc"
)

// ResponseMetadata holds pagination metadata
type ResponseMetadata struct {
	Cursor string `json:"next_cursor"`
}

func (t *ResponseMetadata) initialize() *ResponseMetadata {
	if t != nil {
		return t
	}

	return &ResponseMetadata{}
}

// Client for the slack api.
type Client struct {
	token  string
	client misc.AsterClient
	debug  bool
	log    misc.Ilogger
	method misc.SerialFunc
}

// NewClient returns a pointer of slack api client
func NewClient(accessToken string) *Client {
	return &Client{
		token:  accessToken,
		method: json.Unmarshal,
		client: &http.Client{},
		debug:  false,
		log:    log.New(os.Stderr, "asteraceae/slackapi", log.LstdFlags|log.Lshortfile),
	}
}

// Debugf print a formatted debug line.
func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		api.log.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debugln print a debug line.
func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		api.log.Output(2, fmt.Sprintln(v...))
	}
}

// Debug returns if debug is enabled.
func (api *Client) Debug() bool {
	return api.debug
}

type errorString string

func (t errorString) Error() string {
	return string(t)
}
