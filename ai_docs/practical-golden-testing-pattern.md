# Practical Golden Testing Pattern for External APIs

This document describes the golden testing pattern as successfully implemented in production Go projects for testing external API integrations.

## Overview

Golden file testing is a pattern where expected test outputs are stored in files (golden files) and compared against actual outputs during test runs. For external APIs, this pattern enables:

- Fast, reliable tests without external dependencies
- Documentation of actual API responses
- Easy updates when APIs change
- Confidence that code handles real API responses correctly

## Directory Structure

```
package/
├── client.go                # Your API client implementation
├── client_test.go          # Client tests using golden files
├── golden_test.go          # Golden file flag definition
├── golden_update_test.go   # Test that updates golden files with -update flag
└── testdata/               # Go convention for test data
    ├── case1.json
    ├── case2.json
    └── error_case.json
```

## Implementation Pattern

### 1. Define the Update Flag (golden_test.go)

```go
package mypackage

import "flag"

var update = flag.Bool("update", false, "update golden files")

// This file exists to define the update flag once per package
// The flag is shared across all test files in the package
```

### 2. Create Update Test (golden_update_test.go)

This test makes real API calls when run with `-update` flag:

```go
package mypackage

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)

func TestUpdateGoldenFiles(t *testing.T) {
    if !*update {
        t.Skip("Run with -update flag to update golden files")
    }

    // Check for API credentials
    apiKey := os.Getenv("API_KEY")
    if apiKey == "" {
        t.Skip("API_KEY environment variable not set")
    }

    client := NewClient(apiKey)

    testCases := []struct {
        name       string
        input      string
        goldenFile string
    }{
        {
            name:       "successful_search",
            input:      "test query",
            goldenFile: "successful_search.json",
        },
        {
            name:       "no_results", 
            input:      "nonexistent",
            goldenFile: "no_results.json",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Make real API call
            result, err := client.Search(tc.input)
            require.NoError(t, err)

            // Save response to golden file
            goldenPath := filepath.Join("testdata", tc.goldenFile)
            
            // Ensure testdata directory exists
            err = os.MkdirAll("testdata", 0755)
            require.NoError(t, err)
            
            // Marshal with pretty printing for readability
            data, err := json.MarshalIndent(result, "", "  ")
            require.NoError(t, err)
            
            err = os.WriteFile(goldenPath, data, 0644)
            require.NoError(t, err)
            
            t.Logf("Updated golden file: %s", goldenPath)
        })
    }
}
```

### 3. Create Test Helper (client_test.go)

Create a test helper that returns responses from golden files:

```go
package mypackage

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)

// testClient returns responses from golden files instead of making API calls
type testClient struct {
    t *testing.T
}

func newTestClient(t *testing.T) *testClient {
    t.Helper()
    return &testClient{t: t}
}

func (tc *testClient) Search(query string) (*SearchResult, error) {
    // Map queries to golden files
    var goldenFile string
    switch query {
    case "test query":
        goldenFile = "successful_search.json"
    case "nonexistent":
        goldenFile = "no_results.json"
    default:
        tc.t.Fatalf("No golden file mapped for query: %s", query)
    }

    // Load golden file
    goldenPath := filepath.Join("testdata", goldenFile)
    data, err := os.ReadFile(goldenPath)
    require.NoError(tc.t, err, "golden file should exist: %s", goldenPath)

    var result SearchResult
    err = json.Unmarshal(data, &result)
    require.NoError(tc.t, err)

    return &result, nil
}
```

### 4. Write Tests Using Golden Files

```go
func TestSearch(t *testing.T) {
    t.Parallel()
    
    if *update {
        t.Skip("Skip normal tests when updating golden files")
    }

    tests := []struct {
        name           string
        query          string
        expectedTitle  string
        expectedCount  int
    }{
        {
            name:          "successful_search",
            query:         "test query",
            expectedTitle: "Test Result",
            expectedCount: 1,
        },
        {
            name:          "no_results",
            query:         "nonexistent",
            expectedCount: 0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            
            client := newTestClient(t)
            result, err := client.Search(tt.query)
            require.NoError(t, err)
            
            assert.Equal(t, tt.expectedCount, len(result.Items))
            if tt.expectedCount > 0 {
                assert.Equal(t, tt.expectedTitle, result.Items[0].Title)
            }
        })
    }
}
```

## Usage

### Running Tests Normally
```bash
go test ./mypackage
```

### Updating Golden Files
```bash
# Set API credentials
export API_KEY=your-api-key

# Update golden files
go test ./mypackage -update

# Review changes
git diff testdata/

# Commit if correct
git add testdata/
git commit -m "Update golden files"
```

## Best Practices

### 1. One Golden File Per Test Case
- Makes it easy to understand what each file represents
- Simplifies updates when only specific cases change
- Use descriptive file names that match test names

### 2. Store Raw API Responses
- Store the actual API response, not your processed domain objects
- This tests your parsing/transformation logic
- Documents the real API contract

### 3. Use JSON Format
- Human-readable for easy inspection
- Stable formatting with `json.MarshalIndent`
- Use `assert.JSONEq` for order-independent comparison

### 4. Environment-Based Updates
- Require API credentials only for `-update` mode
- Tests run without credentials using golden files
- Clear skip messages guide developers

### 5. Handle Errors Gracefully
```go
if *update && apiKey == "" {
    t.Skip("API_KEY required for updating golden files")
}
```

### 6. Parallel Testing
- Always use `t.Parallel()` for better test performance
- Helps detect race conditions
- Golden file reads are safe for concurrent access

### 7. Test Both Success and Error Cases
- Create golden files for error responses
- Test error handling paths
- Document various API error scenarios

## Common Patterns

### Testing Internal Types
When you need access to internal types, keep tests in the same package:

```go
package mypackage // not mypackage_test

func TestInternalBehavior(t *testing.T) {
    // Can access unexported types and functions
}
```

### Partial Response Testing
For APIs that return timestamps or changing data:

```go
// Focus on stable fields
assert.Equal(t, expected.Title, actual.Title)
assert.Equal(t, expected.Author, actual.Author)
// Ignore unstable fields like timestamps
```

### Integration Tests
Use build tags for tests that always need real API:

```go
//go:build integration

func TestRealAPI(t *testing.T) {
    // Always makes real API calls
}
```

Run with: `go test -tags=integration`

## Advantages

1. **Fast Tests**: No network calls during normal test runs
2. **Reliable**: No flaky tests due to network issues
3. **Documented**: Golden files show exact API responses
4. **Maintainable**: Easy to update when APIs change
5. **Debuggable**: Can inspect golden files to understand failures
6. **CI/CD Friendly**: Tests run without credentials in CI

## Summary

This pattern provides a practical, maintainable approach to testing external API integrations. It balances the need for realistic tests with the requirements of fast, reliable test suites. The key is making it easy to update golden files when needed while keeping normal test runs fast and dependency-free.