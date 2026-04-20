package http

import (
	"sync"
	"time"
)

// Repeat creates a RequestRepeater that can execute the original HTTP request multiple times with configurable options for delay,
// backoff, and parallel execution.
// Example usage:
//
//	responses, err := http.Get("http://httpbin.org/get").
//		Repeat(5).
//		WithDelay(2 * time.Second).
//		WithBackoff(2.0).
//		Do()
func (req *HttpRequest) Repeat(times int) *RequestRepeater {
	return &RequestRepeater{
		times:           times,
		originalRequest: req,
		options:         DefaultRequestRepeaterOptions(),
	}
}

type RequestRepeater struct {
	times           int
	originalRequest *HttpRequest
	options         RequestRepeaterOptions
}

// RequestRepeaterOptions defines the configuration options for repeating HTTP requests.
type RequestRepeaterOptions struct {
	// Delay is the fixed delay between each request execution.
	// If Backoff is enabled, the delay will be multiplied by the BackoffFactor for each subsequent request.
	Delay time.Duration

	// Backoff enables exponential backoff between request executions.
	// The delay will be multiplied by the BackoffFactor for each subsequent request.
	Backoff bool

	// BackoffFactor is the multiplier for the delay when Backoff is enabled.
	// For example, a factor of 2.0 will double the delay for each subsequent request.
	BackoffFactor float64

	// Parallel enables parallel execution of the repeated requests.
	// The ParallelLimit option can be used to limit the number of concurrent requests.
	Parallel bool

	// ParallelLimit is the maximum number of concurrent requests when Parallel is enabled.
	ParallelLimit int

	// BeforeEach is a function that will be called before each request execution,
	// allowing modification of the HttpRequest before it is sent.
	BeforeEach func(*HttpRequest)
}

// DefaultRequestRepeaterOptions returns a RequestRepeaterOptions struct with default values for all options.
func DefaultRequestRepeaterOptions() RequestRepeaterOptions {
	return RequestRepeaterOptions{
		Delay:         0,
		Backoff:       false,
		BackoffFactor: 2.0,
		Parallel:      false,
		ParallelLimit: 5,
		BeforeEach:    nil,
	}
}

// WithDelay sets a fixed delay between each request execution.
// If Backoff is enabled, the delay will be multiplied by the BackoffFactor for each subsequent request.
// If Parallel is enabled, the delay will be disabled since requests will be executed concurrently.
func (m *RequestRepeater) WithDelay(delay time.Duration) *RequestRepeater {
	m.options.Delay = delay
	return m
}

// WithBackoff enables exponential backoff between request executions.
// The delay will be multiplied by the BackoffFactor for each subsequent request.
// If Parallel is enabled, backoff will be disabled since requests will be executed concurrently.
func (m *RequestRepeater) WithBackoff(factor float64) *RequestRepeater {
	m.options.Backoff = true
	m.options.BackoffFactor = factor
	return m
}

// WithParallel enables parallel execution of the repeated requests.
// The ParallelLimit option can be used to limit the number of concurrent requests.
func (m *RequestRepeater) WithParallel(limit int) *RequestRepeater {
	m.options.Parallel = true
	m.options.ParallelLimit = limit
	return m
}

// WithOptions allows setting all RequestRepeaterOptions at once.
func (m *RequestRepeater) WithOptions(options RequestRepeaterOptions) *RequestRepeater {
	m.options = options
	return m
}

// BeforeEach allows setting a function that will be called before each request execution,
// allowing modification of the HttpRequest before it is sent.
func (m *RequestRepeater) BeforeEach(changerFunc func(*HttpRequest)) *RequestRepeater {
	m.options.BeforeEach = changerFunc
	return m
}

// Do executes the original HTTP request the specified number of times, applying the provided options for delay,
// backoff, and parallel execution.
// It returns a slice of HttpResponse pointers and an error if any request fails.
func (m *RequestRepeater) Do() ([]*HttpResponse, error) {
	if m.options.Parallel {
		return m.doParallel()
	}

	return m.do()
}

func (m *RequestRepeater) do() ([]*HttpResponse, error) {
	var responses []*HttpResponse
	for i := 0; i < m.times; i++ {
		request := m.originalRequest.clone()
		if m.options.BeforeEach != nil {
			m.options.BeforeEach(request)
		}

		response, err := request.Do()
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)

		if i < m.times-1 {
			if m.options.Backoff {
				time.Sleep(m.options.Delay * time.Duration(m.options.BackoffFactor*float64(i+1)))
			} else if m.options.Delay > 0 {
				time.Sleep(m.options.Delay)
			}
		}
	}

	return responses, nil
}

func (m *RequestRepeater) doParallel() ([]*HttpResponse, error) {
	responseChan, errorChan := m.doWithChannel()

	var responses []*HttpResponse
	for responseChan != nil || errorChan != nil {
		select {
		case res, ok := <-responseChan:
			if !ok {
				responseChan = nil
				continue
			}
			responses = append(responses, res)
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
				continue
			}
			return nil, err
		}
	}

	return responses, nil
}

func (m *RequestRepeater) doWithChannel() (<-chan *HttpResponse, <-chan error) {
	responseChan := make(chan *HttpResponse, m.times)
	errorChan := make(chan error, 1)
	done := make(chan struct{})

	var once sync.Once
	cancel := func() { once.Do(func() { close(done) }) }

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		sem := make(chan struct{}, m.options.ParallelLimit)
		var wg sync.WaitGroup
		wg.Add(m.times)

		for i := 0; i < m.times; i++ {
			sem <- struct{}{}

			go func() {
				defer wg.Done()
				defer func() { <-sem }()

				select {
				case <-done:
					return
				default:
				}

				request := m.originalRequest.clone()
				if m.options.BeforeEach != nil {
					m.options.BeforeEach(request)
				}

				response, err := request.Do()
				if err != nil {
					select {
					case errorChan <- err:
					default:
					}
					cancel()
					return
				}

				select {
				case responseChan <- response:
				case <-done:
				}
			}()
		}
		wg.Wait()
	}()

	return responseChan, errorChan
}

// DoWithChannel executes the original HTTP request the specified number of times in parallel,
// returning two channels: one for successful responses and one for errors.
// The caller may select on both channels. Both channels are closed when all requests complete.
func (m *RequestRepeater) DoWithChannel() (<-chan *HttpResponse, <-chan error) {
	return m.doWithChannel()
}
