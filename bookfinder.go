package bookid

import (
	"context"
	"encoding/json"
)

// BookFinder searches for books and returns detailed results
type BookFinder interface {
	// Search performs a book search based on the provided query
	// Returns a list of BookResult with confidence scores
	Search(ctx context.Context, query string) ([]BookResult, error)
}

// BookResult contains all information needed to create Work, Author, and Publication
type BookResult struct {
	// For Work creation
	Title   string   `json:"title"`
	Authors []string `json:"authors"`

	// For Publication creation
	ISBN10              string          `json:"isbn10,omitempty"`
	ISBN13              string          `json:"isbn13,omitempty"`
	Publisher           string          `json:"publisher,omitempty"`
	PublishedYear       int             `json:"published_year,omitempty"`
	Language            string          `json:"language,omitempty"`
	GoogleBooksVolumeID string          `json:"google_books_volume_id,omitempty"`
	ThumbnailURL        string          `json:"thumbnail_url,omitempty"`
	GoogleBooksData     json.RawMessage `json:"google_books_data,omitempty"` // Raw API response

	// Search metadata
	Confidence float64    `json:"confidence"` // 0.0 to 1.0
	SearchType SearchType `json:"search_type"`
}

// SearchType indicates how the search was performed
type SearchType string

const (
	SearchTypeISBN         SearchType = "isbn"
	SearchTypeTitleAuthor  SearchType = "title_author"
	SearchTypeTitle        SearchType = "title"
	SearchTypeGeneralQuery SearchType = "general"
)
