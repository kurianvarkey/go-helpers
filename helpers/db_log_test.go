package helpers

import (
	"testing"
	"time"
)

func TestLogBoundSQL(t *testing.T) {
	// Define a reusable test time in UTC
	testTime, _ := time.Parse(time.RFC3339, "2025-11-28T14:30:00Z")

	// --- Test Cases ---
	name := "Alice O'Brien"
	tests := []struct {
		name     string
		query    string
		args     []any
		expected string
	}{
		{
			name:     "Standard Select with String and Int",
			query:    "SELECT * FROM users WHERE id = $1 AND name = $2",
			args:     []any{123, name},                                                 // Note the single quote in the name
			expected: "SELECT * FROM users WHERE id = 123 AND name = 'Alice O''Brien'", // Should be escaped
		},
		{
			name:     "Time, Null, and Boolean",
			query:    "INSERT INTO logs (created_at, user_id, status) VALUES ($1, $2, $3)",
			args:     []any{testTime, nil, true},
			expected: "INSERT INTO logs (created_at, user_id, status) VALUES ('2025-11-28T14:30:00.000000000Z', NULL, true)",
		},
		{
			name:     "Placeholder Out of Order/Skip",
			query:    "SELECT total FROM orders WHERE id = $3 AND date = $1",
			args:     []any{"2024-01-01", 1, 999}, // Order is args[0], args[1], args[2]
			expected: "SELECT total FROM orders WHERE id = 999 AND date = '2024-01-01'",
		},
		{
			name:     "Missing Placeholder (Not Substituted)",
			query:    "SELECT * FROM products WHERE price > $1 AND stock < $5",
			args:     []any{10.50}, // Only $1 is provided
			expected: "SELECT * FROM products WHERE price > 10.5 AND stock < $5",
		},
		{
			name:     "No Arguments",
			query:    "SELECT count(*) FROM items",
			args:     nil,
			expected: "SELECT count(*) FROM items",
		},
		{
			name:     "With Pointer",
			query:    "SELECT * FROM users WHERE name = $1",
			args:     []any{&name},                                        // Note the single quote in the name
			expected: "SELECT * FROM users WHERE name = 'Alice O''Brien'", // Should be escaped
		},
	}

	// --- Execution ---
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetBoundSQL(tt.query, tt.args...)
			if actual != tt.expected {
				t.Errorf("\nInput Query: %s\nInput Args: %v\nExpected: %s\nActual:   %s",
					tt.query, tt.args, tt.expected, actual)
			}
		})
	}
}
