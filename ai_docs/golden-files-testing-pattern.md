# Golden Files Testing Pattern

## Overview

The golden files testing pattern, popularized by Mitchell Hashimoto in his "Advanced Testing with Go" talk, is a technique for testing functions that produce complex output. Instead of hardcoding expected results in test code, the expected output is stored in separate files called "golden files."

## When to Use Golden Files

Golden files are particularly effective for:
- Testing functions that generate complex formatted output (JSON, HTML, Markdown)
- Verifying HTTP responses and headers
- Comparing large data structures
- Testing code generation or transformation tools
- Testing AI/ML model responses where output may vary slightly but should remain consistent
- Contract testing with external APIs

## Recommended File Organization

For better separation of concerns and maintainability, we recommend separating fixture management from functional tests:

```
mypackage/
├── mypackage.go
├── fixture_test.go      # Golden file updates (with -update flag)
├── mypackage_test.go    # Functional tests using fixtures
└── testdata/
    ├── TestCase1.golden
    ├── TestCase2.golden.json
    └── TestCase3.golden.xml
```

This separation provides:
- Clear distinction between fixture generation and test logic
- Easier code review (fixture changes vs test changes)
- Better organization for large test suites
- Isolated contract testing capabilities

## Basic Implementation

### 1. Simple Golden File Test

```go
package mypackage_test

import (
    "bytes"
    "flag"
    "os"
    "path/filepath"
    "testing"
)

var update = flag.Bool("update", false, "update .golden files")

func TestSomething(t *testing.T) {
    // Run the function under test
    actual := doSomething()
    
    // Define golden file path
    golden := filepath.Join("testdata", t.Name()+".golden")
    
    if *update {
        // Update mode: save actual output as golden file
        err := os.MkdirAll("testdata", 0755)
        if err != nil {
            t.Fatal(err)
        }
        err = os.WriteFile(golden, actual, 0644)
        if err != nil {
            t.Fatal(err)
        }
        t.Logf("Updated golden file: %s", golden)
    }
    
    // Read expected output from golden file
    expected, err := os.ReadFile(golden)
    if err != nil {
        if os.IsNotExist(err) && !*update {
            t.Fatalf("golden file %s does not exist. Run with -update flag to create it", golden)
        }
        t.Fatal(err)
    }
    
    // Compare actual vs expected
    if !bytes.Equal(actual, expected) {
        t.Errorf("output does not match golden file.\nExpected:\n%s\nGot:\n%s", expected, actual)
    }
}
```

### 2. Complex Response Golden Files (JSON)

```go
func TestComplexResponse(t *testing.T) {
    type Response struct {
        Status  string   `json:"status"`
        Data    []string `json:"data"`
        Version string   `json:"version"`
    }
    
    golden := filepath.Join("testdata", t.Name()+".golden.json")
    
    if *update {
        // Generate actual response
        actual := Response{
            Status:  "success",
            Data:    []string{"item1", "item2"},
            Version: "1.0.0",
        }
        
        // Marshal to JSON with pretty printing
        data, err := json.MarshalIndent(actual, "", "  ")
        if err != nil {
            t.Fatal(err)
        }
        
        // Save to golden file
        err = os.MkdirAll("testdata", 0755)
        if err != nil {
            t.Fatal(err)
        }
        err = os.WriteFile(golden, data, 0644)
        if err != nil {
            t.Fatal(err)
        }
        t.Logf("Updated golden file: %s", golden)
    }
    
    // Read and unmarshal golden file
    data, err := os.ReadFile(golden)
    if err != nil {
        if os.IsNotExist(err) && !*update {
            t.Fatalf("golden file %s does not exist. Run with -update flag", golden)
        }
        t.Fatal(err)
    }
    
    var expected Response
    err = json.Unmarshal(data, &expected)
    if err != nil {
        t.Fatal(err)
    }
    
    // Get actual response from your function
    actual := getComplexResponse()
    
    // Compare
    if !reflect.DeepEqual(actual, expected) {
        t.Errorf("response mismatch.\nExpected: %+v\nGot: %+v", expected, actual)
    }
}
```

## Running Golden File Tests

### Normal Test Run
```bash
go test ./...
```
This compares actual output against existing golden files.

### Update Golden Files
When using the `-update` flag, you need to specify only the packages that have golden tests:

```bash
# Update all golden files in packages that have them
go test ./scanner/googlebooks ./scanner/vertexai ./scanner/vision -update

# Or update one package at a time
go test ./scanner/vertexai -update
```

**Note**: You cannot run `go test ./... -update` because not all packages define the `-update` flag. Only packages with `golden_test.go` files have this flag.

### Update Specific Test
```bash
go test ./scanner/vertexai -run TestSomething -update
```

## Best Practices

### 1. Manual Review is Critical
After updating golden files, **always manually review the changes**. The update flag blindly saves whatever the current output is, so you must verify it's correct before committing.

```bash
# After updating golden files
git diff testdata/
# Review each change carefully
```

### 2. Use testdata Directory
Go's build system ignores directories named `testdata`, making it the standard location for test fixtures:

```
mypackage/
├── mypackage.go
├── mypackage_test.go
└── testdata/
    ├── TestCase1.golden
    ├── TestCase2.golden.json
    └── TestCase3.golden.xml
```

### 3. Version Control Golden Files
Always commit golden files to version control. They're part of your test suite and should be reviewed in pull requests when they change.

### 4. Avoid Timestamps and Non-Deterministic Data
If your output includes timestamps or random values, normalize them before comparison:

```go
func normalizeOutput(data []byte) []byte {
    // Replace timestamps with fixed value
    re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`)
    return re.ReplaceAll(data, []byte("2023-01-01T00:00:00"))
}
```

### 5. Platform-Specific Line Endings
Handle line ending differences between platforms:

```go
import "bytes"

func normalizeLineEndings(data []byte) []byte {
    // Convert all line endings to \n
    data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
    data = bytes.ReplaceAll(data, []byte("\r"), []byte("\n"))
    return data
}
```

## Separated Fixture and Test Approach

### fixture_test.go - Golden File Management

This file is responsible for updating golden files with real API responses when run with the `-update` flag:

```go
package mypackage_test

import (
    "encoding/json"
    "flag"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

// TestUpdateGoldenFiles updates fixtures when run with -update flag
func TestUpdateGoldenFiles(t *testing.T) {
    if !*update {
        t.Skip("Run with -update flag to update golden files")
    }

    // This test makes real API calls and saves responses
    testCases := []struct {
        name  string
        input string
    }{
        {
            name:  "simple_query",
            input: "test input",
        },
        // more test cases...
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Make real API call
            response := makeRealAPICall(tc.input)
            
            // Save response to golden file
            golden := filepath.Join("testdata", tc.name+".golden.json")
            data, err := json.MarshalIndent(response, "", "  ")
            require.NoError(t, err)
            
            err = os.WriteFile(golden, data, 0644)
            require.NoError(t, err)
            t.Logf("Updated golden file: %s", golden)
        })
    }
}
```

### mypackage_test.go - Functional Tests

This file contains the actual tests that use the golden files as fixtures:

```go
package mypackage_test

import (
    "encoding/json"
    "os"
    "path/filepath"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestFunctionality(t *testing.T) {
    testCases := []struct {
        name          string
        input         string
        goldenFile    string
        validateFunc  func(t *testing.T, got, want interface{})
    }{
        {
            name:       "simple_query",
            input:      "test input",
            goldenFile: "simple_query.golden.json",
        },
        // more test cases...
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Load golden file
            golden := filepath.Join("testdata", tc.goldenFile)
            data, err := os.ReadFile(golden)
            require.NoError(t, err, "golden file should exist")
            
            var expected Response
            err = json.Unmarshal(data, &expected)
            require.NoError(t, err)
            
            // Run function under test with mocked client
            mockClient := createMockFromGolden(expected)
            result := functionUnderTest(mockClient, tc.input)
            
            // Validate results
            assert.Equal(t, expected, result)
        })
    }
}
```

### Benefits of Separation

1. **Clear Responsibilities**: 
   - `fixture_test.go` handles external API interactions
   - `mypackage_test.go` focuses on business logic testing

2. **Better CI/CD Integration**:
   - Regular test runs don't require API access
   - Golden file updates can be run manually or on schedule
   - Reduces flaky tests due to network issues

3. **Improved Code Review**:
   - Fixture changes are isolated and easy to review
   - Test logic changes don't mix with data updates
   - Contract changes are immediately visible

4. **Contract Documentation**:
   - Golden files serve as living API documentation
   - No separate contract docs to maintain
   - Real examples from actual API responses

## Real-World Example: Contract Testing with Golden Files

Golden files naturally serve as contract documentation by capturing real API responses. Each golden file becomes a living document of the API contract:

```json
// testdata/book_search_gatsby.golden.json
{
  "kind": "books#volumes",
  "totalItems": 1,
  "items": [
    {
      "volumeInfo": {
        "title": "The Great Gatsby",
        "authors": ["F. Scott Fitzgerald"],
        "publishedDate": "1925",
        "industryIdentifiers": [
          {
            "type": "ISBN_13",
            "identifier": "9780743273565"
          }
        ]
      }
    }
  ]
}
```

This approach eliminates the need for separate contract documentation - the fixtures themselves document:
- API response structure
- Field types and formats
- Real-world examples
- Edge cases and variations

## Common Pitfalls

1. **Forgetting to Review**: Always review golden file changes before committing
2. **Binary Files**: Use appropriate comparison for binary files (images, etc.)
3. **Large Files**: Consider using checksums for very large golden files
4. **Missing Files**: Provide clear error messages when golden files are missing

## Alternatives and Libraries

- **sebdah/goldie**: Full-featured golden file testing library
- **xorcare/golden**: Minimal golden file testing package
- **gotestfmt**: Can format test output for better golden file diffs

## Summary

Golden file testing is a powerful pattern for testing complex outputs. It makes tests more maintainable by separating test logic from expected data, and the `-update` flag workflow makes it easy to update expectations as code evolves. Just remember: always manually verify golden file updates before committing them.

## Book Scanner Implementation

The Book Scanner project demonstrates an evolved approach to golden file testing that addresses common pitfalls:

### Key Improvements

1. **Simplified Test Structure**: Instead of complex HTTP-level mocking, tests either:
   - Make real API calls when `-update` flag is set
   - Use saved golden files with domain-level mocks otherwise

2. **Single Flag Definition**: Each package defines the `-update` flag only once in `golden_helper_test.go` to avoid conflicts

3. **Local Utilities**: Following Ben Johnson's Standard Package Layout, golden test helpers are kept local to each package rather than in a shared internal package

### Example from googlebooks Package

```go
// golden_helper_test.go
var update = flag.Bool("update", false, "update golden files")

func compareWithGolden(t *testing.T, filename string, actual interface{}) {
    t.Helper()
    actualJSON, err := json.MarshalIndent(actual, "", "  ")
    require.NoError(t, err)
    
    if *update {
        updateGolden(t, filename, actualJSON)
        return
    }
    
    expectedJSON := loadGolden(t, filename)
    assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}

// golden_test.go
func TestEnrichMetadata_Golden(t *testing.T) {
    if *update {
        // Real API call to update golden files
        client, err := googlebooks.New(apiKey)
        require.NoError(t, err)
        defer client.Close()
        
        result, err := client.EnrichMetadata(ctx, "The Great Gatsby")
        require.NoError(t, err)
        compareWithGolden(t, "great_gatsby.json", result)
        return
    }
    
    // Normal test run - just validate golden file exists and is valid
    goldenData := loadGolden(t, "great_gatsby.json")
    var book bookscanner.Book
    err := json.Unmarshal(goldenData, &book)
    require.NoError(t, err)
    
    // Basic validation
    assert.NotEmpty(t, book.Title)
    assert.NotEmpty(t, book.Author)
}
```

This approach ensures tests are fast, reliable, and maintain proper separation between testing the API contract (golden files) and testing business logic (unit tests with mocks).