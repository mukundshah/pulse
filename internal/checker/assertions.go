package checker

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gorm.io/datatypes"
)

func processAssertions(assertions datatypes.JSON, statusCode int, body []byte) (datatypes.JSON, bool) {
	if len(assertions) == 0 {
		return emptyJSON(), true
	}

	var assertionRules []map[string]interface{}
	if err := json.Unmarshal(assertions, &assertionRules); err != nil {
		return mustMarshalJSON(map[string]interface{}{"error": "invalid assertions format"}), false
	}

	results := make(map[string]interface{})
	allPassed := true

	for i, rule := range assertionRules {
		ruleType, _ := rule["type"].(string)
		passed := false

		switch ruleType {
		case "status_code":
			expected, _ := rule["value"].(float64)
			passed = statusCode == int(expected)
		case "body_contains":
			expected, _ := rule["value"].(string)
			passed = bytes.Contains(body, []byte(expected))
		case "body_not_contains":
			expected, _ := rule["value"].(string)
			passed = !bytes.Contains(body, []byte(expected))
		case "response_time":
			// Requires explicit latency input; not implemented yet
			passed = false
		}

		if !passed {
			allPassed = false
		}

		results[fmt.Sprintf("assertion_%d", i)] = map[string]interface{}{
			"type":    ruleType,
			"passed":  passed,
			"message": rule["message"],
		}
	}

	return mustMarshalJSON(map[string]interface{}{
		"all_passed": allPassed,
		"results":    results,
	}), allPassed
}
