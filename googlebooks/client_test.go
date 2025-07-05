package googlebooks_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
				// Verify thumbnail URL uses HTTPS
				if result.ThumbnailURL != "" {
					assert.True(t, strings.HasPrefix(result.ThumbnailURL, "https://"), "Thumbnail URL should use HTTPS")
				}
			},
		},
		{
			name:               "title_author_search",
			goldenFile:         "title_author_gatsby.json",
			expectedResults:    10,
			expectedFirstTitle: "The Great Gatsby",
			validateFields: func(t *testing.T, result bookid.BookResult) {
				t.Helper()
				assert.NotEmpty(t, result.Authors)
				assert.Equal(t, bookid.SearchTypeGeneralQuery, result.SearchType)
				assert.InDelta(t, 0.70, result.Confidence, 0.01)
				// Verify thumbnail URL uses HTTPS
				if result.ThumbnailURL != "" {
					assert.True(t, strings.HasPrefix(result.ThumbnailURL, "https://"), "Thumbnail URL should use HTTPS")
				}
			},
		},
		{
			name:               "title_only_search",
			goldenFile:         "title_only_pride.json",
			expectedResults:    10,
			expectedFirstTitle: "Pride and Prejudice",
			validateFields: func(t *testing.T, result bookid.BookResult) {
				t.Helper()
				assert.NotEmpty(t, result.Title)
				assert.Equal(t, bookid.SearchTypeGeneralQuery, result.SearchType)
				assert.InDelta(t, 0.70, result.Confidence, 0.06)
				// Verify thumbnail URL uses HTTPS
				if result.ThumbnailURL != "" {
					assert.True(t, strings.HasPrefix(result.ThumbnailURL, "https://"), "Thumbnail URL should use HTTPS")
				}
			},
		},
		{
			name:            "no_results",
			goldenFile:      "no_results.json",
			expectedResults: 0,
		},
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

// TestClient_Search_Errors tests error handling
func TestClient_Search_Errors(t *testing.T) {
	t.Parallel()

	// Create client without API key
	client, err := googlebooks.NewClient("")
	require.NoError(t, err)

	// Test empty query
	_, err = client.Search(context.Background(), "")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "query cannot be empty")
}
