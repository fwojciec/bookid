package sqlite_test

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/fwojciec/bookid/sqlite"
)

var dump = flag.Bool("dump", false, "save work data")

// Ensure the test database can open & close.
func TestDB(t *testing.T) {
	t.Parallel()
	db := MustOpenDB(t)
	MustCloseDB(t, db)
}

// MustOpenDB returns a new, open DB. Fatal on error.
func MustOpenDB(tb testing.TB) *sqlite.DB {
	tb.Helper()

	// Write to an in-memory database by default.
	// If the -dump flag is set, generate a temp file for the database.
	dsn := ":memory:"
	if *dump {
		dir, err := os.MkdirTemp("", "")
		if err != nil {
			tb.Fatal(err)
		}
		dsn = filepath.Join(dir, "db")
		tb.Logf("DUMP=%s", dsn)
	}

	db := sqlite.NewDB(dsn)
	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}
	return db
}

// MustCloseDB closes the DB. Fatal on error.
func MustCloseDB(tb testing.TB, db *sqlite.DB) {
	tb.Helper()
	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}
