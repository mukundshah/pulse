package checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
)

// AssertionSource represents the source of data for an assertion.
type AssertionSource string

const (
	AssertionSourceStatusCode       AssertionSource = "status_code"
	AssertionSourceResponseTimeMs   AssertionSource = "response_time_ms"
	AssertionSourceResponseBodyText AssertionSource = "response_body_text"
	AssertionSourceResponseBodyJSON AssertionSource = "response_body_json"
	AssertionSourceResponseHeaders  AssertionSource = "response_headers"
)

// AssertionComparison represents the comparison operation to perform.
type AssertionComparison string

const (
	AssertionComparisonEquals                 AssertionComparison = "equals"
	AssertionComparisonNotEquals              AssertionComparison = "not_equals"
	AssertionComparisonContains               AssertionComparison = "contains"
	AssertionComparisonNotContains            AssertionComparison = "not_contains"
	AssertionComparisonIsEmpty                AssertionComparison = "is_empty"
	AssertionComparisonIsNotEmpty             AssertionComparison = "is_not_empty"
	AssertionComparisonIsLessThan             AssertionComparison = "is_less_than"
	AssertionComparisonIsLessThanOrEqualTo    AssertionComparison = "is_less_than_or_equal_to"
	AssertionComparisonIsGreaterThan          AssertionComparison = "is_greater_than"
	AssertionComparisonIsGreaterThanOrEqualTo AssertionComparison = "is_greater_than_or_equal_to"
)

var (
	// ErrEmptyBody is returned when the response body is empty.
	ErrEmptyBody = errors.New("empty response body")
	// ErrInvalidPath is returned when a JSON path is invalid.
	ErrInvalidPath = errors.New("invalid path")
	// ErrPathNotFound is returned when a path cannot be resolved.
	ErrPathNotFound = errors.New("path not found in data")
)

// Assertion defines a single assertion to be evaluated.
type Assertion struct {
	Source     AssertionSource `json:"source"`
	Property   *string         `json:"property,omitempty"`
	Comparison string          `json:"comparison"`
	Target     string          `json:"target"`
}

// AssertionResult contains the result of evaluating an assertion.
type AssertionResult struct {
	Assertion
	Received interface{} `json:"received"`
	Passed   bool        `json:"passed"`
}

// responseContext holds cached response data to avoid re-reading.
type responseContext struct {
	resp         *http.Response
	responseTime time.Duration
	bodyText     string
	bodyJSON     interface{}
	bodyRead     bool
	jsonParsed   bool
	jsonErr      error
}

// ProcessAssertions evaluates a list of assertions against an HTTP response.
func ProcessAssertions(assertions datatypes.JSON, resp *http.Response, responseTime time.Duration) ([]AssertionResult, error) {
	var assertionsList []Assertion
	if err := json.Unmarshal(assertions, &assertionsList); err != nil {
		return nil, fmt.Errorf("failed to unmarshal assertions: %w", err)
	}

	ctx := &responseContext{
		resp:         resp,
		responseTime: responseTime,
	}

	results := make([]AssertionResult, len(assertionsList))
	for i, assertion := range assertionsList {
		results[i] = ctx.evaluateAssertion(assertion)
	}

	return results, nil
}

// readBody reads and caches the response body.
func (rc *responseContext) readBody() error {
	if rc.bodyRead {
		return nil
	}
	rc.bodyRead = true

	body, err := io.ReadAll(rc.resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	rc.bodyText = string(body)
	return nil
}

// getBodyJSON parses and caches the response body as JSON.
func (rc *responseContext) getBodyJSON() (interface{}, error) {
	if rc.jsonParsed {
		return rc.bodyJSON, rc.jsonErr
	}
	rc.jsonParsed = true

	if err := rc.readBody(); err != nil {
		rc.jsonErr = err
		return nil, rc.jsonErr
	}

	if rc.bodyText == "" {
		rc.jsonErr = ErrEmptyBody
		return nil, rc.jsonErr
	}

	if err := json.Unmarshal([]byte(rc.bodyText), &rc.bodyJSON); err != nil {
		rc.jsonErr = fmt.Errorf("failed to parse JSON: %w", err)
		return nil, rc.jsonErr
	}

	return rc.bodyJSON, nil
}

// evaluateAssertion evaluates a single assertion.
func (rc *responseContext) evaluateAssertion(a Assertion) AssertionResult {
	result := AssertionResult{Assertion: a}

	switch a.Source {
	case AssertionSourceStatusCode:
		result.Received = rc.resp.StatusCode
		result.Passed = evaluateNumber(a.Comparison, float64(rc.resp.StatusCode), a.Target)

	case AssertionSourceResponseTimeMs:
		responseTimeMs := int(rc.responseTime / time.Millisecond)
		result.Received = responseTimeMs
		result.Passed = evaluateNumber(a.Comparison, float64(responseTimeMs), a.Target)

	case AssertionSourceResponseBodyText:
		if err := rc.readBody(); err != nil {
			result.Received = nil
			result.Passed = false
			return result
		}
		result.Received = rc.bodyText
		result.Passed = evaluateString(a.Comparison, rc.bodyText, a.Target)

	case AssertionSourceResponseBodyJSON:
		bodyJSON, err := rc.getBodyJSON()
		if err != nil {
			result.Received = nil
			result.Passed = false
			return result
		}

		value, err := resolvePath(bodyJSON, a.Property)
		if err != nil {
			result.Received = nil
			result.Passed = false
			return result
		}
		result.Received = value
		result.Passed = evaluateDynamic(a.Comparison, value, a.Target)

	case AssertionSourceResponseHeaders:
		headersMap := headersToMap(rc.resp.Header)
		value, err := resolvePath(headersMap, a.Property)
		if err != nil {
			result.Received = nil
			result.Passed = false
			return result
		}
		result.Received = value
		result.Passed = evaluateDynamic(a.Comparison, value, a.Target)

	default:
		result.Received = nil
		result.Passed = false
	}

	return result
}

// evaluateNumber performs numeric comparisons.
func evaluateNumber(comparison string, actual float64, expected interface{}) bool {
	expectedNum, ok := toFloat(expected)
	if !ok {
		return false
	}

	switch AssertionComparison(comparison) {
	case AssertionComparisonEquals:
		return actual == expectedNum
	case AssertionComparisonNotEquals:
		return actual != expectedNum
	case AssertionComparisonIsLessThan:
		return actual < expectedNum
	case AssertionComparisonIsLessThanOrEqualTo:
		return actual <= expectedNum
	case AssertionComparisonIsGreaterThan:
		return actual > expectedNum
	case AssertionComparisonIsGreaterThanOrEqualTo:
		return actual >= expectedNum
	default:
		return false
	}
}

// evaluateString performs string comparisons.
func evaluateString(comparison string, actual string, expected interface{}) bool {
	expectedStr := fmt.Sprint(expected)

	switch AssertionComparison(comparison) {
	case AssertionComparisonEquals:
		return actual == expectedStr
	case AssertionComparisonNotEquals:
		return actual != expectedStr
	case AssertionComparisonContains:
		return strings.Contains(actual, expectedStr)
	case AssertionComparisonNotContains:
		return !strings.Contains(actual, expectedStr)
	case AssertionComparisonIsEmpty:
		return actual == ""
	case AssertionComparisonIsNotEmpty:
		return actual != ""
	default:
		return false
	}
}

// evaluateDynamic performs comparisons on values of unknown types.
func evaluateDynamic(comparison string, actual, expected interface{}) bool {
	comp := AssertionComparison(comparison)

	// Handle numeric comparisons
	if isNumericComparison(comp) {
		actualNum, okActual := toFloat(actual)
		expectedNum, okExpected := toFloat(expected)
		if !okActual || !okExpected {
			return false
		}
		return evaluateNumber(comparison, actualNum, expectedNum)
	}

	switch comp {
	case AssertionComparisonEquals:
		return compareEquals(actual, expected)
	case AssertionComparisonNotEquals:
		return !compareEquals(actual, expected)
	case AssertionComparisonContains:
		return containsValue(actual, expected)
	case AssertionComparisonNotContains:
		return !containsValue(actual, expected)
	case AssertionComparisonIsEmpty:
		return isEmpty(actual)
	case AssertionComparisonIsNotEmpty:
		return !isEmpty(actual)
	default:
		return false
	}
}

// compareEquals checks equality with numeric type coercion.
func compareEquals(actual, expected interface{}) bool {
	// Try numeric comparison first
	if actualNum, okActual := toFloat(actual); okActual {
		if expectedNum, okExpected := toFloat(expected); okExpected {
			return actualNum == expectedNum
		}
	}
	return reflect.DeepEqual(actual, expected)
}

// toFloat converts various numeric types to float64.
func toFloat(v interface{}) (float64, bool) {
	switch num := v.(type) {
	case float64:
		return num, true
	case float32:
		return float64(num), true
	case int:
		return float64(num), true
	case int64:
		return float64(num), true
	case int32:
		return float64(num), true
	case int16:
		return float64(num), true
	case int8:
		return float64(num), true
	case uint:
		return float64(num), true
	case uint64:
		return float64(num), true
	case uint32:
		return float64(num), true
	case uint16:
		return float64(num), true
	case uint8:
		return float64(num), true
	case json.Number:
		f, err := num.Float64()
		return f, err == nil
	case string:
		// Try parsing string as number
		f, err := strconv.ParseFloat(num, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

// isNumericComparison checks if the comparison requires numeric operands.
func isNumericComparison(comparison AssertionComparison) bool {
	switch comparison {
	case AssertionComparisonIsLessThan,
		AssertionComparisonIsLessThanOrEqualTo,
		AssertionComparisonIsGreaterThan,
		AssertionComparisonIsGreaterThanOrEqualTo:
		return true
	default:
		return false
	}
}

// containsValue checks if actual contains expected.
func containsValue(actual, expected interface{}) bool {
	expectedStr := fmt.Sprint(expected)

	switch v := actual.(type) {
	case string:
		return strings.Contains(v, expectedStr)
	case []interface{}:
		for _, item := range v {
			if reflect.DeepEqual(item, expected) || fmt.Sprint(item) == expectedStr {
				return true
			}
		}
		return false
	default:
		return strings.Contains(fmt.Sprint(actual), expectedStr)
	}
}

// isEmpty checks if a value is considered empty.
func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case map[string]interface{}:
		return len(v) == 0
	default:
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map, reflect.Array, reflect.String:
			return rv.Len() == 0
		}
		return false
	}
}

// headersToMap converts HTTP headers to a map with lowercase keys.
func headersToMap(headers http.Header) map[string]interface{} {
	result := make(map[string]interface{}, len(headers))
	for k, v := range headers {
		lower := strings.ToLower(k)
		if len(v) == 1 {
			result[lower] = v[0]
		} else {
			result[lower] = v
		}
	}
	return result
}

// pathSegment represents a single segment in a JSON path.
type pathSegment struct {
	key   string
	index *int
}

// resolvePath resolves a path in a nested data structure.
func resolvePath(input interface{}, path *string) (interface{}, error) {
	if path == nil || *path == "" {
		return nil, ErrInvalidPath
	}

	segments, err := parsePath(*path)
	if err != nil {
		return nil, err
	}

	current := input
	for _, seg := range segments {
		// Resolve map key
		if seg.key != "" {
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("%w: expected map at key %s", ErrPathNotFound, seg.key)
			}

			val, exists := m[seg.key]
			if !exists {
				// Try case-insensitive lookup
				val, exists = m[strings.ToLower(seg.key)]
				if !exists {
					return nil, fmt.Errorf("%w: key %s not found", ErrPathNotFound, seg.key)
				}
			}
			current = val
		}

		// Resolve array index
		if seg.index != nil {
			arr, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("%w: expected array at index %d", ErrPathNotFound, *seg.index)
			}

			idx := *seg.index
			if idx < 0 || idx >= len(arr) {
				return nil, fmt.Errorf("%w: index %d out of bounds (len=%d)", ErrPathNotFound, idx, len(arr))
			}
			current = arr[idx]
		}
	}

	return current, nil
}

// parsePath parses a dot-notation path with array indices.
// Examples: "user.name", "users[0].email", "data.items[0][1]"
func parsePath(path string) ([]pathSegment, error) {
	if path == "" {
		return nil, ErrInvalidPath
	}

	rawSegments := strings.Split(path, ".")
	segments := make([]pathSegment, 0, len(rawSegments))

	for _, raw := range rawSegments {
		if raw == "" {
			return nil, fmt.Errorf("%w: empty segment in path", ErrInvalidPath)
		}

		remain := raw
		for remain != "" {
			seg := pathSegment{}

			// Find bracket position
			bracket := strings.Index(remain, "[")
			if bracket == -1 {
				// No brackets, entire remaining string is a key
				seg.key = remain
				remain = ""
			} else {
				// Extract key before bracket (if any)
				if bracket > 0 {
					seg.key = remain[:bracket]
				}
				remain = remain[bracket:]
			}

			// Add key segment if present
			if seg.key != "" {
				segments = append(segments, seg)
			}

			// Process array indices
			for strings.HasPrefix(remain, "[") {
				end := strings.Index(remain, "]")
				if end == -1 {
					return nil, fmt.Errorf("%w: unmatched bracket in path", ErrInvalidPath)
				}

				idxStr := remain[1:end]
				idx, err := strconv.Atoi(idxStr)
				if err != nil {
					return nil, fmt.Errorf("%w: invalid array index %q: %v", ErrInvalidPath, idxStr, err)
				}

				segments = append(segments, pathSegment{index: &idx})
				remain = remain[end+1:]
			}
		}
	}

	return segments, nil
}
