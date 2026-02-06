// package main provides functions for parsing Accept-Language HTTP headers.
package main

import(
	"fmt"
	"strings"
)

// parseAcceptLanguage parses the Accept-Language header and returns
// a filtered list of languages that are both requested by the client
// and supported by the server, in the original preference order

// parameters:
// header: the Accept-Language header value as a string (e.g., "en-US, fr-CA, fr-FR")
// supported: a slice of language tags that the server supports

// returns:
// a slice of language tags that satisfy both client preference and server support
// returns empty slice if no matches found or invalid input

// key characteristics:
// language tags are case-sensitive (RFC 7231)
// whitespace around tags is trimmed 
// original preference order is preserved
// duplicates in header result in single occurrence in output
// exact matching is required (no partial matching like "fr" matching "fr-FR")

func parseAcceptLanguage(header string, supported []string) []string {
	// edge case: empty header or no supported languages
	if header == "" || len(supported) == 0{
		return []string{}
	}


// create a map of supported languages for 0(1) lookups

// using a map[string]struct{} for memory efficiency as we only need existence checks
supportedSet := make(map[string]struct{})
for _, lang := range supported {
	supportedSet[lang] = struct{}{}
}

// split the header by commas to get individual language preferences
// according to RFC 7231, Accept-Language is comma-separated

clientPrefs := strings.Split(header, ",")
result := make([]string, 0)

// use a map to track languages we've already added to avoid duplicates
// this maintains the first occurrence position while preventing duplicates
seen := make(map[string]bool)

// process each client preference in order (preserving preference order)
for _, pref := range clientPrefs {
	// trim whitespace as per HTTP header specifications
	// headers can have spaces after commas: "en-US, fr-CA"
	
	lang := strings.TrimSpace(pref)

	// skip empty entries (could occur from trailing commas or extra spaces)
	if lang == "" {
		continue
	}

	// check if this language is both supported and not already added
	if _, isSupported := supportedSet[lang]; isSupported && !seen[lang] {
		result = append(result, lang)
		seen[lang] = true
	}
}
return result
}

// example usage demonstrating the function
func main(){
	fmt.Println("Accept-Language Parser Examples:")
	fmt.Println("_________________________________")

	// example 1: basic case from problem statement
	header1 := "en-US, fr-CA, fr-FR"
	supported1 := []string{"fr-FR", "en-US"}
	result1 := parseAcceptLanguage(header1, supported1)
	fmt.Printf("Header: %q\n", header1)
	fmt.Printf("Supported: %v\n", supported1)
	fmt.Printf("Result: %v\n\n", result1)

	// example 2: partial match
	header2 := "fr-CA, fr-FR"
	supported2 := []string{"en-US", "fr-FR"}
	result2 := parseAcceptLanguage(header2, supported2)
	fmt.Printf("Header: %q\n", header2)
	fmt.Printf("Supported: %v\n", supported2)
	fmt.Printf("Result: %v\n\n", result2)

	// example 3: exact match
	header3 := "en-US"
	supported3 := []string{"en-US", "fr-CA"}
	result3 := parseAcceptLanguage(header3, supported3)
	fmt.Printf("Header: %q\n", header3)
	fmt.Printf("Supported: %v\n", supported3)
	fmt.Printf("Result: %v\n\n", result3)

	// example 4: no matches
	header4 := "de-DE, es-ES"
	supported4 := []string{"en-US", "fr-CA"}
	result4 := parseAcceptLanguage(header4, supported4)
	fmt.Printf("Header: %q\n", header4)
	fmt.Printf("Supported: %v\n", supported4)
	fmt.Printf("Result: %v (no matches)\n\n", result4)

	// example 5: with whitespace
	header5 := "  en-US  ,  fr-CA  ,  fr-FR  "
	supported5 := []string{"fr-FR", "en-US"}
	result5 := parseAcceptLanguage(header5, supported5)
	fmt.Printf("Header: %q\n", header5)
	fmt.Printf("Supported: %v\n", supported5)
	fmt.Printf("Result: %v (with whitespace handling)\n", result5)
}