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
			// Create real client
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
