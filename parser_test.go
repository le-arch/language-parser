package main

import (
	"reflect"
	"testing"
)

// TestParseAcceptLanguage covers both Part 1 and Part 2 requirements
func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name        string
		header      string
		supported   []string
		expected    []string
		description string
	}{
		// PART 1: BASIC EXACT MATCHING 
		{
			name:        "Part1 - Example 1: Multiple preferences with matches",
			header:      "en-US, fr-CA, fr-FR",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Client prefers en-US (1st), fr-CA (2nd), fr-FR (3rd). Server supports fr-FR, en-US. Should return [en-US, fr-FR] in preference order.",
		},
		{
			name:        "Part1 - Example 2: Partial match",
			header:      "fr-CA, fr-FR",
			supported:   []string{"en-US", "fr-FR"},
			expected:    []string{"fr-FR"},
			description: "Client prefers fr-CA (1st), fr-FR (2nd). Server supports en-US, fr-FR. Should return [fr-FR] (fr-CA not supported).",
		},
		{
			name:        "Part1 - Example 3: Exact match",
			header:      "en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US"},
			description: "Client prefers en-US. Server supports en-US, fr-CA. Should return [en-US].",
		},
		{
			name:        "Part1 - Empty header",
			header:      "",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Empty header should return empty slice.",
		},
		{
			name:        "Part1 - No supported languages",
			header:      "en-US, fr-CA",
			supported:   []string{},
			expected:    []string{},
			description: "No supported languages should return empty slice.",
		},
		{
			name:        "Part1 - No matches",
			header:      "de-DE, es-ES",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "When no client preferences match supported languages, return empty slice.",
		},
		{
			name:        "Part1 - All matches",
			header:      "en-US, fr-CA, fr-FR",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "When all client preferences are supported, return all in original order.",
		},
		{
			name:        "Part1 - Whitespace handling",
			header:      "  en-US  ,  fr-CA  ,  fr-FR  ",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Should properly trim whitespace around language tags.",
		},
		{
			name:        "Part1 - Extra commas and spaces",
			header:      "en-US, , fr-CA, , fr-FR",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Should skip empty entries between commas.",
		},
		{
			name:        "Part1 - Duplicate prevention",
			header:      "en-US, fr-CA, en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US", "fr-CA"},
			description: "Should not duplicate languages even if they appear multiple times in header.",
		},
		{
			name:        "Part1 - Case sensitivity",
			header:      "en-us, fr-ca",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Language tags are case-sensitive. 'en-us' ≠ 'en-US'.",
		},

		// PART 2: GENERIC TAG SUPPORT 
		{
			name:        "Part2 - Generic 'en' matches English variants",
			header:      "en",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"en-US"},
			description: "Generic 'en' should match all English variants (en-US).",
		},
		{
			name:        "Part2 - Generic 'fr' matches French variants",
			header:      "fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"fr-CA", "fr-FR"},
			description: "Generic 'fr' should match all French variants (fr-CA, fr-FR).",
		},
		{
			name:        "Part2 - Generic after specific - order preserved",
			header:      "fr-FR, fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"fr-FR", "fr-CA"},
			description: "Client prefers fr-FR first, then any French - should return fr-FR (exact) then fr-CA (generic match).",
		},
		{
			name:        "Part2 - Multiple generic tags",
			header:      "en, fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "Client wants English then French - should return all English then all French variants.",
		},
		{
			name:        "Part2 - Generic with multiple variants",
			header:      "es",
			supported:   []string{"es-ES", "es-MX", "es-AR", "en-US"},
			expected:    []string{"es-ES", "es-MX", "es-AR"},
			description: "Client wants Spanish, server has multiple variants - should return all Spanish variants.",
		},
		{
			name:        "Part2 - Exact match takes precedence over generic",
			header:      "fr-CA, fr",
			supported:   []string{"fr-FR", "fr-CA"},
			expected:    []string{"fr-CA", "fr-FR"},
			description: "Client explicitly wants fr-CA first, then any French - fr-CA should appear before fr-FR.",
		},
		{
			name:        "Part2 - Generic after no exact match",
			header:      "fr-BE, fr",
			supported:   []string{"fr-FR", "fr-CA"},
			expected:    []string{"fr-FR", "fr-CA"},
			description: "Client wants fr-BE (not supported), then any French - should return all French variants.",
		},
		{
			name:        "Part2 - Generic with no matches",
			header:      "zh",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Client wants Chinese, server has no Chinese variants - return empty.",
		},
		{
			name:        "Part2 - Mixed case with generic",
			header:      "EN",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Generic tags are case-sensitive - 'EN' doesn't match 'en-US'.",
		},

		// COMBINED SCENARIOS (Part 1 + Part 2)
		{
			name:        "Combined - Complex scenario with exact and generic",
			header:      "en-US, fr, de-DE, es",
			supported:   []string{"en-US", "en-GB", "fr-CA", "fr-FR", "de-DE", "es-ES", "es-MX"},
			expected:    []string{"en-US", "fr-CA", "fr-FR", "de-DE", "es-ES", "es-MX"},
			description: "Should handle mix of exact and generic tags correctly.",
		},
		{
			name:        "Combined - Duplicate prevention across exact and generic",
			header:      "en-US, en, en-GB",
			supported:   []string{"en-US", "en-GB", "en-AU"},
			expected:    []string{"en-US", "en-GB", "en-AU"},
			description: "Generic 'en' should match all English variants (en-GB, en-AU) after exact match en-US.",
		},
		{
			name:        "Combined - Order preservation with multiple variants",
			header:      "fr, en, es",
			supported:   []string{"es-ES", "en-US", "fr-FR", "fr-CA", "en-GB", "es-MX"},
			expected:    []string{"fr-FR", "fr-CA", "en-US", "en-GB", "es-ES", "es-MX"},
			description: "Should preserve order of generic tags and their variants.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAcceptLanguage(tt.header, tt.supported)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf(`
Test: %s
Description: %s
Header: %q
Supported: %v
Expected: %v
Got: %v
`, tt.name, tt.description, tt.header, tt.supported, tt.expected, result)
			}
		})
	}
}

// TestLanguageParser tests the struct-based implementation
func TestLanguageParser(t *testing.T) {
	parser := NewLanguageParser()
	
	tests := []struct {
		name     string
		header   string
		supported []string
		expected []string
	}{
		{
			name:     "Parser - Basic exact match",
			header:   "en-US, fr-CA, fr-FR",
			supported: []string{"fr-FR", "en-US"},
			expected: []string{"en-US", "fr-FR"},
		},
		{
			name:     "Parser - Generic match",
			header:   "fr",
			supported: []string{"en-US", "fr-CA", "fr-FR"},
			expected: []string{"fr-CA", "fr-FR"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.Parse(tt.header, tt.supported)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Parser.Parse(%q, %v) = %v, expected %v", 
					tt.header, tt.supported, result, tt.expected)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkParseAcceptLanguage(b *testing.B) {
	header := "en-US, fr-CA, fr-FR, de-DE, es-ES, it-IT, ja-JP, ko-KR, zh-CN, zh-TW"
	supported := []string{"en-US", "fr-FR", "ja-JP", "zh-CN", "es-ES", "de-DE"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseAcceptLanguage(header, supported)
	}
}

func BenchmarkParseAcceptLanguageWithGeneric(b *testing.B) {
	header := "en, fr, de, es, ja, zh"
	supported := []string{"en-US", "fr-FR", "de-DE", "es-ES", "ja-JP", "zh-CN"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseAcceptLanguage(header, supported)
	}
}