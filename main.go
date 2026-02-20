// Package main provides functions for parsing Accept-Language HTTP headers
// with support for exact matching, generic tags, and wildcards.
package main

import (
	"fmt"
	"strings"
)

// parseAcceptLanguage parses the Accept-Language header and returns
// supported languages in the client's preference order.
//
// Part 1: Basic exact matching
// Part 2: Generic tag support (e.g., "en" matches "en-US", "en-GB")
// Part 3: Wildcard support (* matches all remaining languages)
//
// Parameters:
//   - header: The Accept-Language header value as a string
//   - supported: A slice of language tags that the server supports
//
// Returns:
//   - A slice of language tags that satisfy client preference and server support
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
	
	// Also keep the original order for wildcard matching (Part 3)
	supportedOrder := make([]string, len(supported))
	copy(supportedOrder, supported)
	
	for _, lang := range supported {
		// Store for exact matching (Part 1)
		supportedExact[lang] = true
		
		// Extract generic part for generic matching (Part 2)
		parts := strings.SplitN(lang, "-", 2)
		generic := parts[0]
		supportedGeneric[generic] = append(supportedGeneric[generic], lang)
	}

	// Parse the header by splitting on commas and trimming whitespace
	clientPrefs := strings.Split(header, ",")
	result := make([]string, 0)
	seen := make(map[string]bool) // Track added languages to prevent duplicates

	// Process each client preference in order
	for _, pref := range clientPrefs {
		lang := strings.TrimSpace(pref)
		if lang == "" {
			continue // Skip empty entries
		}

		// PART 3: Check for wildcard
		if lang == "*" {
			// Wildcard matches all remaining supported languages not yet seen
			for _, supportedLang := range supportedOrder {
				if !seen[supportedLang] {
					result = append(result, supportedLang)
					seen[supportedLang] = true
				}
			}
			continue
		}

		// PART 1: Check for exact match first (higher priority)
		if supportedExact[lang] && !seen[lang] {
			result = append(result, lang)
			seen[lang] = true
			continue
		}

		// PART 2: Check if this is a generic tag (no hyphen)
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
	// Could add configuration options here
}

// NewLanguageParser creates a new instance of LanguageParser
func NewLanguageParser() *LanguageParser {
	return &LanguageParser{}
}

// Parse implements the same logic as parseAcceptLanguage
func (lp *LanguageParser) Parse(header string, supported []string) []string {
	return parseAcceptLanguage(header, supported)
}

// main demonstrates all examples from Parts 1, 2, and 3
func main() {
	fmt.Println("ACCEPT-LANGUAGE PARSER")
	fmt.Println("-----------------------")
	fmt.Println("Part 1: Exact Matching")
	fmt.Println("Part 2: Generic Tag Support") 
	fmt.Println("Part 3: Wildcard Support (*)")
	fmt.Println()

	// PART 1 EXAMPLES 
	fmt.Println("PART 1: Basic Exact Matching")
	fmt.Println("-----------------------------")

	header1 := "en-US, fr-CA, fr-FR"
	supported1 := []string{"fr-FR", "en-US"}
	result1 := parseAcceptLanguage(header1, supported1)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v\n\n", header1, supported1, result1)

	header2 := "fr-CA, fr-FR"
	supported2 := []string{"en-US", "fr-FR"}
	result2 := parseAcceptLanguage(header2, supported2)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v\n\n", header2, supported2, result2)

	// PART 2 EXAMPLES 
	fmt.Println("PART 2: Generic Tag Support")
	fmt.Println("---------------------------")

	header3 := "en"
	supported3 := []string{"en-US", "fr-CA", "fr-FR"}
	result3 := parseAcceptLanguage(header3, supported3)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v\n\n", header3, supported3, result3)

	header4 := "fr"
	supported4 := []string{"en-US", "fr-CA", "fr-FR"}
	result4 := parseAcceptLanguage(header4, supported4)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v\n\n", header4, supported4, result4)

	header5 := "fr-FR, fr"
	supported5 := []string{"en-US", "fr-CA", "fr-FR"}
	result5 := parseAcceptLanguage(header5, supported5)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v\n\n", header5, supported5, result5)

	// PART 3 EXAMPLES
	fmt.Println("PART 3: Wildcard Support (*)")
	fmt.Println("-----------------------------")

	// Example 1: Wildcard after specific match
	header6 := "en-US, *"
	supported6 := []string{"en-US", "fr-CA", "fr-FR"}
	result6 := parseAcceptLanguage(header6, supported6)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v", header6, supported6, result6)
	fmt.Println(" (Wildcard matches fr-CA, fr-FR)")
	fmt.Println()

	// Example 2: Wildcard with generic and exact
	header7 := "fr-FR, fr, *"
	supported7 := []string{"en-US", "fr-CA", "fr-FR"}
	result7 := parseAcceptLanguage(header7, supported7)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v", header7, supported7, result7)
	fmt.Println(" (Wildcard matches en-US)")
	fmt.Println()

	// Example 3: Wildcard alone
	header8 := "*"
	supported8 := []string{"en-US", "fr-CA", "fr-FR"}
	result8 := parseAcceptLanguage(header8, supported8)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v", header8, supported8, result8)
	fmt.Println(" (Wildcard matches all)")
	fmt.Println()

	// Example 4: Wildcard after unmatched preference
	header9 := "de-DE, *"
	supported9 := []string{"en-US", "fr-CA", "fr-FR"}
	result9 := parseAcceptLanguage(header9, supported9)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v", header9, supported9, result9)
	fmt.Println(" (de-DE not supported, wildcard matches all)")
	fmt.Println()

	// Example 5: Complex scenario
	header10 := "en, fr-FR, *"
	supported10 := []string{"en-US", "en-GB", "fr-CA", "fr-FR", "de-DE"}
	result10 := parseAcceptLanguage(header10, supported10)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v", header10, supported10, result10)
	fmt.Println(" (en matches en-US,en-GB; fr-FR exact; wildcard matches fr-CA,de-DE)")
	fmt.Println()

	// Example 6: Multiple wildcards
	header11 := "en-US, *, fr-FR, *"
	supported11 := []string{"en-US", "fr-CA", "fr-FR", "de-DE"}
	result11 := parseAcceptLanguage(header11, supported11)
	fmt.Printf("Header: %q\nSupported: %v\nResult: %v", header11, supported11, result11)
	fmt.Println(" (Second wildcard has nothing left to match)")
}
