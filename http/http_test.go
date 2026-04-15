package http

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

type testPayload struct {
	Message string `json:"message" xml:"message"`
}

type testXMLPayload struct {
	XMLName xml.Name `xml:"root"`
	Value   string   `xml:"value"`
}

func TestGet(t *testing.T) {
	req := Get("http://example.com")
	if req.Method != nethttp.MethodGet {
		t.Fatalf("expected method GET, got %s", req.Method)
	}
	if req.url != "http://example.com" {
		t.Fatalf("unexpected url: got %s", req.url)
	}
}

func TestPost(t *testing.T) {
	req := Post("http://example.com")
	if req.Method != nethttp.MethodPost {
		t.Fatalf("expected method POST, got %s", req.Method)
	}
}

func TestPut(t *testing.T) {
	req := Put("http://example.com")
	if req.Method != nethttp.MethodPut {
		t.Fatalf("expected method PUT, got %s", req.Method)
	}
}

func TestPatch(t *testing.T) {
	req := Patch("http://example.com")
	if req.Method != nethttp.MethodPatch {
		t.Fatalf("expected method PATCH, got %s", req.Method)
	}
}

func TestDelete(t *testing.T) {
	req := Delete("http://example.com")
	if req.Method != nethttp.MethodDelete {
		t.Fatalf("expected method DELETE, got %s", req.Method)
	}
}

func TestHeaders(t *testing.T) {
	req := Get("http://example.com").
		Headers(map[string]string{"X-Foo": "bar", "X-Baz": "qux"})

	if req.headers["X-Foo"] != "bar" {
		t.Fatalf("expected X-Foo: bar, got %s", req.headers["X-Foo"])
	}
	if req.headers["X-Baz"] != "qux" {
		t.Fatalf("expected X-Baz: qux, got %s", req.headers["X-Baz"])
	}
}

func TestHeaders_Merge(t *testing.T) {
	req := Get("http://example.com").
		Headers(map[string]string{"X-A": "1"}).
		Headers(map[string]string{"X-B": "2"})

	if req.headers["X-A"] != "1" || req.headers["X-B"] != "2" {
		t.Fatalf("headers not merged correctly: %v", req.headers)
	}
}

func TestBodyJSON(t *testing.T) {
	src := testPayload{Message: "hello"}
	req := Post("http://example.com").BodyJSON(src)

	want, _ := json.Marshal(src)
	if string(req.body) != string(want) {
		t.Fatalf("unexpected body: got %s, want %s", req.body, want)
	}
	if req.headers["Content-Type"] != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %s", req.headers["Content-Type"])
	}
	if req.headers["Accept"] != "application/json" {
		t.Fatalf("expected Accept application/json, got %s", req.headers["Accept"])
	}
}

func TestBodyJSON_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for unmarshalable value, got none")
		}
	}()
	Post("http://example.com").BodyJSON(make(chan int))
}

func TestBodyXML(t *testing.T) {
	src := testXMLPayload{Value: "hello"}
	req := Post("http://example.com").BodyXML(src)

	want, _ := xml.Marshal(src)
	if string(req.body) != string(want) {
		t.Fatalf("unexpected body: got %s, want %s", req.body, want)
	}
	if req.headers["Content-Type"] != "application/xml" {
		t.Fatalf("expected Content-Type application/xml, got %s", req.headers["Content-Type"])
	}
	if req.headers["Accept"] != "application/xml" {
		t.Fatalf("expected Accept application/xml, got %s", req.headers["Accept"])
	}
}

func TestBody(t *testing.T) {
	raw := []byte("raw body")
	req := Post("http://example.com").Body(raw)
	if string(req.body) != "raw body" {
		t.Fatalf("unexpected body: got %s", req.body)
	}
}

func TestQueryParams(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.URL.Query().Get("page") != "2" || r.URL.Query().Get("size") != "50" {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}

		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	resp, err := Get(server.URL).
		QueryParams(map[string]string{"page": "2", "size": "50"}).
		Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != nethttp.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestRetry_RetriesOn429(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		current := atomic.AddInt32(&attempts, 1)
		if current < 3 {
			w.WriteHeader(nethttp.StatusTooManyRequests)
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, `{"message":"ok"}`)
	}))
	defer server.Close()

	var dest testPayload
	err := Get(server.URL).Retry(2, 0).JSON(&dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}

	if dest.Message != "ok" {
		t.Fatalf("unexpected message: got %s, want ok", dest.Message)
	}
}

func TestRetry_DoesNotRetryOn400(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(nethttp.StatusBadRequest)
	}))
	defer server.Close()

	var dest testPayload
	err := Get(server.URL).Retry(3, 0).JSON(&dest)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if attempts != 1 {
		t.Fatalf("expected 1 attempt for non-retryable status, got %d", attempts)
	}
}

func TestTimeout(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	_, err := Get(server.URL).Timeout(10 * time.Millisecond).Do()
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestDo_Success(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, "response body")
	}))
	defer server.Close()

	resp, err := Get(server.URL).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != nethttp.StatusOK {
		t.Fatalf("unexpected status code: got %d, want 200", resp.StatusCode)
	}
	if string(resp.Body) != "response body" {
		t.Fatalf("unexpected body: got %s, want 'response body'", resp.Body)
	}
}

func TestDo_SendsHeaders(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Header.Get("X-Custom") != "value" {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	resp, err := Get(server.URL).Headers(map[string]string{"X-Custom": "value"}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != nethttp.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDo_SendsBody(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		b, _ := io.ReadAll(r.Body)
		if string(b) != `{"message":"hello"}` {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	resp, err := Post(server.URL).BodyJSON(testPayload{Message: "hello"}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != nethttp.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestDo_InvalidURL(t *testing.T) {
	_, err := Get("://invalid-url").Do()
	if err == nil {
		t.Fatal("expected error for invalid URL, got none")
	}
}

func TestJSON_Success(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, `{"message":"hello"}`)
	}))
	defer server.Close()

	var dest testPayload
	err := Get(server.URL).JSON(&dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Message != "hello" {
		t.Fatalf("unexpected message: got %s, want hello", dest.Message)
	}
}

func TestJSON_StatusCreated(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusCreated)
		fmt.Fprint(w, `{"message":"created"}`)
	}))
	defer server.Close()

	var dest testPayload
	err := Post(server.URL).JSON(&dest)
	if err != nil {
		t.Fatalf("unexpected error for 201 Created: %v", err)
	}
	if dest.Message != "created" {
		t.Fatalf("unexpected message: got %s, want created", dest.Message)
	}
}

func TestJSON_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusInternalServerError)
	}))
	defer server.Close()

	var dest testPayload
	err := Get(server.URL).JSON(&dest)
	if err == nil {
		t.Fatal("expected error for 500 status, got none")
	}
}

func TestJSON_EmptyBody(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	var dest testPayload
	err := Get(server.URL).JSON(&dest)
	if err != nil {
		t.Fatalf("unexpected error for empty body: %v", err)
	}
}

func TestJSON_InvalidBody(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, "not json")
	}))
	defer server.Close()

	var dest testPayload
	err := Get(server.URL).JSON(&dest)
	if err == nil {
		t.Fatal("expected error for invalid JSON body, got none")
	}
}

func TestXML_Success(t *testing.T) {
	payload := testXMLPayload{Value: "world"}
	body, _ := xml.Marshal(payload)

	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(nethttp.StatusOK)
		w.Write(body)
	}))
	defer server.Close()

	var dest testXMLPayload
	err := Get(server.URL).XML(&dest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Value != "world" {
		t.Fatalf("unexpected value: got %s, want world", dest.Value)
	}
}

func TestXML_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusNotFound)
	}))
	defer server.Close()

	var dest testXMLPayload
	err := Get(server.URL).XML(&dest)
	if err == nil {
		t.Fatal("expected error for 404 status, got none")
	}
}

func TestDo_MethodSentCorrectly(t *testing.T) {
	for _, method := range []string{
		nethttp.MethodGet,
		nethttp.MethodPost,
		nethttp.MethodPut,
		nethttp.MethodPatch,
		nethttp.MethodDelete,
	} {
		method := method
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
				if r.Method != method {
					w.WriteHeader(nethttp.StatusMethodNotAllowed)
					return
				}
				w.WriteHeader(nethttp.StatusOK)
			}))
			defer server.Close()

			req := newHttpRequest(method, server.URL)
			resp, err := req.Do()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp.StatusCode != nethttp.StatusOK {
				t.Fatalf("expected 200, got %d for method %s", resp.StatusCode, method)
			}
		})
	}
}

func ExampleGet() {
	// Example demonstrates a fluent chainable request with query params and headers
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.URL.Query().Get("page") == "1" && r.Header.Get("Authorization") == "Bearer token" {
			w.WriteHeader(nethttp.StatusOK)
			fmt.Fprint(w, `{"data":"success"}`)
			return
		}
		w.WriteHeader(nethttp.StatusBadRequest)
	}))
	defer server.Close()

	var result map[string]string
	err := Get(server.URL).
		QueryParams(map[string]string{"page": "1"}).
		Headers(map[string]string{"Authorization": "Bearer token"}).
		JSON(&result)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result["data"])
	// Output: success
}
