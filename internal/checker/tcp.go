package checker

import (
	"context"
	"fmt"
	"net"
	"time"

	"pulse/internal/models"
)

// tcpCheckExecutor executes TCP checks with all necessary configuration.
type tcpCheckExecutor struct {
	check            *models.Check
	timings          *tcpTimingTracker
	ipVersion        string
	ipAddress        string
	connectionReused bool
}

// tcpTimingTracker tracks TCP connection timing events.
type tcpTimingTracker struct {
	requestStart time.Time
	dnsStart     time.Time
	dnsDone      time.Time
	connectStart time.Time
	connectDone  time.Time
	responseEnd  time.Time
}

// ExecuteTCPCheck performs a TCP check and returns the result.
// It handles DNS resolution, connection establishment, and timing tracking.
func ExecuteTCPCheck(ctx context.Context, check *models.Check) Result {
	executor := newTCPCheckExecutor(check)
	return executor.execute(ctx)
}

// newTCPCheckExecutor creates a new TCP check executor.
func newTCPCheckExecutor(check *models.Check) *tcpCheckExecutor {
	return &tcpCheckExecutor{
		check:            check,
		timings:          &tcpTimingTracker{},
		connectionReused: false,
	}
}

// execute runs the TCP check and returns the result.
func (e *tcpCheckExecutor) execute(ctx context.Context) Result {
	// Start timer before connection attempt
	e.timings.requestStart = time.Now().UTC()

	// Resolve address
	address, err := e.resolveAddress(ctx)
	if err != nil {
		return e.createErrorResult(err)
	}

	// Establish TCP connection
	conn, err := e.connect(ctx, address)
	if err != nil {
		return e.createErrorResult(err)
	}
	defer conn.Close()

	// Connection successful - record end time
	e.timings.responseEnd = time.Now().UTC()

	// Extract IP information from connection
	e.extractIPInfo(conn)

	// Build result
	return e.buildResult()
}

// resolveAddress resolves the hostname to an IP address with strict IP version enforcement.
func (e *tcpCheckExecutor) resolveAddress(ctx context.Context) (string, error) {
	host := e.check.Host
	port := e.check.Port

	// Check if host is already an IP address
	ip := net.ParseIP(host)
	if ip != nil {
		// Direct IP address - verify it matches the required version
		isIPv4 := ip.To4() != nil
		requiresIPv4 := e.check.IPVersion == models.IPVersionTypeIPv4

		if isIPv4 != requiresIPv4 {
			if requiresIPv4 {
				return "", fmt.Errorf("IP version mismatch: required IPv4 but got IPv6 address %s", host)
			}
			return "", fmt.Errorf("IP version mismatch: required IPv6 but got IPv4 address %s", host)
		}

		// IP version matches, format address for connection
		address := net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port))
		return address, nil
	}

	// Track DNS resolution timing
	e.timings.dnsStart = time.Now().UTC()

	// Resolve address
	resolver := net.DefaultResolver
	addresses, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		e.timings.dnsDone = time.Now().UTC()
		return "", fmt.Errorf("dns resolution failed: %w", err)
	}

	e.timings.dnsDone = time.Now().UTC()

	// Select appropriate IP address based on strict version requirement
	var selectedIP net.IP
	if e.check.IPVersion == models.IPVersionTypeIPv6 {
		// Filter for IPv6 addresses only
		for _, addr := range addresses {
			if addr.IP.To4() == nil {
				selectedIP = addr.IP
				break
			}
		}

		if selectedIP == nil {
			return "", fmt.Errorf("IP version mismatch: no IPv6 addresses found for host %s", host)
		}
	} else {
		// Filter for IPv4 addresses only
		for _, addr := range addresses {
			if addr.IP.To4() != nil {
				selectedIP = addr.IP
				break
			}
		}

		if selectedIP == nil {
			return "", fmt.Errorf("IP version mismatch: no IPv4 addresses found for host %s", host)
		}
	}

	// Format address for connection
	address := net.JoinHostPort(selectedIP.String(), fmt.Sprintf("%d", port))
	return address, nil
}

// connect establishes a TCP connection to the address with strict IP version enforcement.
func (e *tcpCheckExecutor) connect(ctx context.Context, address string) (net.Conn, error) {
	// Track connection timing
	e.timings.connectStart = time.Now().UTC()

	// Validate IP version from address before connecting
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		e.timings.connectDone = time.Now().UTC()
		return nil, fmt.Errorf("invalid address format: %w", err)
	}

	ip := net.ParseIP(host)
	if ip != nil {
		// Verify the resolved IP matches the required version
		isIPv4 := ip.To4() != nil
		requiresIPv4 := e.check.IPVersion == models.IPVersionTypeIPv4

		if isIPv4 != requiresIPv4 {
			e.timings.connectDone = time.Now().UTC()
			if requiresIPv4 {
				return nil, fmt.Errorf("IP version mismatch: required IPv4 but got IPv6 address %s", host)
			}
			return nil, fmt.Errorf("IP version mismatch: required IPv6 but got IPv4 address %s", host)
		}
	}

	// Create dialer with timeout
	timeout := defaultTimeout
	if e.check.FailedThresholdDuration() > 0 {
		timeout = e.check.FailedThresholdDuration()
	}

	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// Determine network type based on IP version (strict enforcement)
	var network string
	if e.check.IPVersion == models.IPVersionTypeIPv6 {
		network = "tcp6"
	} else {
		network = "tcp4"
	}

	// Establish connection
	conn, err := dialer.DialContext(ctx, network, address)
	if err != nil {
		e.timings.connectDone = time.Now().UTC()
		return nil, err
	}

	e.timings.connectDone = time.Now().UTC()
	return conn, nil
}

// extractIPInfo extracts IP version and address from the connection.
func (e *tcpCheckExecutor) extractIPInfo(conn net.Conn) {
	if conn == nil {
		return
	}

	remoteAddr := conn.RemoteAddr()
	if remoteAddr == nil {
		return
	}

	host, _, err := net.SplitHostPort(remoteAddr.String())
	if err != nil {
		return
	}

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

// buildNetworkTimings creates a map of detailed network timing information.
func (e *tcpCheckExecutor) buildNetworkTimings() map[string]interface{} {
	timings := make(map[string]interface{})

	// Store raw timestamps
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
	if !e.timings.responseEnd.IsZero() {
		timings["response_end"] = e.timings.responseEnd.Format(time.RFC3339Nano)
	}

	// Compute durations (in microseconds)
	// DNS duration: dns_done - dns_start
	if !e.timings.dnsStart.IsZero() && !e.timings.dnsDone.IsZero() && e.timings.dnsDone.After(e.timings.dnsStart) {
		if us := durationUs(e.timings.dnsStart, e.timings.dnsDone); us > 0 {
			timings["dns_duration_us"] = us
		}
	}
	// TCP connection duration: tcp_done - tcp_start
	if !e.timings.connectStart.IsZero() && !e.timings.connectDone.IsZero() && e.timings.connectDone.After(e.timings.connectStart) {
		if us := durationUs(e.timings.connectStart, e.timings.connectDone); us > 0 {
			timings["tcp_duration_us"] = us
		}
	}
	// Total connection time: response_end - request_start
	if !e.timings.requestStart.IsZero() && !e.timings.responseEnd.IsZero() && e.timings.responseEnd.After(e.timings.requestStart) {
		responseTime := e.responseTime()
		if responseTime > 0 {
			timings["response_time_us"] = int(responseTime / time.Microsecond)
		}
	}

	return timings
}

// responseTime calculates the total response time from timestamps.
func (e *tcpCheckExecutor) responseTime() time.Duration {
	if e.timings.requestStart.IsZero() || e.timings.responseEnd.IsZero() {
		return 0
	}
	return e.timings.responseEnd.Sub(e.timings.requestStart)
}

// determineStatus calculates the check status based on response time and thresholds.
func (e *tcpCheckExecutor) determineStatus(responseTime time.Duration) models.CheckRunStatus {
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
func (e *tcpCheckExecutor) buildResult() Result {
	// Compute response time from timestamps
	responseTime := e.responseTime()

	// Build network timings
	networkTimings := e.buildNetworkTimings()

	// Determine status
	status := e.determineStatus(responseTime)

	// Determine failure reason if failed
	var failureReason *models.FailureReason
	if status == models.CheckRunStatusFailing {
		failureReason = e.determineFailureReason(responseTime)
	}

	// Validate timeline invariants
	if !e.timings.requestStart.IsZero() && !e.timings.responseEnd.IsZero() {
		if e.timings.responseEnd.Before(e.timings.requestStart) {
			// Invalid timeline - mark as agent error
			status = models.CheckRunStatusFailing
			failureReason = failureReasonPtr(models.FailureAgent)
		}
	}

	return Result{
		Status:            status,
		FailureReason:     failureReason,
		ResponseStatus:    nil, // TCP doesn't have HTTP status codes
		RequestStartedAt:  e.timings.requestStart,
		FirstByteAt:       e.timings.connectDone, // For TCP, connection established is like "first byte"
		ResponseEndedAt:   e.timings.responseEnd,
		ConnectionReused:  e.connectionReused,
		IPVersion:         e.ipVersion,
		IPAddress:         e.ipAddress,
		ResponseSizeBytes: 0,                 // TCP connection check doesn't transfer data
		AssertionResults:  emptyJSONObject(), // TCP checks don't have assertions yet
		PlaywrightReport:  emptyJSONObject(),
		NetworkTimings:    mustMarshalJSON(networkTimings),
		Error:             nil,
	}
}

// createErrorResult creates a result for a failed check.
func (e *tcpCheckExecutor) createErrorResult(err error) Result {
	// Determine failure reason from error type
	failureReason := e.classifyError(err)

	// Timestamps may be partial
	requestStart := e.timings.requestStart
	if requestStart.IsZero() {
		requestStart = time.Now().UTC()
	}

	return Result{
		Status:            models.CheckRunStatusFailing,
		FailureReason:     failureReason,
		ResponseStatus:    nil,
		RequestStartedAt:  requestStart,
		FirstByteAt:       e.timings.connectDone, // May be zero
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

// classifyError determines the failure reason from an error.
func (e *tcpCheckExecutor) classifyError(err error) *models.FailureReason {
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
		return failureReasonPtr(models.FailureConnectionTimeout)
	}
	if contains(errStr, "connection") && contains(errStr, "reset") {
		return failureReasonPtr(models.FailureTCP)
	}
	if contains(errStr, "network is unreachable") {
		return failureReasonPtr(models.FailureNetworkUnreachable)
	}
	if contains(errStr, "no route to host") {
		return failureReasonPtr(models.FailureNetworkUnreachable)
	}

	return failureReasonPtr(models.FailureUnknown)
}

// determineFailureReason determines the failure reason based on response time.
func (e *tcpCheckExecutor) determineFailureReason(responseTime time.Duration) *models.FailureReason {
	// Check timeouts
	failedThreshold := e.check.FailedThresholdDuration()
	if failedThreshold > 0 && responseTime > failedThreshold {
		return failureReasonPtr(models.FailureConnectionTimeout)
	}

	// Default to unknown
	return failureReasonPtr(models.FailureUnknown)
}
