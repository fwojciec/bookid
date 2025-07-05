package googlebooks_test

import (
	"testing"

	"github.com/fwojciec/bookid"
	"github.com/fwojciec/bookid/googlebooks"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		input         string
		expectedQuery string
		expectedType  bookid.SearchType
		expectedISBN  string
	}{
		{
			name:          "isbn10_only",
			input:         "0743273567",
			expectedQuery: "isbn:0743273567",
			expectedType:  bookid.SearchTypeISBN,
			expectedISBN:  "0743273567",
		},
		{
			name:          "isbn13_only",
			input:         "9780743273565",
			expectedQuery: "isbn:9780743273565",
			expectedType:  bookid.SearchTypeISBN,
			expectedISBN:  "9780743273565",
		},
		{
			name:          "isbn10_with_dashes",
			input:         "0-7432-7356-7",
			expectedQuery: "isbn:0743273567",
			expectedType:  bookid.SearchTypeISBN,
			expectedISBN:  "0743273567",
		},
		{
			name:          "isbn13_with_dashes",
			input:         "978-0-7432-7356-5",
			expectedQuery: "isbn:9780743273565",
			expectedType:  bookid.SearchTypeISBN,
			expectedISBN:  "9780743273565",
		},
		{
			name:          "isbn_in_mixed_text",
			input:         "The Great Gatsby ISBN: 9780743273565",
			expectedQuery: "isbn:9780743273565",
			expectedType:  bookid.SearchTypeISBN,
			expectedISBN:  "9780743273565",
		},
		{
			name:          "title_and_author",
			input:         "The Great Gatsby by F. Scott Fitzgerald",
			expectedQuery: "The Great Gatsby by F. Scott Fitzgerald",
			expectedType:  bookid.SearchTypeTitleAuthor,
		},
		{
			name:          "title_only",
			input:         "The Great Gatsby",
			expectedQuery: "The Great Gatsby",
			expectedType:  bookid.SearchTypeTitle,
		},
		{
			name:          "general_query",
			input:         "classic american literature 1920s",
			expectedQuery: "classic american literature 1920s",
			expectedType:  bookid.SearchTypeGeneralQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			query, searchType, isbn := googlebooks.ParseQuery(tt.input)
			assert.Equal(t, tt.expectedQuery, query)
			assert.Equal(t, tt.expectedType, searchType)
			if tt.expectedISBN != "" {
				assert.Equal(t, tt.expectedISBN, isbn)
			}
		})
	}
}
