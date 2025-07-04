package googlebooks

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/fwojciec/bookid"
	"google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

// Client implements the BookFinder interface for Google Books API
type Client struct {
	service *books.Service
}

// NewClient creates a new Google Books API client
func NewClient(apiKey string) (*Client, error) {
	ctx := context.Background()

	opts := []option.ClientOption{}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	} else {
		// Explicitly disable authentication when no API key is provided
		opts = append(opts, option.WithoutAuthentication())
	}

	service, err := books.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		service: service,
	}, nil
}

// NewClientWithService creates a new client with a custom service (for testing)
func NewClientWithService(service *books.Service) *Client {
	return &Client{
		service: service,
	}
}

// Search performs a book search based on the provided query
func (c *Client) Search(ctx context.Context, query string) ([]bookid.BookResult, error) {
	if query == "" {
		return nil, errors.New("query cannot be empty")
	}

	// Parse the query to determine search type
	searchQuery, searchType, detectedISBN := ParseQuery(query)

	// Build and execute the search
	call := c.service.Volumes.List(searchQuery)
	call.MaxResults(10)
	call.Context(ctx)

	resp, err := call.Do()
	if err != nil {
		// Return the original error to preserve context
		// Callers can check for googleapi.Error if they need specific handling
		return nil, err
	}

	// Convert to BookResult
	results := make([]bookid.BookResult, 0, len(resp.Items))
	for _, volume := range resp.Items {
		// Marshal the volume to JSON for GoogleBooksData field
		volumeJSON, err := json.Marshal(volume)
		if err != nil {
			// Log error but continue - don't fail the whole search
			volumeJSON = nil
		}

		result := volumeToBookResult(volume, searchType, detectedISBN)
		result.GoogleBooksData = volumeJSON
		results = append(results, result)
	}

	return results, nil
}

// volumeToBookResult converts a Google Books Volume to our BookResult
func volumeToBookResult(volume *books.Volume, searchType bookid.SearchType, detectedISBN string) bookid.BookResult {
	result := bookid.BookResult{
		Title:               volume.VolumeInfo.Title,
		Authors:             volume.VolumeInfo.Authors,
		GoogleBooksVolumeID: volume.Id,
		SearchType:          searchType,
	}

	// Extract ISBNs
	for _, identifier := range volume.VolumeInfo.IndustryIdentifiers {
		switch identifier.Type {
		case "ISBN_10":
			result.ISBN10 = identifier.Identifier
		case "ISBN_13":
			result.ISBN13 = identifier.Identifier
		}
	}

	// If we detected an ISBN in the query but didn't find it in identifiers,
	// it might still be the right book
	if detectedISBN != "" && result.ISBN10 == "" && result.ISBN13 == "" {
		if len(detectedISBN) == 10 {
			result.ISBN10 = detectedISBN
		} else if len(detectedISBN) == 13 {
			result.ISBN13 = detectedISBN
		}
	}

	// Extract publication details
	result.Publisher = volume.VolumeInfo.Publisher
	result.Language = volume.VolumeInfo.Language

	// Parse published year
	if volume.VolumeInfo.PublishedDate != "" {
		year := extractYear(volume.VolumeInfo.PublishedDate)
		result.PublishedYear = year
	}

	// Extract thumbnail URL and ensure HTTPS
	if volume.VolumeInfo.ImageLinks != nil {
		if volume.VolumeInfo.ImageLinks.Thumbnail != "" {
			result.ThumbnailURL = ensureHTTPS(volume.VolumeInfo.ImageLinks.Thumbnail)
		} else if volume.VolumeInfo.ImageLinks.SmallThumbnail != "" {
			result.ThumbnailURL = ensureHTTPS(volume.VolumeInfo.ImageLinks.SmallThumbnail)
		}
	}

	// Calculate confidence score
	result.Confidence = calculateConfidence(searchType, volume)

	return result
}

// calculateConfidence determines the confidence score based on search type and result quality
func calculateConfidence(searchType bookid.SearchType, volume *books.Volume) float64 {
	baseConfidence := map[bookid.SearchType]float64{
		bookid.SearchTypeISBN:         0.95,
		bookid.SearchTypeGeneralQuery: 0.70,
	}

	confidence := baseConfidence[searchType]

	// Adjust based on data completeness
	if volume.VolumeInfo != nil {
		dataPoints := 0
		totalPoints := 0

		// Check important fields
		if volume.VolumeInfo.Title != "" {
			dataPoints++
		}
		totalPoints++

		if len(volume.VolumeInfo.Authors) > 0 {
			dataPoints++
		}
		totalPoints++

		if len(volume.VolumeInfo.IndustryIdentifiers) > 0 {
			dataPoints++
		}
		totalPoints++

		if volume.VolumeInfo.Publisher != "" {
			dataPoints++
		}
		totalPoints++

		// Adjust confidence based on data completeness
		completeness := float64(dataPoints) / float64(totalPoints)
		confidence = confidence * (0.7 + 0.3*completeness)
	}

	return confidence
}

// extractYear extracts the year from various date formats
func extractYear(dateStr string) int {
	// Try to parse as year only
	if year, err := strconv.Atoi(dateStr); err == nil && year > 1000 && year < 3000 {
		return year
	}

	// Try to extract year from date string
	parts := strings.Split(dateStr, "-")
	if len(parts) > 0 {
		if year, err := strconv.Atoi(parts[0]); err == nil && year > 1000 && year < 3000 {
			return year
		}
	}

	return 0
}

// ensureHTTPS converts an HTTP URL to HTTPS by replacing only the leading scheme.
// If the URL already uses HTTPS or is empty, it is returned unchanged.
// This function only modifies the URL scheme and preserves any other occurrences
// of "http://" that might appear elsewhere in the URL (e.g., in query parameters).
func ensureHTTPS(url string) string {
	if url == "" {
		return ""
	}
	// Replace only the URL scheme if it starts with http://
	if strings.HasPrefix(url, "http://") {
		return "https://" + url[len("http://"):]
	}
	return url
}
