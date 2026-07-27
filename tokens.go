// tokens.go — SQLite-backed device-token store for captcha verification.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

type TokenStore struct {
	mu sync.Mutex
	db *sql.DB
}

// OpenTokenStore opens (or creates) the token database at dbPath.
// If the file does not exist it is created with the correct schema so the
// server starts cleanly even before any tokens have been collected.
func OpenTokenStore(dbPath string) (*TokenStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open token store %s: %w", dbPath, err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// Ensure the tokens table exists (idempotent).
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tokens (
		id    INTEGER PRIMARY KEY AUTOINCREMENT,
		token TEXT    NOT NULL UNIQUE
	);`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("init token store schema: %w", err)
	}

	return &TokenStore{db: db}, nil
}

func (ts *TokenStore) Close() error {
	if ts.db != nil {
		return ts.db.Close()
	}
	return nil
}

// Count returns the number of tokens currently in the store.
// Returns 0 if the store is nil or the query fails.
func (ts *TokenStore) Count() int {
	if ts == nil || ts.db == nil {
		return 0
	}
	var n int
	_ = ts.db.QueryRow("SELECT COUNT(*) FROM tokens;").Scan(&n)
	return n
}

func (ts *TokenStore) getNext() (string, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	var token string
	err := ts.db.QueryRow("SELECT token FROM tokens ORDER BY id LIMIT 1;").Scan(&token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			captchaLog("No device tokens available in table 'tokens'")
		} else {
			captchaLog("Failed to query token: " + err.Error())
		}
		return "", false
	}
	return token, true
}

func (ts *TokenStore) remove(token string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	_, err := ts.db.Exec("DELETE FROM tokens WHERE token = ?;", token)
	if err != nil {
		captchaLog("Failed to delete consumed token: " + err.Error())
	}
}