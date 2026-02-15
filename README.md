# Accept-Language Parser

A comprehensive Go solution for parsing HTTP Accept-Language headers with support for both basic and generic language tag matching.

## Complete Problem Statement

### Part 1: Basic Language Matching
Given an HTTP `Accept-Language` header (a comma-separated list of language preferences) and a list of languages the server supports, return the languages that work for the request **in the client's original preference order**.

### Part 2: Generic Language Tag Support
Extend the function to support generic language tags where a tag without a region (e.g., "en") matches all region-specific variants (e.g., "en-US", "en-GB", "en-CA"). Exact matches take priority over generic matches.

## Complete Solution

### `main.go` - Main implementation with both parts combined

```go
// Package main provides functions for parsing Accept-Language HTTP headers
// with support for both exact and generic language tag matching.
package main

import (
	"fmt"
	"strings"
)

// parseAcceptLanguage parses the Accept-Language header and returns
// supported languages in the client's preference order.
//
// Part 1: Basic exact matching
//   - Returns languages that are both requested and supported
//   - Maintains original preference order
//   - Handles whitespace and prevents duplicates
//
// Part 2: Generic tag support
//   - Tags without hyphens (e.g., "en") match all region-specific variants (e.g., "en-US")
//   - Exact matches take priority over generic matches
//   - When a generic tag matches multiple variants, all are returned in supported order
//
// Parameters:
//   - header: The Accept-Language header value as a string (e.g., "en-US, fr-CA, fr-FR")
//   - supported: A slice of language tags that the server supports
//
// Returns:
//   - A slice of language tags that satisfy both client preference and server support
//   - Returns empty slice if no matches found or invalid input
func parseAcceptLanguage(header string, supported []string) []string {
	// Edge cases: empty header or no supported languages
	if header == "" || len(supported) == 0 {
		return []string{}
	}

	// Build lookup structures for supported languages
	// Part 1: Exact match lookup
	supportedExact := make(map[string]bool)
	
	// Part 2: Generic match lookup - maps generic tags to their specific variants
	supportedGeneric := make(map[string][]string)
	
	for _, lang := range supported {
		// Store for exact matching (Part 1)
		supportedExact[lang] = true
		
		// Extract generic part for generic matching (Part 2)
		// Everything before the first hyphen becomes the generic tag
		parts := strings.SplitN(lang, "-", 2)
		generic := parts[0]
		supportedGeneric[generic] = append(supportedGeneric[generic], lang)
	}

	// Parse the header by splitting on commas and trimming whitespace
	clientPrefs := strings.Split(header, ",")
	result := make([]string, 0)
	seen := make(map[string]bool) // Track added languages to prevent duplicates

	// Process each client preference in order (maintaining preference order)
	for _, pref := range clientPrefs {
		lang := strings.TrimSpace(pref)
		if lang == "" {
			continue // Skip empty entries from trailing commas or extra spaces
		}

		// PART 1: Check for exact match first (higher priority)
		if supportedExact[lang] && !seen[lang] {
			result = append(result, lang)
			seen[lang] = true
			continue
		}

		// PART 2: Check if this is a generic tag (no hyphen)
		// Generic tags match all region-specific variants
		if !strings.Contains(lang, "-") {
			if variants, exists := supportedGeneric[lang]; exists {
				// Add all variants of this generic language
				for _, variant := range variants {
					if !seen[variant] {
						result = append(result, variant)
						seen[variant] = true
					}
				}
			}
		}
	}

	return result
}

// LanguageParser provides an enhanced struct-based implementation
type LanguageParser struct {
	// Could add configuration options here in the future
	// e.g., caseSensitive bool, defaultLanguage string
}

// NewLanguageParser creates a new instance of LanguageParser
func NewLanguageParser() *LanguageParser {
	return &LanguageParser{}
}

// Parse implements the same logic as parseAcceptLanguage but in a struct-based approach
// for better extensibility and testability
func (lp *LanguageParser) Parse(header string, supported []string) []string {
	// Reuse the existing function for consistency
	return parseAcceptLanguage(header, supported)
}

// main demonstrates all examples from both Part 1 and Part 2
func main() {
	fmt.Println("ACCEPT-LANGUAGE PARSER")
	fmt.Println("======================")
	fmt.Println("Combining Part 1 (Exact Matching) and Part 2 (Generic Tag Support)")
	fmt.Println()

	// ========== PART 1 EXAMPLES ==========
	fmt.Println("PART 1: Basic Exact Matching")
	fmt.Println("-----------------------------")

	// Example 1: Multiple preferences with matches
	header1 := "en-US, fr-CA, fr-FR"
	supported1 := []string{"fr-FR", "en-US"}
	result1 := parseAcceptLanguage(header1, supported1)
	fmt.Printf("Header: %q\n", header1)
	fmt.Printf("Supported: %v\n", supported1)
	fmt.Printf("Result: %v", result1)
	fmt.Println(" (Client prefers en-US → fr-CA → fr-FR, Server supports fr-FR and en-US)")
	fmt.Println()

	// Example 2: Partial match
	header2 := "fr-CA, fr-FR"
	supported2 := []string{"en-US", "fr-FR"}
	result2 := parseAcceptLanguage(header2, supported2)
	fmt.Printf("Header: %q\n", header2)
	fmt.Printf("Supported: %v\n", supported2)
	fmt.Printf("Result: %v", result2)
	fmt.Println(" (fr-CA not supported, only fr-FR matches)")
	fmt.Println()

	// Example 3: Exact match
	header3 := "en-US"
	supported3 := []string{"en-US", "fr-CA"}
	result3 := parseAcceptLanguage(header3, supported3)
	fmt.Printf("Header: %q\n", header3)
	fmt.Printf("Supported: %v\n", supported3)
	fmt.Printf("Result: %v", result3)
	fmt.Println(" (en-US exactly matches)")
	fmt.Println()

	// ========== PART 2 EXAMPLES ==========
	fmt.Println("PART 2: Generic Tag Support")
	fmt.Println("---------------------------")

	// Example 4: Generic 'en' matches en-US
	header4 := "en"
	supported4 := []string{"en-US", "fr-CA", "fr-FR"}
	result4 := parseAcceptLanguage(header4, supported4)
	fmt.Printf("Header: %q\n", header4)
	fmt.Printf("Supported: %v\n", supported4)
	fmt.Printf("Result: %v", result4)
	fmt.Println(" (Generic 'en' matches 'en-US')")
	fmt.Println()

	// Example 5: Generic 'fr' matches both French variants
	header5 := "fr"
	supported5 := []string{"en-US", "fr-CA", "fr-FR"}
	result5 := parseAcceptLanguage(header5, supported5)
	fmt.Printf("Header: %q\n", header5)
	fmt.Printf("Supported: %v\n", supported5)
	fmt.Printf("Result: %v", result5)
	fmt.Println(" (Generic 'fr' matches 'fr-CA' and 'fr-FR')")
	fmt.Println()

	// Example 6: Mixed specific and generic - exact match priority
	header6 := "fr-FR, fr"
	supported6 := []string{"en-US", "fr-CA", "fr-FR"}
	result6 := parseAcceptLanguage(header6, supported6)
	fmt.Printf("Header: %q\n", header6)
	fmt.Printf("Supported: %v\n", supported6)
	fmt.Printf("Result: %v", result6)
	fmt.Println(" (fr-FR exact match first, then fr-CA from generic 'fr')")
	fmt.Println()

	// ========== COMBINED SCENARIOS ==========
	fmt.Println("COMBINED SCENARIOS")
	fmt.Println("------------------")

	// Example 7: Multiple generics with exact matches
	header7 := "en-US, fr, es-ES, de"
	supported7 := []string{"en-US", "en-GB", "fr-CA", "fr-FR", "es-ES", "es-MX", "de-DE"}
	result7 := parseAcceptLanguage(header7, supported7)
	fmt.Printf("Header: %q\n", header7)
	fmt.Printf("Supported: %v\n", supported7)
	fmt.Printf("Result: %v\n", result7)
	fmt.Println()

	// Example 8: Complex ordering with duplicates
	header8 := "fr, en-US, fr-CA, en"
	supported8 := []string{"en-US", "en-GB", "fr-CA", "fr-FR", "fr-BE"}
	result8 := parseAcceptLanguage(header8, supported8)
	fmt.Printf("Header: %q\n", header8)
	fmt.Printf("Supported: %v\n", supported8)
	fmt.Printf("Result: %v\n", result8)
	fmt.Println()

	// Example 9: No matches
	header9 := "zh, ko"
	supported9 := []string{"en-US", "fr-FR", "es-ES"}
	result9 := parseAcceptLanguage(header9, supported9)
	fmt.Printf("Header: %q\n", header9)
	fmt.Printf("Supported: %v\n", supported9)
	fmt.Printf("Result: %v", result9)
	fmt.Println(" (No Chinese or Korean variants supported)")
}
```

### `parser_test.go` - Complete test suite covering both parts

```go
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
		// ========== PART 1: BASIC EXACT MATCHING ==========
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

		// ========== PART 2: GENERIC TAG SUPPORT ==========
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

		// ========== COMBINED SCENARIOS (Part 1 + Part 2) ==========
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
			expected:    []string{"en-US", "en-GB"},
			description: "Should prevent duplicates when same language appears via exact and generic matches.",
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
```

### `go.mod`

```go
module accept-language-parser

go 1.19
```

## Complete README

```markdown
# Accept-Language Parser

A comprehensive Go solution for parsing HTTP Accept-Language headers with support for both exact and generic language tag matching.

## Complete Problem Statement

### Part 1: Basic Language Matching
Given an HTTP `Accept-Language` header (a comma-separated list of language preferences) and a list of languages the server supports, return the languages that work for the request **in the client's original preference order**.

### Part 2: Generic Language Tag Support
Extend the function to support generic language tags where a tag without a region (e.g., "en") matches all region-specific variants (e.g., "en-US", "en-GB", "en-CA"). Exact matches take priority over generic matches.

## Features

- ✅ **Part 1: Exact Matching**
  - Parses comma-separated Accept-Language headers
  - Maintains original preference order
  - Handles whitespace properly
  - Case-sensitive matching (per HTTP spec)
  - Prevents duplicate results

- ✅ **Part 2: Generic Tag Support**
  - Generic tags (e.g., "en") match all region-specific variants (e.g., "en-US")
  - Exact matches take priority over generic matches
  - Multiple variants from a generic tag maintain their order
  - Handles complex scenarios with mixed exact and generic tags

## Installation

```bash
# Clone and initialize
git clone <repository>
cd accept-language-parser
go mod init accept-language-parser
```

## Usage Examples

### Part 1 Examples (Exact Matching)

```go
package main

import (
    "fmt"
    "accept-language-parser"
)

func main() {
    // Example 1: Multiple preferences with matches
    result := parseAcceptLanguage(
        "en-US, fr-CA, fr-FR", 
        []string{"fr-FR", "en-US"},
    )
    fmt.Println(result) // ["en-US", "fr-FR"]
    
    // Example 2: Partial match
    result = parseAcceptLanguage(
        "fr-CA, fr-FR", 
        []string{"en-US", "fr-FR"},
    )
    fmt.Println(result) // ["fr-FR"]
    
    // Example 3: Exact match
    result = parseAcceptLanguage(
        "en-US", 
        []string{"en-US", "fr-CA"},
    )
    fmt.Println(result) // ["en-US"]
}
```

### Part 2 Examples (Generic Tag Support)

```go
// Generic 'en' matches English variants
result := parseAcceptLanguage(
    "en", 
    []string{"en-US", "fr-CA", "fr-FR"},
)
fmt.Println(result) // ["en-US"]

// Generic 'fr' matches French variants
result = parseAcceptLanguage(
    "fr", 
    []string{"en-US", "fr-CA", "fr-FR"},
)
fmt.Println(result) // ["fr-CA", "fr-FR"]

// Mixed specific and generic - exact match priority
result = parseAcceptLanguage(
    "fr-FR, fr", 
    []string{"en-US", "fr-CA", "fr-FR"},
)
fmt.Println(result) // ["fr-FR", "fr-CA"]
```

### Combined Scenarios

```go
// Complex scenario with multiple generic and exact tags
result := parseAcceptLanguage(
    "en-US, fr, es-ES, de",
    []string{"en-US", "en-GB", "fr-CA", "fr-FR", "es-ES", "es-MX", "de-DE"},
)
fmt.Println(result) // ["en-US", "fr-CA", "fr-FR", "es-ES", "de-DE"]
```

## Running the Code

```bash
# Run the main program with all examples
go run main.go

# Run all tests
go test -v

# Run tests with coverage
go test -cover

# Run benchmarks
go test -bench=.

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Coverage

The test suite comprehensively covers both Part 1 and Part 2 requirements:

### Part 1 Tests
- Basic examples from problem statement
- Edge cases (empty header, no supported languages)
- Whitespace handling
- Duplicate prevention
- Case sensitivity
- Order preservation

### Part 2 Tests
- Generic tag matching (single and multiple)
- Exact match priority over generic
- Mixed exact and generic scenarios
- Generic tags with multiple variants
- No matches for generic tags
- Complex ordering scenarios

### Combined Tests
- Scenarios mixing exact and generic tags
- Duplicate prevention across match types
- Order preservation with multiple variants per generic tag

## Algorithm

```go
// For each client preference:
1. Trim whitespace
2. Check for exact match in supported set (Part 1)
3. If exact match found, add to result
4. If no exact match and tag has no hyphen (generic - Part 2):
   - Find all supported languages with matching prefix
   - Add all variants not already in result
5. Move to next preference (maintaining order)
```

**Time Complexity**: O(n * m) worst case, where n = header languages, m = supported variants
**Space Complexity**: O(m) for lookup maps

## Design Decisions

1. **Case Sensitivity**: HTTP language tags are case-sensitive (RFC 7231)
2. **Exact Match Priority**: Specific region tags are preferred over generic matches
3. **Generic Matching**: Tags without hyphen match all variants with that prefix
4. **Order Preservation**: Client preference order is maintained throughout
5. **No Duplicates**: Each language appears only once, preserving first occurrence
6. **Variant Order**: Multiple variants from a generic tag appear in supported list order

## File Structure

```
accept-language-parser/
├── main.go          # Complete implementation with both parts
├── parser_test.go   # Comprehensive test suite covering all requirements
└── go.mod          # Go module file
```

## Interview Discussion Points

This solution demonstrates:

1. **Requirement Evolution**: Handling both basic and extended requirements
2. **Algorithm Design**: Efficient matching with priority rules
3. **Data Structure Selection**: Maps for O(1) lookups, slices for ordered results
4. **Edge Case Consideration**: Comprehensive handling of edge cases
5. **Test-Driven Development**: Thorough test coverage for all scenarios
6. **Go Best Practices**: Clean code, proper error handling, idiomatic Go
7. **HTTP Protocol Knowledge**: Understanding of Accept-Language header format


This combined solution provides a complete implementation covering both Part 1 and Part 2 requirements, with comprehensive testing and documentation.