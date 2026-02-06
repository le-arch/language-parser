// package main_test contains comprehensive tests for the Accept-Language parser.
// these tests verify correctness, edge cases, and performance characteristics.

package main

import (
	"reflect"
	"testing"
	"fmt"
)

// TestParseAcceptLanguage tests the basic parseAcceptLanguage function with all the examples from the problem statement and additional edge cases.
func TestParseAcceptLanguage(t *testing.T) {
	// define test cases using table-driven testing pattern
	// this makes it easy to add new test cases and maintain existing ones
	tests := []struct {
		name        string           // test case identifier
		header      string           // input Accept-Language header
		supported   []string         // server-supported languages
		expected    []string         // expected output
		description string           // human-readable description
	}{
		// EXAMPLES FROM PROBLEM STATEMENT 
		{
			name:        "Example 1: Multiple preferences with matches",
			header:      "en-US, fr-CA, fr-FR",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Client prefers: en-US (1st), fr-CA (2nd), fr-FR (3rd). Server supports: fr-FR, en-US. Should return [en-US, fr-FR] in preference order.",
		},
		{
			name:        "Example 2: Partial match",
			header:      "fr-CA, fr-FR",
			supported:   []string{"en-US", "fr-FR"},
			expected:    []string{"fr-FR"},
			description: "Client prefers: fr-CA (1st), fr-FR (2nd). Server supports: en-US, fr-FR. Should return [fr-FR] (fr-CA not supported).",
		},
		{
			name:        "Example 3: Exact match",
			header:      "en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US"},
			description: "Client prefers: en-US. Server supports: en-US, fr-CA. Should return [en-US].",
		},

		// EDGE CASES
		{
			name:        "Empty header",
			header:      "",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Empty header should return empty slice.",
		},
		{
			name:        "No supported languages",
			header:      "en-US, fr-CA",
			supported:   []string{},
			expected:    []string{},
			description: "No supported languages should return empty slice.",
		},
		{
			name:        "No matches",
			header:      "de-DE, es-ES",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "When no client preferences match supported languages, return empty slice.",
		},
		{
			name:        "All matches",
			header:      "en-US, fr-CA, fr-FR",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "When all client preferences are supported, return all in original order.",
		},

		// WHITESPACE HANDLING
		{
			name:        "Whitespace around languages",
			header:      "  en-US  ,  fr-CA  ,  fr-FR  ",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Should properly trim whitespace around language tags.",
		},
		{
			name:        "Extra commas and spaces",
			header:      "en-US, , fr-CA, , fr-FR",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Should skip empty entries between commas.",
		},
		{
			name:        "Single language with spaces",
			header:      "  en-US  ",
			supported:   []string{"en-US"},
			expected:    []string{"en-US"},
			description: "Should handle single language with surrounding spaces.",
		},

		// DUPLICATE HANDLING 
		{
			name:        "Duplicate in header",
			header:      "en-US, fr-CA, en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US", "fr-CA"},
			description: "Should not duplicate languages even if they appear multiple times in header.",
		},
		{
			name:        "Multiple duplicates",
			header:      "en-US, fr-CA, en-US, fr-CA, en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US", "fr-CA"},
			description: "Should handle multiple duplicates gracefully.",
		},

		// CASE SENSITIVITY
		{
			name:        "Case sensitivity - lowercase",
			header:      "en-us, fr-ca",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Language tags are case-sensitive. 'en-us' ≠ 'en-US'.",
		},
		{
			name:        "Case sensitivity - mixed case",
			header:      "En-Us, Fr-Ca",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Case must match exactly per RFC specification.",
		},

		// COMPLEX SCENARIOS 
		{
			name:        "Reversed preference order",
			header:      "fr-FR, en-US, fr-CA",
			supported:   []string{"en-US", "fr-FR", "fr-CA"},
			expected:    []string{"fr-FR", "en-US", "fr-CA"},
			description: "Should maintain exact preference order from header.",
		},
		{
			name:        "Subset of supported",
			header:      "en-US",
			supported:   []string{"en-US", "fr-CA", "de-DE", "es-ES", "ja-JP"},
			expected:    []string{"en-US"},
			description: "Client requests one language that is among many supported.",
		},
		{
			name:        "Long header",
			header:      "zh-CN, en-US, ja-JP, ko-KR, fr-FR, de-DE, es-ES, it-IT",
			supported:   []string{"en-US", "fr-FR", "ja-JP"},
			expected:    []string{"en-US", "ja-JP", "fr-FR"},
			description: "Should handle long headers and maintain preference order.",
		},
	}

	// run all test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// execute the function
			result := parseAcceptLanguage(tt.header, tt.supported)
			
			// compare result with expected using reflect.DeepEqual
			// this compares both slice contents and order
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf(`
Test Failed: %s
Description: %s
Input Header: %q
Input Supported: %v
Expected Result: %v
Actual Result: %v
				`, tt.name, tt.description, tt.header, tt.supported, tt.expected, result)
			}
		})
	}
}

// TestLanguageParser tests the struct-based LanguageParser
// this ensures both implementations behave identically
func TestLanguageParser(t *testing.T) {
	parser := NewLanguageParser()
	
	tests := []struct {
		name      string
		header    string
		supported []string
		expected  []string
	}{
		{
			name:      "LanguageParser basic",
			header:    "en-US, fr-CA, fr-FR",
			supported: []string{"fr-FR", "en-US"},
			expected:  []string{"en-US", "fr-FR"},
		},
		{
			name:      "LanguageParser no match",
			header:    "de-DE, es-ES",
			supported: []string{"en-US", "fr-CA"},
			expected:  []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.Parse(tt.header, tt.supported)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("LanguageParser.Parse(%q, %v) = %v, expected %v", 
					tt.header, tt.supported, result, tt.expected)
			}
		})
	}
}

// TestEquivalence ensures both function implementations produce same results
func TestEquivalence(t *testing.T) {
	testCases := []struct {
		header    string
		supported []string
	}{
		{"en-US, fr-CA, fr-FR", []string{"fr-FR", "en-US"}},
		{"fr-CA, fr-FR", []string{"en-US", "fr-FR"}},
		{"en-US", []string{"en-US", "fr-CA"}},
		{"", []string{"en-US", "fr-CA"}},
		{"en-US, fr-CA", []string{}},
		{"en-US, en-US, fr-CA", []string{"en-US", "fr-CA"}},
	}
	
	parser := NewLanguageParser()
	
	for i, tc := range testCases {
		result1 := parseAcceptLanguage(tc.header, tc.supported)
		result2 := parser.Parse(tc.header, tc.supported)
		
		if !reflect.DeepEqual(result1, result2) {
			t.Errorf("Test case %d: Implementations differ\nFunction: %v\nParser: %v", 
				i, result1, result2)
		}
	}
}

// benchmark tests for performance comparison
func BenchmarkParseAcceptLanguage(b *testing.B) {
	header := "en-US, fr-CA, fr-FR, de-DE, es-ES, it-IT, ja-JP, ko-KR, zh-CN, zh-TW"
	supported := []string{"en-US", "fr-FR", "ja-JP", "zh-CN"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseAcceptLanguage(header, supported)
	}
}

func BenchmarkLanguageParser(b *testing.B) {
	header := "en-US, fr-CA, fr-FR, de-DE, es-ES, it-IT, ja-JP, ko-KR, zh-CN, zh-TW"
	supported := []string{"en-US", "fr-FR", "ja-JP", "zh-CN"}
	parser := NewLanguageParser()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(header, supported)
	}
}

// example test that appears in Go documentation
func ExampleparseAcceptLanguage() {
	result := parseAcceptLanguage("en-US, fr-CA, fr-FR", []string{"fr-FR", "en-US"})
	for _, lang := range result {
		fmt.Println(lang)
	}
	// Output:
	// en-US
	// fr-FR
}