# Golden Files Testing Pattern

## Overview

The golden files testing pattern is a technique for testing functions that produce complex output, particularly useful for external API integrations. Instead of hardcoding expected results in test code, the expected output is stored in separate files called "golden files."

This pattern has evolved from simple file comparison to a sophisticated approach for API contract testing while maintaining clean package boundaries and testability.

## Core Principles

1. **Package Boundary Testing**: All tests must use the `_test` package suffix to ensure you're testing through the public API
2. **No Test Code in Production**: Never expose test-only functions or flags in your production code
3. **Separate Concerns**: Golden file generation (with real API calls) is separate from validation tests
4. **Domain Types in Golden Files**: Store your domain types, not raw API responses, to maintain clean boundaries

## Implementation Pattern

### Directory Structure

```
googlebooks/
├── client.go           # Production code
├── parser.go           # Production code
├── client_test.go      # Contract validation tests (package googlebooks_test)
├── parser_test.go      # Unit tests (package googlebooks_test)
├── golden_update_test.go # Golden file generation (package googlebooks_test)
└── testdata/
    ├── isbn_9780743273565.json
    ├── title_author_gatsby.json
    └── no_results.json
```

### 1. Golden File Update Test

This test is responsible for making real API calls and storing the results as golden files:

```go
// golden_update_test.go
package googlebooks_test

import (
    "context"
    "encoding/json"
    "flag"
    "os"
    "path/filepath"
    "testing"

    "github.com/fwojciec/bookid/googlebooks"
    "github.com/stretchr/testify/require"
)

// Internal flag for golden file updates
var update = flag.Bool("update", false, "update golden files")

// TestUpdateGoldenFiles updates golden files when run with -update flag
// This test should only be run manually when API responses need to be updated
func TestUpdateGoldenFiles(t *testing.T) {
    t.Parallel()

    if !*update {
        t.Skip("Run with -update flag to update golden files")
    }

    // Get API key from environment
    apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
    if apiKey == "" {
        t.Skip("GOOGLE_BOOKS_API_KEY environment variable not set")
    }

    t.Log("Updating golden files with real API responses...")

    // Test cases that will generate golden files
    testCases := []struct {
        query      string
        goldenFile string
    }{
        {
            query:      "9780743273565",
            goldenFile: "isbn_9780743273565.json",
        },
        {
            query:      "The Great Gatsby by F. Scott Fitzgerald",
            goldenFile: "title_author_gatsby.json",
        },
        {
            query:      "nonexistentbook12345",
            goldenFile: "no_results.json",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.goldenFile, func(t *testing.T) {
            t.Parallel()
            // Create real client using public API
            client, err := googlebooks.NewClient(apiKey)
            require.NoError(t, err)

            // Make real API call
            ctx := context.Background()
            results, err := client.Search(ctx, tc.query)
            require.NoError(t, err)

            // Save to golden file
            goldenPath := filepath.Join("testdata", tc.goldenFile)

            // Ensure testdata directory exists
            err = os.MkdirAll("testdata", 0755)
            require.NoError(t, err)

            // Marshal with pretty printing
            data, err := json.MarshalIndent(results, "", "  ")
            require.NoError(t, err)

            err = os.WriteFile(goldenPath, data, 0644)
            require.NoError(t, err)

            t.Logf("Updated golden file: %s", goldenPath)
            t.Logf("Results count: %d", len(results))
        })
    }
}
```

### 2. Contract Validation Tests

These tests validate that golden files contain the expected structure and serve as contract tests:

```go
// client_test.go
package googlebooks_test

import (
    "context"
    "encoding/json"
    "os"
    "path/filepath"
    "testing"

    "github.com/fwojciec/bookid"
    "github.com/fwojciec/bookid/googlebooks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// TestGoldenFiles validates that the golden files contain the expected data structure
// This serves as a contract test - when golden files are updated with -update flag,
// this test ensures the API responses still contain the fields we depend on
func TestGoldenFiles(t *testing.T) {
    t.Parallel()

    // Check if golden files exist - if not, skip validation
    // This allows golden file generation to happen independently
    testDataPath := filepath.Join("testdata", "isbn_9780743273565.json")
    if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
        t.Skip("Golden files not found - run tests with -update flag in the googlebooks package first")
    }

    tests := []struct {
        name               string
        goldenFile         string
        expectedResults    int
        expectedFirstTitle string
        validateFields     func(t *testing.T, result bookid.BookResult)
    }{
        {
            name:               "isbn_search",
            goldenFile:         "isbn_9780743273565.json",
            expectedResults:    1,
            expectedFirstTitle: "The Great Gatsby",
            validateFields: func(t *testing.T, result bookid.BookResult) {
                t.Helper()
                assert.NotEmpty(t, result.ISBN10)
                assert.NotEmpty(t, result.ISBN13)
                assert.NotEmpty(t, result.Publisher)
                assert.NotEmpty(t, result.GoogleBooksVolumeID)
                assert.Equal(t, bookid.SearchTypeISBN, result.SearchType)
                assert.InDelta(t, 0.95, result.Confidence, 0.01)
            },
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()

            // Load golden file
            goldenPath := filepath.Join("testdata", tt.goldenFile)
            data, err := os.ReadFile(goldenPath)
            require.NoError(t, err, "golden file should exist: %s", goldenPath)

            var results []bookid.BookResult
            err = json.Unmarshal(data, &results)
            require.NoError(t, err, "golden file should contain valid BookResult array")

            // Validate structure
            assert.Len(t, results, tt.expectedResults)

            if tt.expectedResults > 0 {
                assert.Equal(t, tt.expectedFirstTitle, results[0].Title)

                // Run field-specific validations
                if tt.validateFields != nil {
                    tt.validateFields(t, results[0])
                }
            }
        })
    }
}
```

## Key Design Decisions

### 1. Store Domain Types, Not Raw API Responses

Instead of storing raw API responses, transform them to your domain types before saving. This:
- Allows all tests to remain in the `_test` package
- Tests the actual transformation logic during golden file generation
- Makes golden files serve as documentation of your domain model
- Enables contract testing without exposing internals

### 2. Make Tests Independent

The validation tests check if golden files exist before running. This prevents test failures when golden files haven't been generated yet and makes the test suite more robust.

### 3. Use Public API Only

All tests, including golden file generation, use only the public API of your package. This ensures you're testing what consumers of your package will actually use.

## Running Tests

### Normal Test Run
```bash
# Run all tests (validates against existing golden files)
go test ./...
```

### Generate/Update Golden Files
```bash
# Set up environment (if needed)
export GOOGLE_BOOKS_API_KEY="your-api-key"

# Update golden files for a specific package
go test ./googlebooks -update

# Update golden files with verbose output
go test ./googlebooks -update -v
```

### Running in CI/CD

```yaml
# CI pipeline example
test:
  steps:
    # Regular tests run without API access
    - name: Run Tests
      run: go test -race ./...
    
    # Golden file updates run separately (e.g., nightly)
    - name: Update Golden Files
      if: github.event_name == 'schedule'
      env:
        GOOGLE_BOOKS_API_KEY: ${{ secrets.GOOGLE_BOOKS_API_KEY }}
      run: |
        go test ./googlebooks -update
        # Commit changes if any
```

## Best Practices

### 1. Always Use t.Parallel()

Enable parallel test execution and race detection:

```go
func TestSomething(t *testing.T) {
    t.Parallel()
    // test code...
}
```

### 2. Review Golden File Changes

After updating golden files, always review the changes:

```bash
# Update golden files
go test ./googlebooks -update

# Review changes
git diff testdata/

# Ensure changes are expected before committing
```

### 3. Use t.Helper() in Test Helpers

Mark test helper functions to improve error reporting:

```go
func validateFields(t *testing.T, result BookResult) {
    t.Helper()
    assert.NotEmpty(t, result.Title)
    // ... more validations
}
```

### 4. Keep Golden Files Small and Focused

Each golden file should represent a specific test case. Don't create huge golden files that test multiple scenarios.

### 5. Handle Missing Golden Files Gracefully

Provide clear error messages when golden files are missing:

```go
if _, err := os.Stat(testDataPath); os.IsNotExist(err) {
    t.Skip("Golden files not found - run tests with -update flag first")
}
```

## Common Pitfalls and Solutions

### Pitfall 1: Test-Only Code in Production

**Wrong:**
```go
// client.go
func IsUpdatingGoldenFiles() bool {
    return *update
}
```

**Right:**
Keep all test code in test files and use the public API.

### Pitfall 2: Internal Package Tests for Golden Files

**Wrong:**
```go
// golden_test.go
package googlebooks // internal package test

func TestUpdateGolden(t *testing.T) {
    // Can access private fields, but breaks testpackage linting rule
}
```

**Right:**
Use the `_test` package and store domain types in golden files.

### Pitfall 3: Hardcoding Expected Values

**Wrong:**
```go
// Hardcoding expected values in tests
assert.Equal(t, "The Great Gatsby", results[0].Title)
assert.Equal(t, []string{"F. Scott Fitzgerald"}, results[0].Authors)
```

**Right:**
Store expectations in golden files and validate structure, not specific values.

## Advanced Patterns

### Mock Client Using Golden Files

Create a mock client that returns data from golden files:

```go
type MockClient struct {
    responses map[string][]bookid.BookResult
}

func NewMockClientFromGoldenFiles() (*MockClient, error) {
    mock := &MockClient{
        responses: make(map[string][]bookid.BookResult),
    }
    
    // Load all golden files
    files, err := filepath.Glob("testdata/*.json")
    if err != nil {
        return nil, err
    }
    
    for _, file := range files {
        data, err := os.ReadFile(file)
        if err != nil {
            return nil, err
        }
        
        var results []bookid.BookResult
        if err := json.Unmarshal(data, &results); err != nil {
            return nil, err
        }
        
        // Use filename as key
        key := filepath.Base(file)
        mock.responses[key] = results
    }
    
    return mock, nil
}
```

### Contract Evolution Testing

Track API contract changes over time:

```go
func TestContractCompatibility(t *testing.T) {
    // Load current golden file
    current := loadGoldenFile(t, "current_response.json")
    
    // Load previous version
    previous := loadGoldenFile(t, "previous_response.json")
    
    // Verify backward compatibility
    assertBackwardCompatible(t, previous, current)
}
```

## Summary

The golden files pattern, when properly implemented:
1. Maintains clean package boundaries with all tests in `_test` packages
2. Keeps test concerns out of production code
3. Provides living documentation of API contracts
4. Enables fast, reliable tests that don't require external dependencies
5. Makes contract changes visible in code reviews

By storing domain types instead of raw API responses, we achieve the best of both worlds: clean architecture and practical testing.