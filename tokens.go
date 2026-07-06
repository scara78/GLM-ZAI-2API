// tokens.go — SQLite-backed device-token store for captcha verification.
// Replaces the globalDB + dbMu pattern from the original token-helper.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sync"

	_ "modernc.org/sqlite"
)

type TokenStore struct {
	mu sync.Mutex
	db *sql.DB
}

func OpenTokenStore(dbPath string) (*TokenStore, error) {
	if _, err := os.Stat(dbPath); err != nil {
		return nil, fmt.Errorf("token database not found: %s", dbPath)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return &TokenStore{db: db}, nil
}

func (ts *TokenStore) Close() error {
	if ts.db != nil {
		return ts.db.Close()
	}
	return nil
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
