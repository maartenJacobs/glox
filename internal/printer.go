package internal

import (
	"fmt"
)

// stringify is the default printer for Lox values.
func stringify(loxValue interface{}) string {
	if loxValue == nil {
		return "nil"
	}

	switch v := loxValue.(type) {
	case float64:
		return fmt.Sprintf("%f", v)
	case string:
		return v
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	default:
		// Catch this case and add more type cases.
		return fmt.Sprintf("_%v", v)
	}
}
