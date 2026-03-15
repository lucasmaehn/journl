package logstore

import (
	"fmt"
	"time"
)

type LogStore interface {
	Commit(text string, opts ...LogOption) error
	List() ([]LogEntry, error)
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
	Context   string    `json:"context"`
}

func (e LogEntry) String() string {
	return fmt.Sprintf("---\n@%v\n%v\n---\n", e.Timestamp, e.Text)
}

type entryConfig struct {
	attachments []string
	stdin       string
}

type LogOption func(*entryConfig)

func WithAttachment(filepath string) LogOption {
	return func(c *entryConfig) {
		c.attachments = append(c.attachments, filepath)
	}
}

func WithStdin(stdin string) LogOption {
	return func(c *entryConfig) {
		c.stdin = stdin
	}
}
