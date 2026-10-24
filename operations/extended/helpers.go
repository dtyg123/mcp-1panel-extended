package extended

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// strArg returns the string value for key, or def when absent/not a string.
func strArg(args map[string]interface{}, key, def string) string {
	if args != nil {
		if v, ok := args[key].(string); ok {
			return v
		}
	}
	return def
}

// numArg returns the numeric value for key (JSON numbers decode to float64), or def.
func numArg(args map[string]interface{}, key string, def float64) float64 {
	if args != nil {
		if v, ok := args[key].(float64); ok {
			return v
		}
	}
	return def
}

// boolArg returns the boolean value for key, or def.
func boolArg(args map[string]interface{}, key string, def bool) bool {
	if args != nil {
		if v, ok := args[key].(bool); ok {
			return v
		}
	}
	return def
}

// requiredStrArg returns the string value for key, or an error when missing/empty.
func requiredStrArg(args map[string]interface{}, key string) (string, error) {
	v := strArg(args, key, "")
	if v == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return v, nil
}

// prettyJSON re-indents raw JSON while preserving original number lexemes.
// Returns the input unchanged when it is not valid JSON.
func prettyJSON(raw []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, "", "  "); err == nil {
		return buf.String()
	}
	return string(raw)
}
