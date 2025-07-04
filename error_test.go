package bookid_test

import (
	"errors"
	"testing"

	"github.com/fwojciec/bookid"
)

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		e := &bookid.Error{Code: bookid.EINVALID, Message: "Invalid input"}
		if got, want := e.Error(), "bookid error: code=invalid message=Invalid input"; got != want {
			t.Fatalf("Error()=%q, want %q", got, want)
		}
	})

	t.Run("ErrorCode", func(t *testing.T) {
		t.Parallel()
		// Application error
		e := &bookid.Error{Code: bookid.ENOTFOUND, Message: "Not found"}
		if got, want := bookid.ErrorCode(e), bookid.ENOTFOUND; got != want {
			t.Fatalf("ErrorCode()=%q, want %q", got, want)
		}

		// Non-application error
		if got, want := bookid.ErrorCode(errors.New("disk error")), bookid.EINTERNAL; got != want {
			t.Fatalf("ErrorCode()=%q, want %q", got, want)
		}

		// Nil error
		if got, want := bookid.ErrorCode(nil), ""; got != want {
			t.Fatalf("ErrorCode()=%q, want %q", got, want)
		}
	})

	t.Run("ErrorMessage", func(t *testing.T) {
		t.Parallel()
		// Application error
		e := &bookid.Error{Code: bookid.ENOTFOUND, Message: "User not found"}
		if got, want := bookid.ErrorMessage(e), "User not found"; got != want {
			t.Fatalf("ErrorMessage()=%q, want %q", got, want)
		}

		// Non-application error
		if got, want := bookid.ErrorMessage(errors.New("disk error")), "Internal error."; got != want {
			t.Fatalf("ErrorMessage()=%q, want %q", got, want)
		}

		// Nil error
		if got, want := bookid.ErrorMessage(nil), ""; got != want {
			t.Fatalf("ErrorMessage()=%q, want %q", got, want)
		}
	})

	t.Run("Errorf", func(t *testing.T) {
		t.Parallel()
		e := bookid.Errorf(bookid.EINVALID, "Invalid field: %s", "email")
		if got, want := e.Code, bookid.EINVALID; got != want {
			t.Fatalf("Code=%q, want %q", got, want)
		}
		if got, want := e.Message, "Invalid field: email"; got != want {
			t.Fatalf("Message=%q, want %q", got, want)
		}
	})
}
