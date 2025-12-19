package checker

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"

	"gorm.io/datatypes"
)

const (
	// MaxResponseBodySize is the maximum size of HTTP response body to store (1MB)
	MaxResponseBodySize = 1 * 1024 * 1024
	// MaxTextBodySize is the maximum size for text bodies before base64 encoding (500KB)
	MaxTextBodySize = 500 * 1024
)

// ResponseBuilder helps build uniform response structures across check types.
type ResponseBuilder struct{}

// BuildHTTPResponse builds a uniform HTTP response structure.
func (rb *ResponseBuilder) BuildHTTPResponse(headers map[string][]string, body []byte, contentType, proto string) datatypes.JSON {
	response := make(map[string]interface{})
	response["type"] = "http"
	response["headers"] = headers
	response["content_type"] = contentType
	response["proto"] = proto

	bodySize := len(body)
	response["body_size_bytes"] = bodySize

	// Handle body with size limits
	if bodySize == 0 {
		response["body"] = ""
		response["body_encoding"] = "text"
	} else if bodySize > MaxResponseBodySize {
		// Truncate large bodies
		truncatedBody := body[:MaxResponseBodySize]
		if rb.isTextContent(contentType) && utf8.Valid(truncatedBody) {
			response["body"] = string(truncatedBody)
			response["body_encoding"] = "text"
		} else {
			response["body"] = base64.StdEncoding.EncodeToString(truncatedBody)
			response["body_encoding"] = "base64"
		}
		response["body_truncated"] = true
		response["body_original_size_bytes"] = bodySize
	} else {
		// Store full body
		isText := rb.isTextContent(contentType)
		if isText && bodySize <= MaxTextBodySize && utf8.Valid(body) {
			response["body"] = string(body)
			response["body_encoding"] = "text"
		} else {
			response["body"] = base64.StdEncoding.EncodeToString(body)
			response["body_encoding"] = "base64"
		}
		response["body_truncated"] = false
	}

	return mustMarshalJSON(response)
}

// BuildTCPResponse builds a uniform TCP response structure.
func (rb *ResponseBuilder) BuildTCPResponse() datatypes.JSON {
	response := make(map[string]interface{})
	response["type"] = "tcp"
	response["connection_status"] = "established"

	return mustMarshalJSON(response)
}

// BuildDNSResponse builds a uniform DNS response structure.
func (rb *ResponseBuilder) BuildDNSResponse(records interface{}, dnsServer string, rawFormat interface{}, jsonFormat map[string]interface{}) datatypes.JSON {
	response := make(map[string]interface{})
	response["type"] = "dns"
	response["records"] = records
	response["dns_server"] = dnsServer

	// Add formatted outputs if available
	if rawFormat != nil || jsonFormat != nil {
		formats := make(map[string]interface{})
		if rawFormat != nil {
			formats["raw"] = rawFormat
		}
		if jsonFormat != nil {
			formats["json"] = jsonFormat
		}
		response["formats"] = formats
	}

	return mustMarshalJSON(response)
}

// isTextContent checks if the content type indicates text-based content.
func (rb *ResponseBuilder) isTextContent(contentType string) bool {
	if contentType == "" {
		return false
	}

	textTypes := []string{
		"text/",
		"application/json",
		"application/xml",
		"application/xhtml+xml",
		"application/javascript",
		"application/x-javascript",
		"application/ecmascript",
		"application/x-www-form-urlencoded",
		"application/xml-dtd",
		"application/xop+xml",
		"application/atom+xml",
		"application/rss+xml",
		"application/soap+xml",
		"application/xslt+xml",
		"application/mathml+xml",
		"application/svg+xml",
		"application/x-yaml",
		"application/yaml",
		"application/toml",
		"application/csv",
	}

	contentTypeLower := strings.ToLower(contentType)
	for _, textType := range textTypes {
		if strings.HasPrefix(contentTypeLower, textType) {
			return true
		}
	}

	return false
}

// EmptyResponse returns an empty response object.
func EmptyResponse() datatypes.JSON {
	return emptyJSONObject()
}
