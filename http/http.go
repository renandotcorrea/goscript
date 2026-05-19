// Package http provides a fluent HTTP request builder with chainable methods.
package http

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// HttpRequest represents an HTTP request being built with a fluent chainable API.
// Use the fluent methods (Headers, BodyJSON, QueryParams, etc.) to configure the request,
// and call Do or JSON to execute it.
type HttpRequest struct {
	url     string
	headers map[string]string
	query   map[string]string
	Method  string
	body    []byte
	retries int
	backoff time.Duration
	timeout time.Duration
	client  *http.Client
}

// HttpResponse represents the response from an HTTP request.
// StatusCode is the HTTP status code, Body is the response body, and Headers are the response headers.
type HttpResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// UnmarshalerFunc is a function type that unmarshals data into a value.
// It's used by JSON and XML methods to decode response bodies.
type UnmarshalerFunc func(data []byte, v any) error

// Get creates a new HttpRequest with the GET method and the specified URL.
// Example usage:
//
//	err := http.Get("http://httpbin.org/get").
//		Headers(map[string]string{"Custom-Header": "value"}).
//		JSON(&dest)
func Get(url string) *HttpRequest {
	return newHttpRequest(http.MethodGet, url)
}

// Post creates a new HttpRequest with the POST method and the specified URL.
// Example usage:
//
//	err := http.Post("http://httpbin.org/post").
//		Headers(map[string]string{"Custom-Header": "value"}).
//		BodyJSON(map[string]string{"key": "value"}).
//		JSON(&dest)
func Post(url string) *HttpRequest {
	return newHttpRequest(http.MethodPost, url)
}

// Put creates a new HttpRequest with the PUT method and the specified URL.
// Example usage:
//
//	err := http.Put("http://httpbin.org/put").
//		Headers(map[string]string{"Custom-Header": "value"}).
//		BodyJSON(map[string]string{"key": "value"}).
//		JSON(&dest)
func Put(url string) *HttpRequest {
	return newHttpRequest(http.MethodPut, url)
}

// Patch creates a new HttpRequest with the PATCH method and the specified URL.
// Example usage:
//
//	err := http.Patch("http://httpbin.org/patch").
//		Headers(map[string]string{"Custom-Header": "value"}).
//		BodyJSON(map[string]string{"key": "value"}).
//		JSON(&dest)
func Patch(url string) *HttpRequest {
	return newHttpRequest(http.MethodPatch, url)
}

// Delete creates a new HttpRequest with the DELETE method and the specified URL.
// Example usage:
//
//	err := http.Delete("http://httpbin.org/delete").
//		Headers(map[string]string{"Custom-Header": "value"}).Do()
func Delete(url string) *HttpRequest {
	return newHttpRequest(http.MethodDelete, url)
}

func newHttpRequest(method, url string) *HttpRequest {
	return &HttpRequest{
		url:     url,
		Method:  method,
		headers: make(map[string]string),
		query:   make(map[string]string),
		client:  &http.Client{},
	}
}

// clone returns a deep copy of the HttpRequest, duplicating the headers and query maps
// so that mutations by one goroutine do not affect others.
func (req *HttpRequest) clone() *HttpRequest {
	c := *req
	c.headers = make(map[string]string, len(req.headers))
	for k, v := range req.headers {
		c.headers[k] = v
	}
	c.query = make(map[string]string, len(req.query))
	for k, v := range req.query {
		c.query[k] = v
	}
	return &c
}

// Headers sets the headers for the HttpRequest and returns the modified HttpRequest.
func (req *HttpRequest) Headers(headers map[string]string) *HttpRequest {
	for k, v := range headers {
		req.headers[k] = v
	}

	return req
}

// QueryParams adds query string parameters to the request URL and returns the modified HttpRequest.
func (req *HttpRequest) QueryParams(params map[string]string) *HttpRequest {
	for k, v := range params {
		req.query[k] = v
	}

	return req
}

// Retry configures the request to retry transient responses (429 and 5xx) and transient network errors.
func (req *HttpRequest) Retry(n int, backoff time.Duration) *HttpRequest {
	if n < 0 {
		n = 0
	}

	req.retries = n
	req.backoff = backoff

	return req
}

// Timeout sets a per-request timeout.
func (req *HttpRequest) Timeout(d time.Duration) *HttpRequest {
	req.timeout = d
	return req
}

// Client sets a custom http.Client for the request. If not set, a default http.Client will be used.
func (req *HttpRequest) Client(client *http.Client) *HttpRequest {
	req.client = client
	return req
}

// BodyJSON sets the body of the HttpRequest to the JSON representation of the provided src and returns the modified HttpRequest.
func (req *HttpRequest) BodyJSON(src any) *HttpRequest {
	body, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}

	return req.jsonBodyContentType().setBody(body)
}

// BodyXML sets the body of the HttpRequest to the XML representation of the provided src and returns the modified HttpRequest.
func (req *HttpRequest) BodyXML(src any) *HttpRequest {
	body, err := xml.Marshal(src)
	if err != nil {
		panic(err)
	}

	return req.xmlBodyContentType().setBody(body)
}

// Body sets the body of the HttpRequest to the provided byte slice and returns the modified HttpRequest.
func (req *HttpRequest) Body(body []byte) *HttpRequest {
	return req.setBody(body)
}

func (req *HttpRequest) setBody(body []byte) *HttpRequest {
	req.body = body
	return req
}

func (req *HttpRequest) jsonBodyContentType() *HttpRequest {
	req.headers["Content-Type"] = "application/json"
	req.headers["Accept"] = "application/json"

	return req
}

// JSON executes the HttpRequest and unmarshals the response body into the provided dest using JSON unmarshaling.
// It returns an error if the request fails or if the response status code is not 200 OK or 201 Created.
func (req *HttpRequest) JSON(dest any) error {
	return req.doAndUnmarshal(json.Unmarshal, dest)
}

func (req *HttpRequest) xmlBodyContentType() *HttpRequest {
	req.headers["Content-Type"] = "application/xml"
	req.headers["Accept"] = "application/xml"
	return req
}

// XML executes the HttpRequest and unmarshals the response body into the provided dest using XML unmarshaling.
// It returns an error if the request fails or if the response status code is not 200 OK or 201 Created.
func (req *HttpRequest) XML(dest any) error {
	return req.doAndUnmarshal(xml.Unmarshal, dest)
}

// Do executes the HttpRequest and returns an HttpResponse containing the status code, body, and headers of the response.
// It returns an error if the request fails.
func (req *HttpRequest) Do() (*HttpResponse, error) {
	return req.do()
}

func (req *HttpRequest) do() (*HttpResponse, error) {
	requestURL, err := req.requestURL()
	if err != nil {
		return nil, err
	}

	attempts := req.retries + 1
	for attempt := 0; attempt < attempts; attempt++ {
		response, doErr := req.doAttempt(requestURL)
		if doErr != nil {
			if attempt == attempts-1 || !isRetryableError(doErr) {
				return nil, doErr
			}
		} else {
			if !isRetryableStatus(response.StatusCode) || attempt == attempts-1 {
				return response, nil
			}
		}

		if req.backoff > 0 {
			time.Sleep(req.backoff)
		}
	}

	return nil, fmt.Errorf("unexpected request execution flow")
}

func (req *HttpRequest) doAttempt(requestURL string) (*HttpResponse, error) {
	request, err := http.NewRequest(req.Method, requestURL, bytes.NewBuffer(req.body))
	if err != nil {
		return nil, err
	}

	for k, v := range req.headers {
		request.Header.Set(k, v)
	}

	if req.timeout > 0 {
		req.client.Timeout = req.timeout
	}

	response, err := req.client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		StatusCode: response.StatusCode,
		Body:       b,
		Headers:    req.headers,
	}, nil
}

func (req *HttpRequest) requestURL() (string, error) {
	if len(req.query) == 0 {
		return req.url, nil
	}

	parsed, err := url.Parse(req.url)
	if err != nil {
		return "", err
	}

	values := parsed.Query()
	for k, v := range req.query {
		values.Set(k, v)
	}

	parsed.RawQuery = values.Encode()
	return parsed.String(), nil
}

func isRetryableStatus(statusCode int) bool {
	if statusCode == http.StatusTooManyRequests {
		return true
	}

	return statusCode >= http.StatusInternalServerError && statusCode <= http.StatusNetworkAuthenticationRequired
}

func isRetryableError(err error) bool {
	var netErr net.Error
	if !errors.As(err, &netErr) {
		return false
	}

	return netErr.Timeout()
}

func (req *HttpRequest) doAndUnmarshal(unmarshalerFunc UnmarshalerFunc, dest any) error {
	res, err := req.do()
	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
		return fmt.Errorf("err status code: %d", res.StatusCode)
	}

	if res.Body != nil && len(res.Body) > 0 {
		err = unmarshalerFunc(res.Body, dest)
		if err != nil {
			return err
		}
	}

	return nil
}
