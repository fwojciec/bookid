package bookid

import "time"

// Work represents the abstract creative work (the "platonic" book)
type Work struct {
	ID        int64  // Simple auto-increment ID
	Title     string // As it appears on the title page
	Author    string // As credited on the title page
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Author represents a person who created works
type Author struct {
	ID   int64  // Simple auto-increment ID
	Name string // Normalized name for deduplication
}

// WorkAuthor links works to their authors (for searching/indexing)
type WorkAuthor struct {
	WorkID   int64
	AuthorID int64
}

// Publication represents a specific published edition of a Work
type Publication struct {
	ID                  int64
	WorkID              int64
	ISBN10              string
	ISBN13              string
	Publisher           string
	PublishedYear       int
	Language            string
	GoogleBooksVolumeID string
	ThumbnailURL        string
	GoogleBooksData     string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
