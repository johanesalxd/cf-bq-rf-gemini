package bqrfgemini_test

import (
	"encoding/json"
	"reflect"
	"testing"

	bqrfgemini "github.com/johanesalxd/cf-bq-rf-gemini"
)

// TestGenerateJSONResponse tests the GenerateJSONResponse function
func TestGenerateJSONResponse(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		// Test case 1: Ensure a valid struct is correctly serialized
		{
			name:     "Valid struct",
			input:    struct{ Name string }{"John"},
			expected: `{"Name":"John"}`,
		},
		// Test case 2: Ensure a valid map is correctly serialized
		{
			name:     "Valid map",
			input:    map[string]int{"Age": 30},
			expected: `{"Age":30}`,
		},
		// Test case 3: Ensure nil input is handled correctly
		{
			name:     "Nil input",
			input:    nil,
			expected: `null`,
		},
		// Test case 4: Ensure unmarshalable input returns an error message
		{
			name:     "Unmarshalable input",
			input:    make(chan int),
			expected: `{"error": "json: unsupported type: chan int"}`,
		},
	}

	// Run each test case
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

// TestParseModelConfig tests the ParseModelConfig function
func TestParseModelConfig(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		input    string
		expected bqrfgemini.ModelConfig
	}{
		// Test case 1: Ensure a valid full configuration is parsed correctly
		{
			name:  "Valid full config",
			input: `{"temperature":0.2,"maxOutputTokens":8000,"topP":0.8,"topK":40}`,
			expected: bqrfgemini.ModelConfig{
				Temperature:     0.2,
				MaxOutputTokens: 8000,
				TopP:            0.8,
				TopK:            40,
			},
		},
		// Test case 2: Ensure a partial configuration is parsed correctly with default values
		{
			name:  "Partial config",
			input: `{"temperature":0.5,"maxOutputTokens":5000}`,
			expected: bqrfgemini.ModelConfig{
				Temperature:     0.5,
				MaxOutputTokens: 5000,
				TopP:            0.95, // Default value
				TopK:            1,    // Default value
			},
		},
		// Test case 3: Ensure an empty configuration returns all default values
		{
			name:  "Empty config",
			input: `{}`,
			expected: bqrfgemini.ModelConfig{
				Temperature:     1,    // Default value
				MaxOutputTokens: 1000, // Default value
				TopP:            0.95, // Default value
				TopK:            1,    // Default value
			},
		},
		// Test case 4: Ensure invalid JSON returns all default values
		{
			name:  "Invalid JSON",
			input: `{"temperature":0.3,}`,
			expected: bqrfgemini.ModelConfig{
				Temperature:     1,    // Default value
				MaxOutputTokens: 1000, // Default value
				TopP:            0.95, // Default value
				TopK:            1,    // Default value
			},
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bqrfgemini.ParseModelConfig(tt.input)

			// Compare the result with the expected output
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseModelConfig() = %+v, want %+v", result, tt.expected)
			}
		})
	}
}
