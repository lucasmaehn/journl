package logstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/georgysavva/scany/sqlscan"
	"github.com/lucasmaehn/journl/config"
	_ "modernc.org/sqlite"
)

func NewSQLite(contextName string, cfg config.StoreConfig) (*SQLiteLogstore, error) {
	log.Println("opening sqlite @", cfg.Path)
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	ls := &SQLiteLogstore{
		contextName: contextName,
		db:          db,
	}
	if err := ls.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite adapter: %w", err)
	}

	return ls, nil
}

type SQLiteLogstore struct {
	contextName string
	db          *sql.DB
}

type entryRow struct {
	ID        int       `db:"id"`
	Text      string    `db:"text"`
	Context   string    `db:"context"`
	CreatedAt time.Time `db:"created_at"`
}

func (ls *SQLiteLogstore) init() error {
	query := `
CREATE TABLE IF NOT EXISTS journal_entries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	text TEXT NOT NULL,
	context TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);`

	_, err := ls.db.Exec(query)
	return err
}

func (ls *SQLiteLogstore) Commit(text string, opts ...LogOption) error {
	fmt.Println("log", text)
	cfg := entryConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	if len(cfg.attachments) > 0 {
		return errors.New("attachments are not supported")
	}

	if len(cfg.stdin) > 0 {
		text += "\n" + cfg.stdin
	}

	query := `
INSERT INTO journal_entries (text, context)
VALUES (?, ?)
	`
	_, err := ls.db.Exec(query, text, ls.contextName)
	if err != nil {
		return fmt.Errorf("failed to commit journal entry: %w", err)
	}

	return nil
}

func (ls *SQLiteLogstore) List() ([]LogEntry, error) {
	var rows []entryRow
	query := `SELECT created_at, text, context FROM journal_entries WHERE context=? ORDER BY created_at desc`
	err := sqlscan.Select(context.Background(), ls.db, &rows, query, ls.contextName)
	if err != nil {
		return nil, err
	}

	var entries []LogEntry
	for _, row := range rows {
		entries = append(entries, LogEntry{
			Text:      row.Text,
			Context:   row.Context,
			Timestamp: row.CreatedAt,
		})
	}

	return entries, nil
}
