package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fwojciec/bookid"
	"github.com/fwojciec/bookid/googlebooks"
)

const (
	defaultTimeout = 30 * time.Second
)

type Config struct {
	GoogleBooksAPIKey string
	Timeout           time.Duration
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("usage: bookid <search query>")
	}

	// Combine all arguments after the program name as the search query
	query := strings.Join(os.Args[1:], " ")

	// Load configuration
	config := loadConfig()

	// Create Google Books client
	client, err := googlebooks.NewClient(config.GoogleBooksAPIKey)
	if err != nil {
		return fmt.Errorf("creating Google Books client: %w", err)
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Perform search
	results, err := client.Search(ctx, query)
	if err != nil {
		return fmt.Errorf("searching for books: %w", err)
	}

	// Return just the top result if any results were found
	if len(results) == 0 {
		// Return empty response
		response := struct {
			Query  string             `json:"query"`
			Result *bookid.BookResult `json:"result"`
		}{
			Query:  query,
			Result: nil,
		}
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)
		return encoder.Encode(response)
	}

	// Get the top result and strip out the raw Google Books data
	topResult := results[0]
	topResult.GoogleBooksData = nil

	response := struct {
		Query  string            `json:"query"`
		Result bookid.BookResult `json:"result"`
	}{
		Query:  query,
		Result: topResult,
	}

	// Output as pretty-printed JSON without HTML escaping
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(response); err != nil {
		return fmt.Errorf("encoding JSON output: %w", err)
	}

	return nil
}

func loadConfig() Config {
	config := Config{
		GoogleBooksAPIKey: os.Getenv("GOOGLE_BOOKS_API_KEY"),
		Timeout:           defaultTimeout,
	}

	// Allow timeout override via environment variable
	if timeoutStr := os.Getenv("BOOKID_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			config.Timeout = timeout
		}
	}

	return config
}
