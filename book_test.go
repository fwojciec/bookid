package bookid_test

import (
	"testing"
	"time"

	"github.com/fwojciec/bookid"
	"github.com/stretchr/testify/assert"
)

func TestWork(t *testing.T) {
	t.Parallel()
	t.Run("create work with all fields", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		work := &bookid.Work{
			ID:        1,
			Title:     "The Go Programming Language",
			Author:    "Alan A. A. Donovan and Brian W. Kernighan",
			CreatedAt: now,
			UpdatedAt: now,
		}

		assert.Equal(t, int64(1), work.ID)
		assert.Equal(t, "The Go Programming Language", work.Title)
		assert.Equal(t, "Alan A. A. Donovan and Brian W. Kernighan", work.Author)
		assert.Equal(t, now, work.CreatedAt)
		assert.Equal(t, now, work.UpdatedAt)
	})
}

func TestAuthor(t *testing.T) {
	t.Parallel()
	t.Run("create author with all fields", func(t *testing.T) {
		t.Parallel()
		author := &bookid.Author{
			ID:   1,
			Name: "Brian W. Kernighan",
		}

		assert.Equal(t, int64(1), author.ID)
		assert.Equal(t, "Brian W. Kernighan", author.Name)
	})
}

func TestWorkAuthor(t *testing.T) {
	t.Parallel()
	t.Run("create work-author relationship", func(t *testing.T) {
		t.Parallel()
		wa := &bookid.WorkAuthor{
			WorkID:   1,
			AuthorID: 2,
		}

		assert.Equal(t, int64(1), wa.WorkID)
		assert.Equal(t, int64(2), wa.AuthorID)
	})
}

func TestPublication(t *testing.T) {
	t.Parallel()
	t.Run("create publication with all fields", func(t *testing.T) {
		t.Parallel()
		now := time.Now()
		pub := &bookid.Publication{
			ID:                  1,
			WorkID:              1,
			ISBN10:              "0134190440",
			ISBN13:              "978-0134190440",
			Publisher:           "Addison-Wesley Professional",
			PublishedYear:       2015,
			Language:            "en",
			GoogleBooksVolumeID: "SJHvCgAAQBAJ",
			ThumbnailURL:        "https://books.google.com/books/content?id=SJHvCgAAQBAJ",
			GoogleBooksData:     `{"title":"The Go Programming Language"}`,
			CreatedAt:           now,
			UpdatedAt:           now,
		}

		assert.Equal(t, int64(1), pub.ID)
		assert.Equal(t, int64(1), pub.WorkID)
		assert.Equal(t, "0134190440", pub.ISBN10)
		assert.Equal(t, "978-0134190440", pub.ISBN13)
		assert.Equal(t, "Addison-Wesley Professional", pub.Publisher)
		assert.Equal(t, 2015, pub.PublishedYear)
		assert.Equal(t, "en", pub.Language)
		assert.Equal(t, "SJHvCgAAQBAJ", pub.GoogleBooksVolumeID)
		assert.Equal(t, "https://books.google.com/books/content?id=SJHvCgAAQBAJ", pub.ThumbnailURL)
		assert.JSONEq(t, `{"title":"The Go Programming Language"}`, pub.GoogleBooksData)
		assert.Equal(t, now, pub.CreatedAt)
		assert.Equal(t, now, pub.UpdatedAt)
	})
}
