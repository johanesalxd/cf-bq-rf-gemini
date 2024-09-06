package bqrfgemini_test

import (
	"encoding/json"
	"testing"

	bqrfgemini "github.com/johanesalxd/cf-bq-rf-gemini"
)

func TestGenerateJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "Valid struct",
			input:    struct{ Name string }{"John"},
			expected: `{"Name":"John"}`,
		},
		{
			name:     "Valid map",
			input:    map[string]int{"Age": 30},
			expected: `{"Age":30}`,
		},
		{
			name:     "Nil input",
			input:    nil,
			expected: `null`,
		},
		{
			name:     "Unmarshalable input",
			input:    make(chan int),
			expected: `{"error": "json: unsupported type: chan int"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bqrfgemini.GenerateJSONResponse(tt.input)

			// Compare JSON strings
			if string(result) != tt.expected {
				t.Errorf("generateJSONResponse() = %v, want %v", string(result), tt.expected)
			}

			// Verify it's valid JSON
			var js json.RawMessage
			if err := json.Unmarshal(result, &js); err != nil {
				t.Errorf("generateJSONResponse() produced invalid JSON: %v", err)
			}
		})
	}
}
