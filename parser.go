// package main provides enhanced parser functionality for Accept-Language headers

package main 

import (
	"strings"
)

// LanguageParser provides methods for parsing Accept-Language headers
// using a struct allows for potential feature enhancements like:
// caching of parsed results
// configuring options (case-insensitive matching, etc)
//statistics collection
type LanguageParser struct {
	// 
}

// NewLanguageParser creates a new instance of LanguageParser
func NewLanguageParser() *LanguageParser {
	return  &LanguageParser{}
}

// parse accepts an Accept-Language header and returns supported languages
// in the client's preference order.
// this is the enhanced version of paseAcceptLanguage with better error handling and separation of concerns
// it implements the same logic but in a more maintainable structure
func (lp *LanguageParser) Parse(header string, supported []string) []string {
	// early returns for edge cases
	if header == "" {
		return  []string{}
	}

	if len(supported) == 0 {
		return []string{}
	}

	// build lookup set for supported languages
	supportedSet := make(map[string]struct{})
	for _, lang := range supported {
		supportedSet[lang] = struct{}{}
	}

	// parse header and filter
	preferences := strings.Split(header, ",")
	result := make([]string, 0)
	seen := make(map[string]bool)

	for _, pref := range preferences {
		lang := strings.TrimSpace(pref)
		if lang == "" {
			continue
		}

		// check if supported and not already in result
		if _, ok := supportedSet[lang]; ok && !seen[lang] {
			result = append(result, lang)
			seen[lang] = true
		}
	}

	return result
}