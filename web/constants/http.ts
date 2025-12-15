export const STATUS_CODES = [
  100, // Continue
  101, // Switching Protocols
  102, // Processing
  103, // Early Hints
  200, // OK
  201, // Created
  202, // Accepted
  203, // Non-Authoritative Information
  204, // No Content
  205, // Reset Content
  206, // Partial Content
  207, // Multi-Status
  208, // Already Reported
  226, // IM Used
  300, // Multiple Choices
  301, // Moved Permanently
  302, // Found
  303, // See Other
  304, // Not Modified
  307, // Temporary Redirect
  308, // Permanent Redirect
  400, // Bad Request
  401, // Unauthorized
  402, // Payment Required
  403, // Forbidden
  404, // Not Found
  405, // Method Not Allowed
  406, // Not Acceptable
  407, // Proxy Authentication Required
  408, // Request Timeout
  409, // Conflict
  410, // Gone
  411, // Length Required
  412, // Precondition Failed
  413, // Content Too Large
  414, // URI Too Long
  415, // Unsupported Media Type
  416, // Range Not Satisfiable
  417, // Expectation Failed
  418, // I'm a teapot
  421, // Misdirected Request
  422, // Unprocessable Content
  423, // Locked
  424, // Failed Dependency
  425, // Too Early
  426, // Upgrade Required
  428, // Precondition Required
  429, // Too Many Requests
  431, // Request Header Fields Too Large
  451, // Unavailable For Legal Reasons
  500, // Internal Server Error
  501, // Not Implemented
  502, // Bad Gateway
  503, // Service Unavailable
  504, // Gateway Timeout
  505, // HTTP Version Not Supported
  506, // Variant Also Negotiates
  507, // Insufficient Storage
  508, // Loop Detected
  510, // Not Extended
  511, // Network Authentication Required
] as const

export const HTTP_METHODS = [
  'GET',
  'POST',
  'PUT',
  'DELETE',
  'PATCH',
  'HEAD',
  'OPTIONS',
] as const

export const IP_VERSIONS = [
  'ipv4',
  'ipv6',
] as const

export const DNS_RECORD_TYPES = [
  'A',
  'AAAA',
  'CNAME',
  'MX',
  'NS',
  'SOA',
  'SRV',
  'TXT',
] as const

export const DNS_RESOLVER_PROTOCOLS = [
  'udp',
  'tcp',
] as const

export const ASSERTION_SOURCES = [
  'status_code',
  'response_body_json',
  'response_body_text',
  'response_headers',
  'response_time_ms',
] as const

export const ASSERTION_COMPARISONS = [
  'equals',
  'not_equals',
  'contains',
  'not_contains',
  'is_empty',
  'is_not_empty',
  'is_less_than',
  'is_less_than_or_equal_to',
  'is_greater_than',
  'is_greater_than_or_equal_to',
] as const

export const ASSERTION_PROPERTIES = {
  status_code: {
    label: 'Status Code',
    type: 'number',
    operators: [
      'equals',
      'not_equals',
      'is_less_than',
      'is_less_than_or_equal_to',
      'is_greater_than',
      'is_greater_than_or_equal_to',
    ],
  },
  response_body_json: {
    label: 'JSON Response',
    type: 'string',
    operators: [
      'equals',
      'not_equals',
      'contains',
      'not_contains',
      'is_empty',
      'is_not_empty',
    ],
  },
  response_body_text: {
    label: 'Text Response',
    type: 'string',
    operators: [
      'equals',
      'not_equals',
      'contains',
      'not_contains',
      'is_empty',
      'is_not_empty',
    ],
  },
  response_headers: {
    label: 'Response Headers',
    type: 'object',
    operators: [
      'equals',
      'not_equals',
      'contains',
      'not_contains',
      'is_empty',
      'is_not_empty',
    ],
  },
  response_time_ms: {
    label: 'Response Time (ms)',
    type: 'number',
    operators: [
      'equals',
      'not_equals',
      'is_less_than',
      'is_less_than_or_equal_to',
      'is_greater_than',
      'is_greater_than_or_equal_to',
    ],
  },
}
