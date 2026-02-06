# Accept-Language Parser

A Go solution to parse HTTP Accept-Language headers and match them against server-supported languages.

## Problem Statement

Given an HTTP `Accept-Language` header (a comma-separated list of language preferences) and a list of languages the server supports, return the languages that work for the request **in the client's original preference order**.

### Example

```go
parseAcceptLanguage("en-US, fr-CA, fr-FR", ["fr-FR", "en-US"])
// Returns: ["en-US", "fr-FR"]
```

The client prefers:
1. English US (most preferred)
2. French Canada
3. French France (least preferred)

The server supports French France and English US, so we return `["en-US", "fr-FR"]` in the original preference order.

## Features

- ✅ Parses comma-separated Accept-Language headers
- ✅ Maintains original preference order
- ✅ Handles whitespace properly
- ✅ Case-sensitive matching (per HTTP spec)
- ✅ Prevents duplicate results
- ✅ Comprehensive test coverage
- ✅ Two implementations (function and struct-based)

## Installation

```bash
# Clone and initialize
git clone <repository>
cd accept-language-parser
go mod init accept-language-parser
```

## Usage

### Basic Function

```go
import "strings"

result := parseAcceptLanguage("en-US, fr-CA, fr-FR", []string{"fr-FR", "en-US"})
// result = ["en-US", "fr-FR"]
```

### Enhanced Parser (Struct-based)

```go
parser := NewLanguageParser()
result := parser.Parse("en-US, fr-CA", []string{"en-US", "fr-CA"})
// result = ["en-US", "fr-CA"]
```

### Run Examples

```bash
go run main.go parser.go
```

## Running Tests

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Run benchmarks
go test -bench=.

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Cases Covered

1. **Basic Examples** from problem statement
2. **Edge Cases**: Empty header, no supported languages, no matches
3. **Whitespace Handling**: Extra spaces, trailing commas
4. **Duplicate Prevention**: Multiple occurrences in header
5. **Case Sensitivity**: Language tags must match exactly
6. **Order Preservation**: Maintains client's preference order

## Algorithm

1. **Parse** the header by commas and trim whitespace
2. **Create a set** of supported languages for O(1) lookups
3. **Filter** client preferences to only include supported languages
4. **Remove duplicates** while maintaining first occurrence position
5. **Return** filtered list in original order

**Time Complexity**: O(n + m) where n = header languages, m = supported languages  
**Space Complexity**: O(m + k) where k = result size

## Design Decisions

1. **Case Sensitivity**: HTTP language tags are case-sensitive (RFC 7231)
2. **Exact Matching**: Requires full tag match (no partial matching like "fr" for "fr-FR")
3. **Order Preservation**: Client preference order is sacred
4. **No Duplicates**: Even if language appears multiple times in header
5. **Empty Handling**: Empty header or no supported languages returns empty slice

## File Structure

```
accept-language-parser/
├── main.go          # Main implementation and examples
├── parser.go        # Enhanced struct-based parser
├── parser_test.go   # Comprehensive test suite
└── go.mod          # Go module file
```