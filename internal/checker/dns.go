package checker

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"pulse/internal/models"

	"github.com/miekg/dns"
)

// dnsCheckExecutor executes DNS checks with all necessary configuration.
type dnsCheckExecutor struct {
	check     *models.Check
	timings   *dnsTimingTracker
	ipVersion string
	ipAddress string
	records   interface{} // Store the DNS records for result
	client    *dns.Client
	server    string   // DNS server address
	reply     *dns.Msg // Store full DNS response
	request   *dns.Msg // Store DNS request
}

// dnsTimingTracker tracks DNS query timing events.
type dnsTimingTracker struct {
	requestStart time.Time
	queryStart   time.Time
	queryDone    time.Time
	responseEnd  time.Time
}

// ExecuteDNSCheck performs a DNS check and returns the result.
// It handles DNS queries, timing tracking, and result building.
func ExecuteDNSCheck(ctx context.Context, check *models.Check) Result {
	executor := newDNSCheckExecutor(check)
	return executor.execute(ctx)
}

// newDNSCheckExecutor creates a new DNS check executor.
func newDNSCheckExecutor(check *models.Check) *dnsCheckExecutor {
	executor := &dnsCheckExecutor{
		check:   check,
		timings: &dnsTimingTracker{},
	}

	// Configure DNS client and server
	executor.configureDNSClient()

	return executor
}

// configureDNSClient configures the DNS client and server based on check settings.
func (e *dnsCheckExecutor) configureDNSClient() {
	// Determine DNS server
	if e.check.DNSResolver != nil && *e.check.DNSResolver != "" {
		// Use custom DNS resolver
		resolver := *e.check.DNSResolver
		port := 53 // Default DNS port

		if e.check.DNSResolverPort != nil && *e.check.DNSResolverPort > 0 {
			port = *e.check.DNSResolverPort
		}

		e.server = net.JoinHostPort(resolver, fmt.Sprintf("%d", port))
	} else {
		// Use system default DNS resolver
		// Try to get system DNS from /etc/resolv.conf or use 8.8.8.8 as fallback
		e.server = "8.8.8.8:53"
	}

	// Determine protocol (UDP or TCP)
	protocol := "udp" // Default to UDP
	if e.check.DNSResolverProtocol != nil {
		if *e.check.DNSResolverProtocol == models.DNSResolverProtocolTCP {
			protocol = "tcp"
		}
	}

	// Create DNS client with timeout
	timeout := defaultTimeout
	if e.check.FailedThresholdDuration() > 0 {
		timeout = e.check.FailedThresholdDuration()
	}

	e.client = &dns.Client{
		Net:     protocol,
		Timeout: timeout,
	}
}

// execute runs the DNS check and returns the result.
func (e *dnsCheckExecutor) execute(ctx context.Context) Result {
	// Start timer before query
	e.timings.requestStart = time.Now().UTC()

	// Validate DNS record type
	if e.check.DNSRecordType == nil {
		return e.createErrorResult(fmt.Errorf("DNS record type is required"))
	}

	// Perform DNS query
	err := e.performQuery(ctx)
	if err != nil {
		return e.createErrorResult(err)
	}

	// Query successful - record end time
	e.timings.responseEnd = time.Now().UTC()

	// Build result
	return e.buildResult()
}

// performQuery performs the DNS query based on record type using a single method.
func (e *dnsCheckExecutor) performQuery(ctx context.Context) error {
	host := e.check.Host

	// Track query timing
	e.timings.queryStart = time.Now().UTC()

	// Perform query using single method
	var err error
	e.records, err = e.queryDNS(ctx, *e.check.DNSRecordType, host)

	e.timings.queryDone = time.Now().UTC()

	if err != nil {
		return fmt.Errorf("DNS query failed: %w", err)
	}

	// Extract IP information if applicable
	e.extractIPInfo()

	return nil
}

// queryDNS performs a DNS query for the specified record type.
func (e *dnsCheckExecutor) queryDNS(ctx context.Context, recordType models.DNSRecordType, host string) (interface{}, error) {
	// Create DNS message
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(host), e.dnsRecordTypeToQType(recordType))
	msg.RecursionDesired = true

	// Create context with timeout if needed
	queryCtx := ctx
	if e.client.Timeout > 0 {
		var cancel context.CancelFunc
		queryCtx, cancel = context.WithTimeout(ctx, e.client.Timeout)
		defer cancel()
	}

	// Store request
	e.request = msg

	// Perform DNS query
	reply, _, err := e.client.ExchangeContext(queryCtx, msg, e.server)
	if err != nil {
		return nil, fmt.Errorf("DNS exchange failed: %w", err)
	}

	// Store full reply
	e.reply = reply

	// Check for DNS errors
	if reply.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS query returned error code: %s", dns.RcodeToString[reply.Rcode])
	}

	// Parse response based on record type
	return e.parseDNSResponse(reply, recordType, host)
}

// dnsRecordTypeToQType converts models.DNSRecordType to dns.QType.
func (e *dnsCheckExecutor) dnsRecordTypeToQType(recordType models.DNSRecordType) uint16 {
	switch recordType {
	case models.DNSRecordTypeA:
		return dns.TypeA
	case models.DNSRecordTypeAAAA:
		return dns.TypeAAAA
	case models.DNSRecordTypeCNAME:
		return dns.TypeCNAME
	case models.DNSRecordTypeMX:
		return dns.TypeMX
	case models.DNSRecordTypeNS:
		return dns.TypeNS
	case models.DNSRecordTypeTXT:
		return dns.TypeTXT
	case models.DNSRecordTypeSRV:
		return dns.TypeSRV
	case models.DNSRecordTypeSOA:
		return dns.TypeSOA
	default:
		return dns.TypeA // Default fallback
	}
}

// parseDNSResponse parses DNS response based on record type.
func (e *dnsCheckExecutor) parseDNSResponse(reply *dns.Msg, recordType models.DNSRecordType, host string) (interface{}, error) {
	switch recordType {
	case models.DNSRecordTypeA:
		return e.parseARecords(reply, host)
	case models.DNSRecordTypeAAAA:
		return e.parseAAAARecords(reply, host)
	case models.DNSRecordTypeCNAME:
		return e.parseCNAMERecords(reply, host)
	case models.DNSRecordTypeMX:
		return e.parseMXRecords(reply, host)
	case models.DNSRecordTypeNS:
		return e.parseNSRecords(reply, host)
	case models.DNSRecordTypeTXT:
		return e.parseTXTRecords(reply, host)
	case models.DNSRecordTypeSRV:
		return e.parseSRVRecords(reply, host)
	case models.DNSRecordTypeSOA:
		return e.parseSOARecords(reply, host)
	default:
		return nil, fmt.Errorf("unsupported DNS record type: %s", recordType)
	}
}

// parseARecords parses A records from DNS response.
func (e *dnsCheckExecutor) parseARecords(reply *dns.Msg, host string) (interface{}, error) {
	var ipv4Addrs []string
	for _, rr := range reply.Answer {
		if a, ok := rr.(*dns.A); ok {
			ipv4Addrs = append(ipv4Addrs, a.A.String())
		}
	}

	if len(ipv4Addrs) == 0 {
		return nil, fmt.Errorf("no A records found for %s", host)
	}

	return ipv4Addrs, nil
}

// parseAAAARecords parses AAAA records from DNS response.
func (e *dnsCheckExecutor) parseAAAARecords(reply *dns.Msg, host string) (interface{}, error) {
	var ipv6Addrs []string
	for _, rr := range reply.Answer {
		if aaaa, ok := rr.(*dns.AAAA); ok {
			ipv6Addrs = append(ipv6Addrs, aaaa.AAAA.String())
		}
	}

	if len(ipv6Addrs) == 0 {
		return nil, fmt.Errorf("no AAAA records found for %s", host)
	}

	return ipv6Addrs, nil
}

// parseCNAMERecords parses CNAME records from DNS response.
func (e *dnsCheckExecutor) parseCNAMERecords(reply *dns.Msg, host string) (interface{}, error) {
	for _, rr := range reply.Answer {
		if cname, ok := rr.(*dns.CNAME); ok {
			return cname.Target, nil
		}
	}

	return nil, fmt.Errorf("no CNAME record found for %s", host)
}

// parseMXRecords parses MX records from DNS response.
func (e *dnsCheckExecutor) parseMXRecords(reply *dns.Msg, host string) (interface{}, error) {
	var mxList []map[string]interface{}
	for _, rr := range reply.Answer {
		if mx, ok := rr.(*dns.MX); ok {
			mxList = append(mxList, map[string]interface{}{
				"host": mx.Mx,
				"pref": mx.Preference,
			})
		}
	}

	if len(mxList) == 0 {
		return nil, fmt.Errorf("no MX records found for %s", host)
	}

	return mxList, nil
}

// parseNSRecords parses NS records from DNS response.
func (e *dnsCheckExecutor) parseNSRecords(reply *dns.Msg, host string) (interface{}, error) {
	var nsList []string
	for _, rr := range reply.Answer {
		if ns, ok := rr.(*dns.NS); ok {
			nsList = append(nsList, ns.Ns)
		}
	}

	if len(nsList) == 0 {
		return nil, fmt.Errorf("no NS records found for %s", host)
	}

	return nsList, nil
}

// parseTXTRecords parses TXT records from DNS response.
func (e *dnsCheckExecutor) parseTXTRecords(reply *dns.Msg, host string) (interface{}, error) {
	var txtList []string
	for _, rr := range reply.Answer {
		if txt, ok := rr.(*dns.TXT); ok {
			// TXT records can have multiple strings, join them
			txtList = append(txtList, txt.Txt...)
		}
	}

	if len(txtList) == 0 {
		return nil, fmt.Errorf("no TXT records found for %s", host)
	}

	return txtList, nil
}

// parseSRVRecords parses SRV records from DNS response.
func (e *dnsCheckExecutor) parseSRVRecords(reply *dns.Msg, host string) (interface{}, error) {
	var srvList []map[string]interface{}
	for _, rr := range reply.Answer {
		if srv, ok := rr.(*dns.SRV); ok {
			srvList = append(srvList, map[string]interface{}{
				"target":   srv.Target,
				"port":     srv.Port,
				"priority": srv.Priority,
				"weight":   srv.Weight,
			})
		}
	}

	if len(srvList) == 0 {
		return nil, fmt.Errorf("no SRV records found for %s", host)
	}

	return srvList, nil
}

// parseSOARecords parses SOA records from DNS response.
func (e *dnsCheckExecutor) parseSOARecords(reply *dns.Msg, host string) (interface{}, error) {
	for _, rr := range reply.Answer {
		if soa, ok := rr.(*dns.SOA); ok {
			return map[string]interface{}{
				"ns":      soa.Ns,
				"mbox":    soa.Mbox,
				"serial":  soa.Serial,
				"refresh": soa.Refresh,
				"retry":   soa.Retry,
				"expire":  soa.Expire,
				"minttl":  soa.Minttl,
			}, nil
		}
	}

	return nil, fmt.Errorf("no SOA record found for %s", host)
}

// extractIPInfo extracts IP version and address from DNS records if applicable.
func (e *dnsCheckExecutor) extractIPInfo() {
	if e.records == nil {
		return
	}

	// Extract IP info for A/AAAA records
	switch *e.check.DNSRecordType {
	case models.DNSRecordTypeA, models.DNSRecordTypeAAAA:
		if ipList, ok := e.records.([]string); ok && len(ipList) > 0 {
			// Use first IP address
			e.ipAddress = ipList[0]
			ip := net.ParseIP(e.ipAddress)
			if ip != nil {
				if ip.To4() != nil {
					e.ipVersion = "IPv4"
				} else {
					e.ipVersion = "IPv6"
				}
			}
		}
	}
}

// buildNetworkTimings creates a map of detailed network timing information.
func (e *dnsCheckExecutor) buildNetworkTimings() map[string]interface{} {
	timings := make(map[string]interface{})

	// Store raw timestamps
	if !e.timings.requestStart.IsZero() {
		timings["request_start"] = e.timings.requestStart.Format(time.RFC3339Nano)
	}
	if !e.timings.queryStart.IsZero() {
		timings["query_start"] = e.timings.queryStart.Format(time.RFC3339Nano)
	}
	if !e.timings.queryDone.IsZero() {
		timings["query_done"] = e.timings.queryDone.Format(time.RFC3339Nano)
	}
	if !e.timings.responseEnd.IsZero() {
		timings["response_end"] = e.timings.responseEnd.Format(time.RFC3339Nano)
	}

	// Compute durations (in microseconds)
	// Query duration: query_done - query_start
	if !e.timings.queryStart.IsZero() && !e.timings.queryDone.IsZero() && e.timings.queryDone.After(e.timings.queryStart) {
		if us := durationUs(e.timings.queryStart, e.timings.queryDone); us > 0 {
			timings["query_duration_us"] = us
		}
	}
	// Total response time: response_end - request_start
	if !e.timings.requestStart.IsZero() && !e.timings.responseEnd.IsZero() && e.timings.responseEnd.After(e.timings.requestStart) {
		responseTime := e.responseTime()
		if responseTime > 0 {
			timings["response_time_us"] = int(responseTime / time.Microsecond)
		}
	}

	// Add DNS server information
	if e.server != "" {
		timings["dns_server"] = e.server
	}

	// Add DNS records to network timings for visibility
	if e.records != nil {
		timings["records"] = e.records
	}

	return timings
}

// responseTime calculates the total response time from timestamps.
func (e *dnsCheckExecutor) responseTime() time.Duration {
	if e.timings.requestStart.IsZero() || e.timings.responseEnd.IsZero() {
		return 0
	}
	return e.timings.responseEnd.Sub(e.timings.requestStart)
}

// determineStatus calculates the check status based on response time and thresholds.
func (e *dnsCheckExecutor) determineStatus(responseTime time.Duration) models.CheckRunStatus {
	// Check response time thresholds
	if responseTime > e.check.FailedThresholdDuration() {
		return models.CheckRunStatusFailing
	}
	if responseTime > e.check.DegradedThresholdDuration() {
		return models.CheckRunStatusDegraded
	}

	return models.CheckRunStatusPassing
}

// buildDNSJSONResponse builds a structured JSON representation of the DNS response.
func (e *dnsCheckExecutor) buildDNSJSONResponse() map[string]interface{} {
	jsonResp := map[string]interface{}{
		"Status": dns.RcodeToString[e.reply.Rcode],
		"TC":     e.reply.Truncated,
		"AD":     e.reply.AuthenticatedData,
		"CD":     e.reply.CheckingDisabled,
		"ID":     e.reply.Id,
	}

	// Question section
	questions := make([]map[string]interface{}, 0, len(e.reply.Question))
	for _, q := range e.reply.Question {
		questions = append(questions, map[string]interface{}{
			"Name": q.Name,
			"Type": dns.TypeToString[q.Qtype],
		})
	}
	jsonResp["Question"] = questions

	// Answer section
	answers := make([]map[string]interface{}, 0)
	for _, rr := range e.reply.Answer {
		answer := map[string]interface{}{
			"name": rr.Header().Name,
			"type": dns.TypeToString[rr.Header().Rrtype],
			"TTL":  rr.Header().Ttl,
		}

		// Extract data based on record type
		switch v := rr.(type) {
		case *dns.A:
			answer["data"] = v.A.String()
		case *dns.AAAA:
			answer["data"] = v.AAAA.String()
		case *dns.CNAME:
			answer["data"] = v.Target
		case *dns.MX:
			answer["data"] = v.Mx
			answer["preference"] = v.Preference
		case *dns.NS:
			answer["data"] = v.Ns
		case *dns.TXT:
			answer["data"] = v.Txt
		case *dns.SRV:
			answer["data"] = v.Target
			answer["port"] = v.Port
			answer["priority"] = v.Priority
			answer["weight"] = v.Weight
		case *dns.SOA:
			answer["data"] = map[string]interface{}{
				"ns":     v.Ns,
				"mbox":   v.Mbox,
				"serial": v.Serial,
			}
		default:
			answer["data"] = rr.String()
		}

		answers = append(answers, answer)
	}
	jsonResp["Answer"] = answers

	// Authority section
	authorities := make([]map[string]interface{}, 0)
	for _, rr := range e.reply.Ns {
		authority := map[string]interface{}{
			"name": rr.Header().Name,
			"type": dns.TypeToString[rr.Header().Rrtype],
			"TTL":  rr.Header().Ttl,
		}
		authorities = append(authorities, authority)
	}
	if len(authorities) > 0 {
		jsonResp["Authority"] = authorities
	}

	// Additional section
	additionals := make([]map[string]interface{}, 0)
	for _, rr := range e.reply.Extra {
		additional := map[string]interface{}{
			"name": rr.Header().Name,
			"type": dns.TypeToString[rr.Header().Rrtype],
			"TTL":  rr.Header().Ttl,
		}
		additionals = append(additionals, additional)
	}
	if len(additionals) > 0 {
		jsonResp["Additional"] = additionals
	}

	return jsonResp
}

// buildDNSRawResponse builds a raw (dig-like) text representation of the DNS response.
// Returns the raw format as a string for storage in response.formats.raw
func (e *dnsCheckExecutor) buildDNSRawResponse() interface{} {
	if e.reply == nil {
		return ""
	}
	if e.reply == nil {
		return ""
	}

	var buf strings.Builder

	// Header line
	buf.WriteString(fmt.Sprintf(";; opcode: %s, status: %s, id: %d\n",
		dns.OpcodeToString[e.reply.Opcode],
		dns.RcodeToString[e.reply.Rcode],
		e.reply.Id))

	// Flags
	flags := []string{}
	if e.reply.Response {
		flags = append(flags, "qr")
	}
	if e.reply.Authoritative {
		flags = append(flags, "aa")
	}
	if e.reply.Truncated {
		flags = append(flags, "tc")
	}
	if e.reply.RecursionDesired {
		flags = append(flags, "rd")
	}
	if e.reply.RecursionAvailable {
		flags = append(flags, "ra")
	}
	if e.reply.Zero {
		flags = append(flags, "z")
	}
	if e.reply.AuthenticatedData {
		flags = append(flags, "ad")
	}
	if e.reply.CheckingDisabled {
		flags = append(flags, "cd")
	}

	buf.WriteString(fmt.Sprintf(";; flags: %s; QUERY: %d, ANSWER: %d, AUTHORITY: %d, ADDITIONAL: %d\n\n",
		strings.Join(flags, " "),
		len(e.reply.Question),
		len(e.reply.Answer),
		len(e.reply.Ns),
		len(e.reply.Extra)))

	// Question section
	if len(e.reply.Question) > 0 {
		buf.WriteString(";; QUESTION SECTION:\n")
		for _, q := range e.reply.Question {
			buf.WriteString(fmt.Sprintf(";%s\t\tIN\t%s\n", q.Name, dns.TypeToString[q.Qtype]))
		}
		buf.WriteString("\n")
	}

	// Answer section
	if len(e.reply.Answer) > 0 {
		buf.WriteString(";; ANSWER SECTION:\n")
		for _, rr := range e.reply.Answer {
			buf.WriteString(fmt.Sprintf("%s\t%d\tIN\t%s\t%s\n",
				rr.Header().Name,
				rr.Header().Ttl,
				dns.TypeToString[rr.Header().Rrtype],
				e.formatRRData(rr)))
		}
		buf.WriteString("\n")
	}

	// Authority section
	if len(e.reply.Ns) > 0 {
		buf.WriteString(";; AUTHORITY SECTION:\n")
		for _, rr := range e.reply.Ns {
			buf.WriteString(fmt.Sprintf("%s\t%d\tIN\t%s\t%s\n",
				rr.Header().Name,
				rr.Header().Ttl,
				dns.TypeToString[rr.Header().Rrtype],
				e.formatRRData(rr)))
		}
		buf.WriteString("\n")
	}

	// Additional section
	if len(e.reply.Extra) > 0 {
		buf.WriteString(";; ADDITIONAL SECTION:\n")
		for _, rr := range e.reply.Extra {
			buf.WriteString(fmt.Sprintf("%s\t%d\tIN\t%s\t%s\n",
				rr.Header().Name,
				rr.Header().Ttl,
				dns.TypeToString[rr.Header().Rrtype],
				e.formatRRData(rr)))
		}
	}

	return buf.String()
}

// formatRRData formats the data portion of a DNS resource record.
func (e *dnsCheckExecutor) formatRRData(rr dns.RR) string {
	switch v := rr.(type) {
	case *dns.A:
		return v.A.String()
	case *dns.AAAA:
		return v.AAAA.String()
	case *dns.CNAME:
		return v.Target
	case *dns.MX:
		return fmt.Sprintf("%d %s", v.Preference, v.Mx)
	case *dns.NS:
		return v.Ns
	case *dns.TXT:
		return strings.Join(v.Txt, " ")
	case *dns.SRV:
		return fmt.Sprintf("%d %d %d %s", v.Priority, v.Weight, v.Port, v.Target)
	case *dns.SOA:
		return fmt.Sprintf("%s %s %d %d %d %d %d",
			v.Ns, v.Mbox, v.Serial, v.Refresh, v.Retry, v.Expire, v.Minttl)
	default:
		return rr.String()
	}
}

// buildResult creates the final result object.
func (e *dnsCheckExecutor) buildResult() Result {
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

	// Build response data using unified builder
	rb := &ResponseBuilder{}
	rawFormat := e.buildDNSRawResponse()
	jsonFormat := e.buildDNSJSONResponse()
	responseData := rb.BuildDNSResponse(e.records, e.server, rawFormat, jsonFormat)

	return Result{
		Status:            status,
		FailureReason:     failureReason,
		ResponseStatus:    nil, // DNS doesn't have HTTP status codes
		RequestStartedAt:  e.timings.requestStart,
		FirstByteAt:       e.timings.queryDone, // For DNS, query done is like "first byte"
		ResponseEndedAt:   e.timings.responseEnd,
		ConnectionReused:  false, // DNS queries don't reuse connections
		IPVersion:         e.ipVersion,
		IPAddress:         e.ipAddress,
		ResponseSizeBytes: 0,                 // DNS queries don't have response size in bytes
		AssertionResults:  emptyJSONObject(), // DNS checks don't have assertions yet
		PlaywrightReport:  emptyJSONObject(),
		NetworkTimings:    mustMarshalJSON(networkTimings),
		Response:          responseData,
		Error:             nil,
	}
}

// createErrorResult creates a result for a failed check.
func (e *dnsCheckExecutor) createErrorResult(err error) Result {
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
		FirstByteAt:       e.timings.queryDone,   // May be zero
		ResponseEndedAt:   e.timings.responseEnd, // May be zero
		ConnectionReused:  false,
		IPVersion:         e.ipVersion,
		IPAddress:         e.ipAddress,
		ResponseSizeBytes: 0,
		AssertionResults:  emptyJSONObject(),
		PlaywrightReport:  emptyJSONObject(),
		NetworkTimings:    emptyJSONObject(),
		Response:          EmptyResponse(),
		Error:             err,
	}
}

// classifyError determines the failure reason from an error.
func (e *dnsCheckExecutor) classifyError(err error) *models.FailureReason {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// DNS-specific errors
	if contains(errStr, "no such host") || contains(errStr, "dns") {
		return failureReasonPtr(models.FailureDNS)
	}
	if contains(errStr, "no such record") || contains(errStr, "not found") {
		return failureReasonPtr(models.FailureDNS)
	}
	if contains(errStr, "timeout") || contains(errStr, "deadline exceeded") {
		return failureReasonPtr(models.FailureRequestTimeout)
	}
	if contains(errStr, "network is unreachable") {
		return failureReasonPtr(models.FailureNetworkUnreachable)
	}
	if contains(errStr, "NXDOMAIN") || contains(errStr, "SERVFAIL") || contains(errStr, "REFUSED") {
		return failureReasonPtr(models.FailureDNS)
	}

	return failureReasonPtr(models.FailureUnknown)
}

// determineFailureReason determines the failure reason based on response time.
func (e *dnsCheckExecutor) determineFailureReason(responseTime time.Duration) *models.FailureReason {
	// Check timeouts
	failedThreshold := e.check.FailedThresholdDuration()
	if failedThreshold > 0 && responseTime > failedThreshold {
		return failureReasonPtr(models.FailureRequestTimeout)
	}

	// Default to unknown
	return failureReasonPtr(models.FailureUnknown)
}
