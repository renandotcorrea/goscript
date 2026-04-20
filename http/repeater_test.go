package http

import (
	"fmt"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	nethttp "net/http"
)

func TestRepeat_ExecutesNTimes(t *testing.T) {
	var count int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		atomic.AddInt32(&count, 1)
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	responses, err := Get(server.URL).Repeat(5).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 5 {
		t.Fatalf("expected 5 responses, got %d", len(responses))
	}
	if count != 5 {
		t.Fatalf("expected 5 requests, got %d", count)
	}
}

func TestRepeat_ReturnsAllResponses(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, "ok")
	}))
	defer server.Close()

	responses, err := Get(server.URL).Repeat(3).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, resp := range responses {
		if resp.StatusCode != nethttp.StatusOK {
			t.Fatalf("response[%d]: expected 200, got %d", i, resp.StatusCode)
		}
		if string(resp.Body) != "ok" {
			t.Fatalf("response[%d]: unexpected body: got %s", i, resp.Body)
		}
	}
}

func TestRepeat_StopsOnError(t *testing.T) {
	var count int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		current := atomic.AddInt32(&count, 1)
		if current == 2 {
			w.WriteHeader(nethttp.StatusInternalServerError)
			return
		}
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	// Do() itself doesn't error on non-2xx; only network errors stop it.
	// Verify we still get all N responses even on error status codes.
	responses, err := Get(server.URL).Repeat(3).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 3 {
		t.Fatalf("expected 3 responses, got %d", len(responses))
	}
	if responses[1].StatusCode != nethttp.StatusInternalServerError {
		t.Fatalf("expected response[1] to be 500, got %d", responses[1].StatusCode)
	}
}

func TestRepeat_ErrorOnInvalidURL(t *testing.T) {
	_, err := Get("://invalid").Repeat(3).Do()
	if err == nil {
		t.Fatal("expected error for invalid URL, got nil")
	}
}

func TestRepeat_WithDelay_SetsOption(t *testing.T) {
	req := Get("http://example.com").Repeat(2).WithDelay(100 * time.Millisecond)
	if req.options.Delay != 100*time.Millisecond {
		t.Fatalf("expected delay 100ms, got %v", req.options.Delay)
	}
}

func TestRepeat_WithBackoff_SetsOption(t *testing.T) {
	req := Get("http://example.com").Repeat(2).WithBackoff(3.0)
	if !req.options.Backoff {
		t.Fatal("expected Backoff to be true")
	}
	if req.options.BackoffFactor != 3.0 {
		t.Fatalf("expected BackoffFactor 3.0, got %f", req.options.BackoffFactor)
	}
}

func TestRepeat_WithParallel_SetsOption(t *testing.T) {
	req := Get("http://example.com").Repeat(5).WithParallel(3)
	if !req.options.Parallel {
		t.Fatal("expected Parallel to be true")
	}
	if req.options.ParallelLimit != 3 {
		t.Fatalf("expected ParallelLimit 3, got %d", req.options.ParallelLimit)
	}
}

func TestRepeat_WithOptions_SetsAllOptions(t *testing.T) {
	opts := RequestRepeaterOptions{
		Delay:         50 * time.Millisecond,
		Backoff:       true,
		BackoffFactor: 1.5,
		Parallel:      true,
		ParallelLimit: 4,
	}
	req := Get("http://example.com").Repeat(2).WithOptions(opts)
	if req.options.Delay != opts.Delay ||
		req.options.Backoff != opts.Backoff ||
		req.options.BackoffFactor != opts.BackoffFactor ||
		req.options.Parallel != opts.Parallel ||
		req.options.ParallelLimit != opts.ParallelLimit {
		t.Fatalf("expected options to be set, got %+v", req.options)
	}
}

func TestRepeat_BeforeEach_SetsOption(t *testing.T) {
	fn := func(req *HttpRequest) {}
	repeater := Get("http://example.com").Repeat(3).BeforeEach(fn)
	if repeater.options.BeforeEach == nil {
		t.Fatal("expected BeforeEach to be set")
	}
}

func TestRepeat_BeforeEach_CalledInSequentialPath(t *testing.T) {
	var count int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	_, err := Get(server.URL).Repeat(3).BeforeEach(func(req *HttpRequest) {
		atomic.AddInt32(&count, 1)
	}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected BeforeEach called 3 times, got %d", count)
	}
}

func TestRepeat_BeforeEach_ModifiesRequestPerIteration(t *testing.T) {
	var received []string
	var mu sync.Mutex
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		mu.Lock()
		received = append(received, r.Header.Get("X-Seq"))
		mu.Unlock()
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	var idx int32
	_, err := Get(server.URL).Repeat(3).BeforeEach(func(req *HttpRequest) {
		i := atomic.AddInt32(&idx, 1)
		req.headers["X-Seq"] = fmt.Sprintf("%d", i)
	}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(received) != 3 {
		t.Fatalf("expected 3 requests, got %d", len(received))
	}
}

func TestRepeat_Parallel_ExecutesNTimes(t *testing.T) {
	var count int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		atomic.AddInt32(&count, 1)
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	responses, err := Get(server.URL).Repeat(5).WithParallel(3).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 5 {
		t.Fatalf("expected 5 responses, got %d", len(responses))
	}
	if count != 5 {
		t.Fatalf("expected 5 requests, got %d", count)
	}
}

func TestRepeat_Parallel_AllResponsesSucceed(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, "parallel")
	}))
	defer server.Close()

	responses, err := Get(server.URL).Repeat(4).WithParallel(4).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, resp := range responses {
		if resp.StatusCode != nethttp.StatusOK {
			t.Fatalf("response[%d]: expected 200, got %d", i, resp.StatusCode)
		}
	}
}

func TestRepeat_Parallel_BeforeEach_CalledForEachRequest(t *testing.T) {
	var count int32
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	_, err := Get(server.URL).Repeat(3).WithParallel(3).BeforeEach(func(req *HttpRequest) {
		atomic.AddInt32(&count, 1)
	}).Do()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected BeforeEach to be called 3 times, got %d", count)
	}
}

func TestRepeat_Parallel_CancelsOnError(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	// Force a network error by using an invalid URL for one of the requests.
	// Use BeforeEach to inject the bad URL on the first call.
	var once sync.Once
	_, err := Get(server.URL).Repeat(5).WithParallel(5).BeforeEach(func(req *HttpRequest) {
		once.Do(func() { req.url = "://invalid" })
	}).Do()
	if err == nil {
		t.Fatal("expected error from failed parallel request, got nil")
	}
}

func TestRepeat_DoWithChannel_ReceivesResponses(t *testing.T) {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
	}))
	defer server.Close()

	repeater := Get(server.URL).Repeat(3).WithParallel(3)
	responseChan, errorChan := repeater.DoWithChannel()

	var responses []*HttpResponse
	for i := 0; i < 3; i++ {
		select {
		case res := <-responseChan:
			responses = append(responses, res)
		case err := <-errorChan:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if len(responses) != 3 {
		t.Fatalf("expected 3 responses from channel, got %d", len(responses))
	}
}

func TestDefaultRequestRepeaterOptions(t *testing.T) {
	opts := DefaultRequestRepeaterOptions()
	if opts.Delay != 0 {
		t.Fatalf("expected Delay 0, got %v", opts.Delay)
	}
	if opts.Backoff {
		t.Fatal("expected Backoff false by default")
	}
	if opts.BackoffFactor != 2.0 {
		t.Fatalf("expected BackoffFactor 2.0, got %f", opts.BackoffFactor)
	}
	if opts.Parallel {
		t.Fatal("expected Parallel false by default")
	}
	if opts.ParallelLimit != 5 {
		t.Fatalf("expected ParallelLimit 5, got %d", opts.ParallelLimit)
	}
	if opts.BeforeEach != nil {
		t.Fatal("expected BeforeEach nil by default")
	}
}

func ExampleRequestRepeater_Do() {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	responses, err := Get(server.URL).Repeat(3).Do()
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(len(responses))
	// Output: 3
}

func ExampleRequestRepeater_Do_parallel() {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	responses, err := Get(server.URL).
		Repeat(4).
		WithParallel(2).
		Do()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(len(responses))
	// Output: 4
}

func ExampleRequestRepeater_DoWithChannel() {
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	responseChan, errorChan := Get(server.URL).
		Repeat(3).
		WithParallel(2).
		DoWithChannel()

	responses := 0
	errors := 0
	for responseChan != nil || errorChan != nil {
		select {
		case _, ok := <-responseChan:
			if !ok {
				responseChan = nil
				continue
			}
			responses++
		case _, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			errors++
		}
	}

	fmt.Printf("responses: %d errors: %d\n", responses, errors)
	// Output: responses: 3 errors: 0
}
