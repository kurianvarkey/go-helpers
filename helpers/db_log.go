package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const customNanoLayout = "2006-01-02T15:04:05.000000000Z"

// GetBoundSQL takes the raw query and its arguments,
// and returns a single string with the arguments substituted into the query.
func GetBoundSQL(query string, args ...any) string {
	// Pattern to find placeholders $1, $2, etc.
	re := regexp.MustCompile(`\$(\d+)`)

	finalQuery := re.ReplaceAllStringFunc(query, func(match string) string {
		indexStr := match[1:]

		index, err := strconv.Atoi(indexStr)
		if err != nil || index == 0 || index > len(args) {
			return match
		}

		arg := args[index-1]

		return formatArgForSQL(arg)
	})

	return finalQuery
}

// formatArgForSQL safely formats a Go variable into a string suitable for SQL logging.
func formatArgForSQL(arg any) string {
	if arg == nil {
		return "NULL"
	}

	v := reflect.ValueOf(arg)
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "NULL" // Pointer is nil
		}
		v = v.Elem() // Dereference the pointer
	}
	// v now holds the concrete (non-pointer) value
	unwrappedArg := v.Interface()

	switch val := unwrappedArg.(type) {
	case string:
		// Escape single quotes and wrap in quotes for SQL string literal
		escaped := strings.ReplaceAll(val, "'", "''")
		return fmt.Sprintf("'%s'", escaped)
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, float32, float64:
		// Numbers are safe
		return fmt.Sprintf("%v", val)
	case bool:
		// Booleans
		return fmt.Sprintf("%t", val)
	case time.Time:
		// Check if the time is UTC or has a fixed offset.
		// If not, convert it to UTC first to guarantee the 'Z' suffix is accurate.
		t := val
		if val.Location().String() != "UTC" {
			t = val.In(time.UTC)
		}

		// Format the time using the custom layout
		return fmt.Sprintf("'%s'", t.Format(customNanoLayout))
	default:
		// Fallback for other types
		return fmt.Sprintf("'%v'", v)
	}
}
