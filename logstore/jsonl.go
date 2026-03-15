package logstore

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"time"
)

func NewJSONL() (*JSONLStore, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	filepath := path.Join(homeDir, ".journl", "default.jsonl")

	if err := os.MkdirAll(path.Dir(filepath), 0o755); err != nil {
		return nil, err
	}

	return &JSONLStore{
		filepath: filepath,
	}, nil
}

type JSONLStore struct {
	filepath string
}

func (ls *JSONLStore) List() ([]LogEntry, error) {
	f, err := os.Open(ls.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(f)

	i := 0
	for scanner.Scan() {
		i += 1
		entry := LogEntry{}
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			log.Printf("failed to read line %d: %v", i, err)
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (ls *JSONLStore) Commit(text string, opts ...LogOption) error {
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

	entry := LogEntry{
		Timestamp: time.Now(),
		Text:      text,
	}

	f, err := os.OpenFile(ls.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	if _, err := f.WriteString(string(bs) + "\n"); err != nil {
		return err
	}

	return nil
}
