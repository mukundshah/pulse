package checker

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
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
	requestStart time.Time
	dnsStart     time.Time
	dnsDone      time.Time
	connectStart time.Time
	connectDone  time.Time
	tlsStart     time.Time
	tlsDone      time.Time
	gotConn      time.Time
	requestSent  time.Time
	firstByte    time.Time
	responseEnd  time.Time
}

// httpCheckExecutor executes HTTP checks with all necessary configuration.
type httpCheckExecutor struct {
	check            *models.Check
	client           *http.Client
	timings          *timingTracker
	responseSize     int64
	connectionReused bool
	ipVersion        string
	ipAddress        string
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

	// Enforce strict IP version usage
	transport.DialContext = createIPVersionDialer(check.IPVersion)

	// Configure TLS based on SkipSSLVerification setting
	if check.SkipSSLVerification {
		// Skip SSL certificate verification
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
	} else {
		// Use secure TLS defaults with certificate verification
		transport.TLSClientConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
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
		check:            check,
		client:           client,
		timings:          &timingTracker{},
		responseSize:     0,
		connectionReused: false,
	}
}

// createIPVersionDialer creates a dialer that enforces strict IP version usage.
func createIPVersionDialer(ipVersion models.IPVersionType) func(context.Context, string, string) (net.Conn, error) {
	baseDialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return func(ctx context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, fmt.Errorf("invalid address format: %w", err)
		}

		// Determine the network type based on IP version requirement
		// Use tcp4 or tcp6 for strict enforcement
		var targetNetwork string
		if ipVersion == models.IPVersionTypeIPv4 {
			targetNetwork = "tcp4"
		} else {
			targetNetwork = "tcp6"
		}

		// Check if host is already an IP address
		ip := net.ParseIP(host)
		if ip != nil {
			// Direct IP address - verify it matches the required version
			isIPv4 := ip.To4() != nil
			requiresIPv4 := ipVersion == models.IPVersionTypeIPv4

			if isIPv4 != requiresIPv4 {
				if requiresIPv4 {
					return nil, fmt.Errorf("IP version mismatch: required IPv4 but got IPv6 address %s", host)
				}
				return nil, fmt.Errorf("IP version mismatch: required IPv6 but got IPv4 address %s", host)
			}

			// IP version matches, proceed with connection using strict network type
			return baseDialer.DialContext(ctx, targetNetwork, address)
		}

		// Hostname - resolve and filter by IP version
		var addrs []net.IP
		if ipVersion == models.IPVersionTypeIPv4 {
			// Resolve only IPv4 addresses
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, fmt.Errorf("DNS lookup failed: %w", err)
			}

			for _, ipAddr := range ips {
				if ipAddr.IP.To4() != nil {
					addrs = append(addrs, ipAddr.IP)
				}
			}

			if len(addrs) == 0 {
				return nil, fmt.Errorf("IP version mismatch: no IPv4 addresses found for host %s", host)
			}
		} else {
			// Resolve only IPv6 addresses
			ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
			if err != nil {
				return nil, fmt.Errorf("DNS lookup failed: %w", err)
			}

			for _, ipAddr := range ips {
				if ipAddr.IP.To4() == nil {
					addrs = append(addrs, ipAddr.IP)
				}
			}

			if len(addrs) == 0 {
				return nil, fmt.Errorf("IP version mismatch: no IPv6 addresses found for host %s", host)
			}
		}

		// Try connecting to each resolved address using strict network type
		var lastErr error
		for _, addr := range addrs {
			addrStr := net.JoinHostPort(addr.String(), port)
			conn, err := baseDialer.DialContext(ctx, targetNetwork, addrStr)
			if err == nil {
				return conn, nil
			}
			lastErr = err
		}

		if lastErr != nil {
			return nil, fmt.Errorf("failed to connect to any %s address for %s: %w", ipVersion, host, lastErr)
		}

		return nil, fmt.Errorf("no %s addresses available for %s", ipVersion, host)
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

	// CRITICAL: Start timer before request
	e.timings.requestStart = time.Now().UTC()

	// Execute request (this returns at TTFB, not after body download)
	resp, httpErr := e.client.Do(req)

	// Handle response error BEFORE reading body
	if httpErr != nil {
		return e.createErrorResult(httpErr)
	}
	defer resp.Body.Close()

	// CRITICAL: Read body fully to measure download time
	// Read body into memory so we can use it for assertions too
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return e.createErrorResult(fmt.Errorf("failed to read response body: %w", err))
	}

	// NOW stop the timer (after body is fully read)
	e.timings.responseEnd = time.Now().UTC()
	e.responseSize = int64(len(bodyBytes))

	// Create new reader for assertions (they may need the body)
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// Process assertions
	responseTime := e.responseTime()
	assertionResults, err := e.processAssertions(resp, responseTime)
	if err != nil {
		return e.createErrorResult(fmt.Errorf("%w: %v", ErrAssertionProcessing, err))
	}

	// Build result with timestamps
	return e.buildResult(resp, assertionResults)
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

	if e.check.QueryParams != nil && len(e.check.QueryParams) > 0 {
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
			e.timings.dnsStart = time.Now().UTC()
		},
		DNSDone: func(httptrace.DNSDoneInfo) {
			e.timings.dnsDone = time.Now().UTC()
		},
		ConnectStart: func(_, _ string) {
			e.timings.connectStart = time.Now().UTC()
		},
		ConnectDone: func(_, _ string, _ error) {
			e.timings.connectDone = time.Now().UTC()
		},
		TLSHandshakeStart: func() {
			e.timings.tlsStart = time.Now().UTC()
		},
		TLSHandshakeDone: func(tls.ConnectionState, error) {
			e.timings.tlsDone = time.Now().UTC()
		},
		GotConn: func(info httptrace.GotConnInfo) {
			e.timings.gotConn = time.Now().UTC()
			e.connectionReused = info.Reused

			// Extract IP address and version from connection
			if info.Conn != nil {
				host, _, err := net.SplitHostPort(info.Conn.RemoteAddr().String())
				if err == nil {
					e.ipAddress = host
					ip := net.ParseIP(host)
					if ip != nil {
						if ip.To4() != nil {
							e.ipVersion = "IPv4"
						} else {
							e.ipVersion = "IPv6"
						}
					}
				}
			}
		},
		WroteRequest: func(httptrace.WroteRequestInfo) {
			e.timings.requestSent = time.Now().UTC()
		},
		GotFirstResponseByte: func() {
			e.timings.firstByte = time.Now().UTC()
		},
	}

	return req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
}

// processAssertions evaluates all assertions against the response.
func (e *httpCheckExecutor) processAssertions(resp *http.Response, responseTime time.Duration) ([]AssertionResult, error) {
	if len(e.check.Assertions) == 0 {
		return []AssertionResult{}, nil
	}

	return ProcessAssertions(e.check.Assertions, resp, responseTime)
}

// responseTime calculates the total response time from timestamps.
func (e *httpCheckExecutor) responseTime() time.Duration {
	if e.timings.requestStart.IsZero() || e.timings.responseEnd.IsZero() {
		return 0
	}
	return e.timings.responseEnd.Sub(e.timings.requestStart)
}

// buildNetworkTimings creates a map of detailed network timing information.
// Stores both raw timestamps (for timeline reconstruction) and durations (for convenience).
func (e *httpCheckExecutor) buildNetworkTimings() map[string]interface{} {
	timings := make(map[string]interface{})

	// Store raw timestamps (for timeline reconstruction and debugging)
	if !e.timings.requestStart.IsZero() {
		timings["request_start"] = e.timings.requestStart.Format(time.RFC3339Nano)
	}
	if !e.timings.dnsStart.IsZero() {
		timings["dns_start"] = e.timings.dnsStart.Format(time.RFC3339Nano)
	}
	if !e.timings.dnsDone.IsZero() {
		timings["dns_done"] = e.timings.dnsDone.Format(time.RFC3339Nano)
	}
	if !e.timings.connectStart.IsZero() {
		timings["tcp_start"] = e.timings.connectStart.Format(time.RFC3339Nano)
	}
	if !e.timings.connectDone.IsZero() {
		timings["tcp_done"] = e.timings.connectDone.Format(time.RFC3339Nano)
	}
	if !e.timings.tlsStart.IsZero() {
		timings["tls_start"] = e.timings.tlsStart.Format(time.RFC3339Nano)
	}
	if !e.timings.tlsDone.IsZero() {
		timings["tls_done"] = e.timings.tlsDone.Format(time.RFC3339Nano)
	}
	if !e.timings.requestSent.IsZero() {
		timings["request_sent"] = e.timings.requestSent.Format(time.RFC3339Nano)
	}
	if !e.timings.firstByte.IsZero() {
		timings["first_byte"] = e.timings.firstByte.Format(time.RFC3339Nano)
	}
	if !e.timings.responseEnd.IsZero() {
		timings["response_end"] = e.timings.responseEnd.Format(time.RFC3339Nano)
	}

	// Compute durations from timestamps (in microseconds)
	// Only compute if both timestamps are valid and end is after start
	// DNS duration: dns_done - dns_start
	if !e.timings.dnsStart.IsZero() && !e.timings.dnsDone.IsZero() && e.timings.dnsDone.After(e.timings.dnsStart) {
		if us := durationUs(e.timings.dnsStart, e.timings.dnsDone); us > 0 {
			timings["dns_duration_us"] = us
		}
	}
	// TCP duration: tcp_done - tcp_start
	if !e.timings.connectStart.IsZero() && !e.timings.connectDone.IsZero() && e.timings.connectDone.After(e.timings.connectStart) {
		if us := durationUs(e.timings.connectStart, e.timings.connectDone); us > 0 {
			timings["tcp_duration_us"] = us
		}
	}
	// TLS duration: tls_done - tls_start
	if !e.timings.tlsStart.IsZero() && !e.timings.tlsDone.IsZero() && e.timings.tlsDone.After(e.timings.tlsStart) {
		if us := durationUs(e.timings.tlsStart, e.timings.tlsDone); us > 0 {
			timings["tls_duration_us"] = us
		}
	}
	// Request send duration: request_sent - tls_done (or request_sent - request_start if no TLS)
	// Only compute if request_sent is after tls_done (or request_start if no TLS)
	if !e.timings.requestSent.IsZero() {
		var requestDurationUs int
		if !e.timings.tlsDone.IsZero() && e.timings.requestSent.After(e.timings.tlsDone) {
			requestDurationUs = durationUs(e.timings.tlsDone, e.timings.requestSent)
		} else if !e.timings.requestStart.IsZero() && e.timings.requestSent.After(e.timings.requestStart) {
			requestDurationUs = durationUs(e.timings.requestStart, e.timings.requestSent)
		}
		if requestDurationUs > 0 {
			timings["request_duration_us"] = requestDurationUs
		}
	}
	// TTFB: first_byte - request_sent (only if both are available and first_byte is after request_sent)
	if !e.timings.requestSent.IsZero() && !e.timings.firstByte.IsZero() && e.timings.firstByte.After(e.timings.requestSent) {
		if us := durationUs(e.timings.requestSent, e.timings.firstByte); us > 0 {
			timings["ttfb_us"] = us
		}
	}
	// Download duration: response_end - first_byte
	if !e.timings.firstByte.IsZero() && !e.timings.responseEnd.IsZero() && e.timings.responseEnd.After(e.timings.firstByte) {
		if us := durationUs(e.timings.firstByte, e.timings.responseEnd); us > 0 {
			timings["download_us"] = us
		}
	}
	// Total response time: response_end - request_start
	if !e.timings.requestStart.IsZero() && !e.timings.responseEnd.IsZero() && e.timings.responseEnd.After(e.timings.requestStart) {
		responseTime := e.responseTime()
		if responseTime > 0 {
			timings["response_time_us"] = int(responseTime / time.Microsecond)
		}
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
func (e *httpCheckExecutor) buildResult(resp *http.Response, assertionResults []AssertionResult) Result {
	// Compute response time from timestamps
	responseTime := e.responseTime()

	// Build network timings (raw timestamps for JSON)
	networkTimings := e.buildNetworkTimings()

	// Determine status
	status := e.determineStatus(responseTime, assertionResults)

	// Determine failure reason if failed
	var failureReason *models.FailureReason
	if status == models.CheckRunStatusFailing {
		failureReason = e.determineFailureReason(resp, responseTime, assertionResults)
	}

	// Validate timeline invariants
	if !e.timings.requestStart.IsZero() && !e.timings.responseEnd.IsZero() {
		if e.timings.responseEnd.Before(e.timings.requestStart) {
			// Invalid timeline - mark as agent error
			status = models.CheckRunStatusFailing
			failureReason = failureReasonPtr(models.FailureAgent)
		}
	}
	if !e.timings.firstByte.IsZero() {
		if e.timings.firstByte.Before(e.timings.requestStart) {
			status = models.CheckRunStatusFailing
			failureReason = failureReasonPtr(models.FailureAgent)
		}
		if !e.timings.responseEnd.IsZero() && e.timings.responseEnd.Before(e.timings.firstByte) {
			status = models.CheckRunStatusFailing
			failureReason = failureReasonPtr(models.FailureAgent)
		}
	}

	// Response status (nullable)
	var responseStatus *int32
	if resp != nil {
		code := int32(resp.StatusCode)
		responseStatus = &code
	}

	return Result{
		Status:            status,
		FailureReason:     failureReason,
		ResponseStatus:    responseStatus,
		RequestStartedAt:  e.timings.requestStart,
		FirstByteAt:       e.timings.firstByte,
		ResponseEndedAt:   e.timings.responseEnd,
		ConnectionReused:  e.connectionReused,
		IPVersion:         e.ipVersion,
		IPAddress:         e.ipAddress,
		ResponseSizeBytes: e.responseSize,
		AssertionResults:  mustMarshalJSON(assertionResults),
		PlaywrightReport:  emptyJSONObject(),
		NetworkTimings:    mustMarshalJSON(networkTimings),
		Error:             nil,
	}
}

// createErrorResult creates a result for a failed check.
func (e *httpCheckExecutor) createErrorResult(err error) Result {
	// Determine failure reason from error type
	failureReason := e.classifyError(err)

	// Response status is nil on error
	var responseStatus *int32

	// Timestamps may be partial
	requestStart := e.timings.requestStart
	if requestStart.IsZero() {
		requestStart = time.Now().UTC()
	}

	return Result{
		Status:            models.CheckRunStatusFailing,
		FailureReason:     failureReason,
		ResponseStatus:    responseStatus,
		RequestStartedAt:  requestStart,
		FirstByteAt:       e.timings.firstByte,   // May be zero
		ResponseEndedAt:   e.timings.responseEnd, // May be zero
		ConnectionReused:  e.connectionReused,
		IPVersion:         e.ipVersion,
		IPAddress:         e.ipAddress,
		ResponseSizeBytes: 0,
		AssertionResults:  emptyJSONObject(),
		PlaywrightReport:  emptyJSONObject(),
		NetworkTimings:    emptyJSONObject(),
		Error:             err,
	}
}

// failureReasonPtr returns a pointer to a FailureReason constant.
func failureReasonPtr(r models.FailureReason) *models.FailureReason {
	return &r
}

// classifyError determines the failure reason from an error.
func (e *httpCheckExecutor) classifyError(err error) *models.FailureReason {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Network errors
	if contains(errStr, "ip version mismatch") {
		return failureReasonPtr(models.FailureIPVersionMismatch)
	}
	if contains(errStr, "no such host") || contains(errStr, "dns") {
		return failureReasonPtr(models.FailureDNS)
	}
	if contains(errStr, "connection refused") {
		return failureReasonPtr(models.FailureConnectionRefused)
	}
	if contains(errStr, "timeout") || contains(errStr, "deadline exceeded") {
		return failureReasonPtr(models.FailureRequestTimeout)
	}
	// TLS/SSL errors - check before other connection errors
	if contains(errStr, "tls") || contains(errStr, "ssl") || contains(errStr, "certificate") ||
		contains(errStr, "handshake") || contains(errStr, "x509") || contains(errStr, "certificate verify failed") {
		return failureReasonPtr(models.FailureTLS)
	}
	if contains(errStr, "connection") && contains(errStr, "reset") {
		return failureReasonPtr(models.FailureTCP)
	}
	if contains(errStr, "network is unreachable") {
		return failureReasonPtr(models.FailureNetworkUnreachable)
	}

	return failureReasonPtr(models.FailureUnknown)
}

// contains checks if a string contains a substring (case-insensitive).
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// determineFailureReason determines the failure reason based on response and assertions.
func (e *httpCheckExecutor) determineFailureReason(resp *http.Response, responseTime time.Duration, assertionResults []AssertionResult) *models.FailureReason {
	// Check assertions first
	for _, result := range assertionResults {
		if !result.Passed {
			return failureReasonPtr(models.FailureAssertionFailed)
		}
	}

	// Check HTTP status
	if resp != nil {
		if resp.StatusCode >= 500 {
			return failureReasonPtr(models.FailureHTTP5xx)
		}
		if resp.StatusCode >= 400 {
			return failureReasonPtr(models.FailureHTTP4xx)
		}
	}

	// Check timeouts (if applicable)
	// Determine specific timeout type based on responseTime and thresholds
	failedThreshold := e.check.FailedThresholdDuration()
	if failedThreshold > 0 && responseTime > failedThreshold {
		// Total request exceeded threshold - determine if it's TTFB or download timeout
		if !e.timings.firstByte.IsZero() && !e.timings.requestStart.IsZero() {
			ttfb := e.timings.firstByte.Sub(e.timings.requestStart)
			downloadTime := responseTime - ttfb

			// TTFB timeout: TTFB exceeds the failed threshold OR TTFB is >80% of total time
			// This indicates the server took too long to start responding
			if ttfb > 0 && (ttfb > failedThreshold || ttfb > responseTime*4/5) {
				return failureReasonPtr(models.FailureTTFBTimeout)
			}

			// Download timeout: download time is significant (>80% of total)
			// AND exceeds a reasonable threshold (50% of failed threshold)
			// This indicates the response body download was too slow
			downloadThreshold := failedThreshold / 2
			if downloadTime > 0 && downloadTime > downloadThreshold && downloadTime > responseTime*4/5 {
				return failureReasonPtr(models.FailureDownloadTimeout)
			}
		}

		// Default to general request timeout when we can't determine specific timeout type
		// or when total time exceeded threshold but timestamps aren't available
		return failureReasonPtr(models.FailureRequestTimeout)
	}

	// Default to unknown
	return failureReasonPtr(models.FailureUnknown)
}

// durationUs calculates duration in microseconds between two times.
// Returns 0 if either time is zero or end is before start.
func durationUs(start, end time.Time) int {
	if start.IsZero() || end.IsZero() || end.Before(start) {
		return 0
	}
	return int(end.Sub(start) / time.Microsecond)
}
