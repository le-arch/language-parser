package main

import (
	"reflect"
	"testing"
)

// TestParseAcceptLanguage covers Parts 1, 2, and 3 requirements
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
			description: "Should return [en-US, fr-FR] in preference order.",
		},
		{
			name:        "Part1 - Example 2: Partial match",
			header:      "fr-CA, fr-FR",
			supported:   []string{"en-US", "fr-FR"},
			expected:    []string{"fr-FR"},
			description: "Should return [fr-FR] (fr-CA not supported).",
		},
		{
			name:        "Part1 - Example 3: Exact match",
			header:      "en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US"},
			description: "Should return [en-US].",
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
			name:        "Part1 - Whitespace handling",
			header:      "  en-US  ,  fr-CA  ,  fr-FR  ",
			supported:   []string{"fr-FR", "en-US"},
			expected:    []string{"en-US", "fr-FR"},
			description: "Should properly trim whitespace.",
		},
		{
			name:        "Part1 - Duplicate prevention",
			header:      "en-US, fr-CA, en-US",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{"en-US", "fr-CA"},
			description: "Should not duplicate languages.",
		},
		{
			name:        "Part1 - Case sensitivity",
			header:      "en-us, fr-ca",
			supported:   []string{"en-US", "fr-CA"},
			expected:    []string{},
			description: "Language tags are case-sensitive.",
		},

		// PART 2: GENERIC TAG SUPPORT 
		{
			name:        "Part2 - Generic 'en' matches English variants",
			header:      "en",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"en-US"},
			description: "Generic 'en' should match all English variants.",
		},
		{
			name:        "Part2 - Generic 'fr' matches French variants",
			header:      "fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"fr-CA", "fr-FR"},
			description: "Generic 'fr' should match all French variants.",
		},
		{
			name:        "Part2 - Generic after specific - order preserved",
			header:      "fr-FR, fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"fr-FR", "fr-CA"},
			description: "Should return fr-FR (exact) then fr-CA (generic match).",
		},
		{
			name:        "Part2 - Multiple generic tags",
			header:      "en, fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "Should return all English then all French variants.",
		},
		{
			name:        "Part2 - Generic with multiple variants",
			header:      "es",
			supported:   []string{"es-ES", "es-MX", "es-AR", "en-US"},
			expected:    []string{"es-ES", "es-MX", "es-AR"},
			description: "Should return all Spanish variants.",
		},

		// PART 3: WILDCARD SUPPORT 
		{
			name:        "Part3 - Wildcard after specific match",
			header:      "en-US, *",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "Wildcard matches all remaining languages after exact match",
		},
		{
			name:        "Part3 - Wildcard with generic and exact",
			header:      "fr-FR, fr, *",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"fr-FR", "fr-CA", "en-US"},
			description: "Wildcard matches remaining languages after exact and generic matches",
		},
		{
			name:        "Part3 - Wildcard alone",
			header:      "*",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "Wildcard alone matches all supported languages",
		},
		{
			name:        "Part3 - Wildcard after unmatched preference",
			header:      "de-DE, *",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "Wildcard matches all when previous preference doesn't match",
		},
		{
			name:        "Part3 - Complex scenario with generic and wildcard",
			header:      "en, fr-FR, *",
			supported:   []string{"en-US", "en-GB", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "en-GB", "fr-FR", "fr-CA", "de-DE"},
			description: "Generic en matches en-US,en-GB; fr-FR exact; wildcard adds remaining",
		},
		{
			name:        "Part3 - Multiple wildcards",
			header:      "en-US, *, fr-FR, *",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			description: "Multiple wildcards - second wildcard has nothing left to match",
		},
		{
			name:        "Part3 - Wildcard with all languages already matched",
			header:      "en-US, fr-CA, fr-FR, *",
			supported:   []string{"en-US", "fr-CA", "fr-FR"},
			expected:    []string{"en-US", "fr-CA", "fr-FR"},
			description: "Wildcard adds nothing when all languages already matched",
		},
		{
			name:        "Part3 - Wildcard preserves supported order",
			header:      "*, en-US",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			description: "Wildcard first matches all in supported order",
		},
		{
			name:        "Part3 - Wildcard with generic after",
			header:      "*, fr",
			supported:   []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "fr-CA", "fr-FR", "de-DE"},
			description: "Wildcard matches all, generic fr has nothing new to add",
		},

		// COMBINED SCENARIOS (All Parts) 
		{
			name:        "Combined - Complex scenario with exact, generic, and wildcard",
			header:      "en-US, fr, de-DE, *",
			supported:   []string{"en-US", "en-GB", "fr-CA", "fr-FR", "de-DE", "es-ES", "it-IT"},
			expected:    []string{"en-US", "fr-CA", "fr-FR", "de-DE", "en-GB", "es-ES", "it-IT"},
			description: "Should handle mix of exact, generic, and wildcard correctly",
		},
		{
			name:        "Combined - Wildcard after generic with multiple variants",
			header:      "en, *",
			supported:   []string{"en-US", "en-GB", "fr-CA", "fr-FR", "de-DE"},
			expected:    []string{"en-US", "en-GB", "fr-CA", "fr-FR", "de-DE"},
			description: "Generic en matches all English variants, wildcard adds the rest",
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
		{
			name:     "Parser - Wildcard match",
			header:   "en-US, *",
			supported: []string{"en-US", "fr-CA", "fr-FR"},
			expected: []string{"en-US", "fr-CA", "fr-FR"},
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

func BenchmarkParseAcceptLanguageWithWildcard(b *testing.B) {
	header := "en-US, fr, *"
	supported := []string{"en-US", "fr-CA", "fr-FR", "de-DE", "es-ES", "it-IT"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseAcceptLanguage(header, supported)
	}
}
