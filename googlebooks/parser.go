package googlebooks

import (
	"regexp"
	"strings"

	"github.com/fwojciec/bookid"
)

var (
	// ISBN-10: exactly 10 digits (with optional dashes)
	isbn10Pattern = regexp.MustCompile(`\b(\d{1,5}[-\s]?\d{1,7}[-\s]?\d{1,7}[-\s]?\d)\b`)

	// ISBN-13: exactly 13 digits starting with 978 or 979 (with optional dashes)
	isbn13Pattern = regexp.MustCompile(`\b(97[89][-\s]?\d{1,5}[-\s]?\d{1,7}[-\s]?\d{1,7}[-\s]?\d)\b`)
)

// ParseQuery analyzes the input string and returns an appropriate Google Books API query
func ParseQuery(input string) (query string, searchType bookid.SearchType, isbn string) {
	// First, check for ISBN
	cleanInput := strings.TrimSpace(input)

	// Check for ISBN-13 first (more specific)
	if matches := isbn13Pattern.FindStringSubmatch(cleanInput); len(matches) > 0 {
		isbn = cleanISBN(matches[1])
		if validateISBN13(isbn) {
			return "isbn:" + isbn, bookid.SearchTypeISBN, isbn
		}
	}

	// Check for ISBN-10
	if matches := isbn10Pattern.FindStringSubmatch(cleanInput); len(matches) > 0 {
		isbn = cleanISBN(matches[1])
		if validateISBN10(isbn) {
			return "isbn:" + isbn, bookid.SearchTypeISBN, isbn
		}
	}

	// For all non-ISBN queries, use natural language search
	// Google Books API handles fuzzy matching better than strict operators
	// This allows for variations in title/author spelling and formatting
	return cleanInput, bookid.SearchTypeGeneralQuery, ""
}

// cleanISBN removes dashes and spaces from ISBN
func cleanISBN(isbn string) string {
	isbn = strings.ReplaceAll(isbn, "-", "")
	isbn = strings.ReplaceAll(isbn, " ", "")
	return isbn
}

// validateISBN10 checks if the ISBN-10 is valid (basic length check)
func validateISBN10(isbn string) bool {
	return len(isbn) == 10 && isAllDigits(isbn)
}

// validateISBN13 checks if the ISBN-13 is valid
func validateISBN13(isbn string) bool {
	if len(isbn) != 13 || !isAllDigits(isbn) {
		return false
	}
	// Must start with 978 or 979
	return strings.HasPrefix(isbn, "978") || strings.HasPrefix(isbn, "979")
}

// isAllDigits checks if string contains only digits
func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
