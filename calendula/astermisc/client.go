package astermisc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SerialFunc unmarshals bytes to interface{}
type SerialFunc func([]byte, interface{}) error

// AsterClient defines the minimal interface needed for an http.Client to be implemented.
type AsterClient interface {
	Do(*http.Request) (*http.Response, error)
}

// RateLimitedError defines a rate limit error
type RateLimitedError struct {
	RetryAfter time.Duration
}

// StatusCodeError represents an http response error.
// type httpStatusCode interface { HTTPStatusCode() int } to handle it.
type statusCodeError struct {
	Code   int
	Status string
}

func (t statusCodeError) Error() string {
	return fmt.Sprintf("%d %s.", t.Code, t.Status)
}

func (t statusCodeError) HTTPStatusCode() int {
	return t.Code
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("Slack rate limit exceeded, retry after %s", e.RetryAfter)
}

func doRequest(ctx context.Context, client AsterClient,
	req *http.Request, intf interface{}, method SerialFunc, d debug) error {
	req = req.WithContext(ctx)
	logRequest(req, d)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitedError{time.Duration(retry) * time.Second}
	}

	// it seems to send an HTML body along with 5xx error codes. Don't parse it.
	if resp.StatusCode != http.StatusOK {
		logResponse(resp, d)
		return statusCodeError{Code: resp.StatusCode, Status: resp.Status}
	}

	return parseResponseBody(resp.Body, intf, method, d)
}

// PostJSON sends POST in JSON.
func PostJSON(ctx context.Context, client AsterClient, endpoint, token string, json []byte, intf interface{}, method SerialFunc, d debug) error {
	reqBody := bytes.NewBuffer(json)
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return doRequest(ctx, client, req, intf, method, d)
}

// send a url encoded form.
func sendForm(ctx context.Context, client AsterClient, endpoint string, values url.Values, intf interface{}, method SerialFunc, d debug) error {
	reqBody := strings.NewReader(values.Encode())
	req, err := http.NewRequest("POST", endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return doRequest(ctx, client, req, intf, method, d)
}

func logRequest(req *http.Request, d debug) error {
	if d.Debug() {
		text, err := httputil.DumpRequest(req, true)
		if err != nil {
			return err
		}
		d.Debugln(string(text))
	}

	return nil
}

func logResponse(resp *http.Response, d debug) error {
	if d.Debug() {
		text, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		d.Debugln(string(text))
	}

	return nil
}

func parseResponseBody(body io.ReadCloser, intf interface{}, method SerialFunc, d debug) error {
	response, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	if d.Debug() {
		d.Debugln("parseResponseBody", string(response))
	}

	return method(response, intf)
}

// GetJSON helps to send a GET request in json
func GetJSON(ctx context.Context, client AsterClient, endpoint, token string, values url.Values, intf interface{}, method SerialFunc, d debug) error {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.URL.RawQuery = values.Encode()
	return doRequest(ctx, client, req, intf, method, d)
}
