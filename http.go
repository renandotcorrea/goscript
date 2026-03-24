package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type HttpRequest struct {
	url     string
	headers map[string]string
	Method  string
	body    []byte
}

type HttpResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

func Get(url string) *HttpRequest {
	return &HttpRequest{url: url, Method: http.MethodGet, headers: make(map[string]string)}
}

func (req *HttpRequest) Headers(headers map[string]string) *HttpRequest {
	for k, v := range headers {
		req.headers[k] = v
	}

	return req
}

func (req *HttpRequest) Body(body []byte) *HttpRequest {
	req.body = body
	return req
}

func (req *HttpRequest) JSON(dest any) error {
	req.headers["Content-Type"] = "application/json"
	req.headers["Accept"] = "application/json"

	res, err := req.do()
	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
		return fmt.Errorf("err status code: %d", res.StatusCode)
	}

	err = json.Unmarshal(res.Body, dest)
	if err != nil {
		return err
	}

	return fmt.Errorf("err status code: %d", res.StatusCode)
}

func (req *HttpRequest) XML(dest any) error {
	req.headers["Content-Type"] = "application/xml"
	req.headers["Accept"] = "application/xml"

	res, err := req.do()
	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated) {
		return fmt.Errorf("err status code: %d", res.StatusCode)
	}

	err = xml.Unmarshal(res.Body, dest)
	if err != nil {
		return err
	}

	return nil
}

func (req *HttpRequest) Do() (*HttpResponse, error) {
	return req.do()
}

func (req *HttpRequest) do() (*HttpResponse, error) {
	request, err := http.NewRequest(req.Method, req.url, bytes.NewBuffer(req.body))
	if err != nil {
		return nil, err
	}

	for k, v := range req.headers {
		request.Header.Set(k, v)
	}

	client := &http.Client{}
	response, err := client.Do(request)
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
