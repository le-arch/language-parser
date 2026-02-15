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
	// 
	
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
	fmt.Println("------------------------")
	fmt.Println("Combining Part 1 (Exact Matching) and Part 2 (Generic Tag Support)")
	fmt.Println()

	// PART 1 EXAMPLES
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

	// PART 2 EXAMPLES 
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

	// COMBINED SCENARIOS 
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