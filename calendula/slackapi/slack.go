package slackapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	misc "github.com/xinnige/asteraceae/calendula/astermisc"
	"github.com/xinnige/asteraceae/calendula/utils"
)

const (
	envDebug = "DEBUG"
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

type marshalFunc func(interface{}) ([]byte, error)

// Client for the slack api.
type Client struct {
	token     string
	client    misc.AsterClient
	debug     bool
	log       misc.Ilogger
	unmarshal misc.SerialFunc
	marshal   marshalFunc
}

// NewClient returns a pointer of slack api client
func NewClient(accessToken string) *Client {
	return &Client{
		token:     accessToken,
		unmarshal: json.Unmarshal,
		marshal:   json.Marshal,
		client:    &http.Client{},
		debug:     utils.GetEnv(envDebug, "false") == "true",
		log:       log.New(os.Stderr, "slackapi", log.LstdFlags|log.Lshortfile),
	}
}

// SetLogger setup logger
func (api *Client) SetLogger(logger misc.Ilogger) {
	if logger == nil {
		return
	}
	api.log = logger
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
