package checker

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"

	"pulse/internal/models"
)

const (
	defaultTimeout      = 30 * time.Second
	defaultMaxRedirects = 10
)

var (
	// ErrRequestCreation is returned when HTTP request creation fails.
	ErrRequestCreation = fmt.Errorf("failed to create HTTP request")
	// ErrAssertionProcessing is returned when assertion processing fails.
	ErrAssertionProcessing = fmt.Errorf("failed to process assertions")
)

// timingTracker tracks HTTP request timing events.
type timingTracker struct {
	dnsStart     time.Time
	dnsDone      time.Time
	connectStart time.Time
	connectDone  time.Time
	tlsStart     time.Time
	tlsDone      time.Time
	gotConn      time.Time
	firstByte    time.Time
}

// httpCheckExecutor executes HTTP checks with all necessary configuration.
type httpCheckExecutor struct {
	check     *models.Check
	client    *http.Client
	timings   *timingTracker
	startTime time.Time
	endTime   time.Time
}

// ExecuteHTTPCheck performs an HTTP check and returns the result.
// It handles request creation, execution, timing tracking, and assertion evaluation.
func ExecuteHTTPCheck(ctx context.Context, check *models.Check) Result {
	executor := newHTTPCheckExecutor(check)
	return executor.execute(ctx)
}

// newHTTPCheckExecutor creates a new HTTP check executor with configured client.
func newHTTPCheckExecutor(check *models.Check) *httpCheckExecutor {
	transport := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: false,
		DisableKeepAlives:  false,
	}

	if check.SkipSSLVerification {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
	}

	client := &http.Client{
		Timeout:   defaultTimeout,
		Transport: transport,
	}

	if !check.FollowRedirects {
		client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &httpCheckExecutor{
		check:   check,
		client:  client,
		timings: &timingTracker{},
	}
}

// execute runs the HTTP check and returns the result.
func (e *httpCheckExecutor) execute(ctx context.Context) Result {
	// Build request
	req, err := e.buildRequest(ctx)
	if err != nil {
		return e.createErrorResult(err)
	}

	// Add tracing
	req = e.addTracing(req)

	// Execute request
	e.startTime = time.Now()
	resp, httpErr := e.client.Do(req)
	e.endTime = time.Now()

	// Handle response error
	if httpErr != nil {
		return e.createErrorResult(httpErr)
	}
	defer resp.Body.Close()

	// Calculate response time
	responseTime := e.responseTime()

	// Process assertions
	assertionResults, err := e.processAssertions(resp, responseTime)
	if err != nil {
		return e.createErrorResult(fmt.Errorf("%w: %v", ErrAssertionProcessing, err))
	}

	// Build result
	return e.buildResult(resp, responseTime, assertionResults)
}

// buildRequest creates an HTTP request from the check configuration.
func (e *httpCheckExecutor) buildRequest(ctx context.Context) (*http.Request, error) {
	targetURL, err := e.buildURL()
	if err != nil {
		return nil, fmt.Errorf("%w: invalid URL: %v", ErrRequestCreation, err)
	}

	bodyReader, err := e.buildBody()
	if err != nil {
		return nil, fmt.Errorf("%w: invalid body: %v", ErrRequestCreation, err)
	}

	req, err := http.NewRequestWithContext(ctx, e.check.Method, targetURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequestCreation, err)
	}

	if err := e.setHeaders(req); err != nil {
		return nil, fmt.Errorf("%w: invalid headers: %v", ErrRequestCreation, err)
	}

	return req, nil
}

// buildURL constructs the target URL from check configuration.
func (e *httpCheckExecutor) buildURL() (string, error) {
	u := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", e.check.Host, e.check.Port),
		Path:   e.check.Path,
	}

	if e.check.Secure {
		u.Scheme = "https"
	}

	if len(e.check.QueryParams) > 0 {
		if err := e.setQueryParams(&u); err != nil {
			return "", err
		}
	}

	return u.String(), nil
}

// setQueryParams adds query parameters to the URL.
func (e *httpCheckExecutor) setQueryParams(u *url.URL) error {
	var queryParams map[string]interface{}
	if err := json.Unmarshal(e.check.QueryParams, &queryParams); err != nil {
		return fmt.Errorf("failed to unmarshal query params: %w", err)
	}

	q := u.Query()
	for k, v := range queryParams {
		q.Set(k, fmt.Sprint(v))
	}
	u.RawQuery = q.Encode()

	return nil
}

// buildBody creates a reader for the request body.
func (e *httpCheckExecutor) buildBody() (io.Reader, error) {
	if len(e.check.Body) == 0 {
		return nil, nil
	}
	return bytes.NewReader(e.check.Body), nil
}

// setHeaders sets HTTP headers from check configuration.
func (e *httpCheckExecutor) setHeaders(req *http.Request) error {
	if e.check.Headers == nil || len(e.check.Headers) == 0 {
		return nil
	}

	var headers map[string]interface{}
	if err := json.Unmarshal(e.check.Headers, &headers); err != nil {
		return fmt.Errorf("failed to unmarshal headers: %w", err)
	}

	for k, v := range headers {
		// Support both string and array of strings
		switch val := v.(type) {
		case string:
			req.Header.Set(k, val)
		case []interface{}:
			for _, item := range val {
				if str, ok := item.(string); ok {
					req.Header.Add(k, str)
				}
			}
		}
	}

	return nil
}

// addTracing adds HTTP trace callbacks to track timing information.
func (e *httpCheckExecutor) addTracing(req *http.Request) *http.Request {
	trace := &httptrace.ClientTrace{
		DNSStart: func(httptrace.DNSStartInfo) {
			e.timings.dnsStart = time.Now()
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			e.timings.dnsDone = time.Now()
		},
		ConnectStart: func(_, _ string) {
			e.timings.connectStart = time.Now()
		},
		ConnectDone: func(_, _ string, _ error) {
			e.timings.connectDone = time.Now()
		},
		TLSHandshakeStart: func() {
			e.timings.tlsStart = time.Now()
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			e.timings.tlsDone = time.Now()
		},
		GotConn: func(httptrace.GotConnInfo) {
			e.timings.gotConn = time.Now()
		},
		GotFirstResponseByte: func() {
			e.timings.firstByte = time.Now()
		},
	}

	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}

// processAssertions evaluates all assertions against the response.
func (e *httpCheckExecutor) processAssertions(resp *http.Response, responseTime time.Duration) ([]AssertionResult, error) {
	if e.check.Assertions == nil || len(e.check.Assertions) == 0 {
		return []AssertionResult{}, nil
	}

	return ProcessAssertions(e.check.Assertions, resp, responseTime)
}

// responseTime calculates the total response time.
func (e *httpCheckExecutor) responseTime() time.Duration {
	return e.endTime.Sub(e.startTime)
}

// buildNetworkTimings creates a map of detailed network timing information.
func (e *httpCheckExecutor) buildNetworkTimings(responseTime time.Duration) map[string]int {
	timings := make(map[string]int)

	if ms := durationMs(e.timings.dnsStart, e.timings.dnsDone); ms > 0 {
		timings["dns_lookup_ms"] = ms
	}
	if ms := durationMs(e.timings.connectStart, e.timings.connectDone); ms > 0 {
		timings["tcp_connection_ms"] = ms
	}
	if ms := durationMs(e.timings.tlsStart, e.timings.tlsDone); ms > 0 {
		timings["tls_handshake_ms"] = ms
	}
	if ms := durationMs(e.startTime, e.timings.firstByte); ms > 0 {
		timings["time_to_first_byte_ms"] = ms
	}
	if ms := durationMs(e.timings.gotConn, e.timings.firstByte); ms > 0 {
		timings["server_processing_ms"] = ms
	}
	if ms := durationMs(e.timings.firstByte, e.endTime); ms > 0 {
		timings["content_transfer_ms"] = ms
	}
	if responseTime > 0 {
		timings["response_time_ms"] = int(responseTime / time.Millisecond)
	}

	return timings
}

// determineStatus calculates the check status based on response and thresholds.
func (e *httpCheckExecutor) determineStatus(responseTime time.Duration, assertionResults []AssertionResult) models.CheckRunStatus {
	// Check if any assertions failed
	for _, result := range assertionResults {
		if !result.Passed {
			return models.CheckRunStatusFailing
		}
	}

	// Check response time thresholds
	if responseTime > e.check.FailedThresholdDuration() {
		return models.CheckRunStatusFailing
	}
	if responseTime > e.check.DegradedThresholdDuration() {
		return models.CheckRunStatusDegraded
	}

	return models.CheckRunStatusPassing
}

// buildResult creates the final result object.
func (e *httpCheckExecutor) buildResult(resp *http.Response, responseTime time.Duration, assertionResults []AssertionResult) Result {
	networkTimings := e.buildNetworkTimings(responseTime)
	status := e.determineStatus(responseTime, assertionResults)

	return Result{
		Status:           status,
		ResponseStatus:   int32(resp.StatusCode),
		TotalTimeMs:      int(responseTime / time.Millisecond),
		AssertionResults: mustMarshalJSON(assertionResults),
		PlaywrightReport: emptyJSONObject(),
		NetworkTimings:   mustMarshalJSON(networkTimings),
		Error:            nil,
	}
}

// createErrorResult creates a result for a failed check.
func (e *httpCheckExecutor) createErrorResult(err error) Result {
	return Result{
		Status:           models.CheckRunStatusFailing,
		ResponseStatus:   0,
		TotalTimeMs:      0,
		AssertionResults: emptyJSONObject(),
		PlaywrightReport: emptyJSONObject(),
		NetworkTimings:   emptyJSONObject(),
		Error:            err,
	}
}

// durationMs calculates duration in milliseconds between two times.
// Returns 0 if either time is zero or end is before start.
func durationMs(start, end time.Time) int {
	if start.IsZero() || end.IsZero() || end.Before(start) {
		return 0
	}
	return int(end.Sub(start) / time.Millisecond)
}
