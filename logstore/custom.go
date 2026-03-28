package logstore

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/lucasmaehn/journl/config"
)

type CustomStore struct {
	path        string
	templ       *template.Template
	contextName string
}

type entryTemplateData struct {
	Title     string
	Body      string
	Context   string
	Timestamp time.Time
}

func newEntryTemplateData(contextName, text string) entryTemplateData {
	title, body := extractTitle(text)
	return entryTemplateData{
		Title:     title,
		Body:      body,
		Context:   contextName,
		Timestamp: time.Now(),
	}
}

func extractTitle(text string) (title, body string) {
	lines := strings.SplitN(text, "\n", 3)
	if len(lines) == 3 {
		if strings.TrimSpace(lines[1]) == "" {
			return strings.TrimSpace(lines[0]), lines[2]
		}
	}

	return "", text
}

func NewCustom(contextName string, cfg config.StoreConfig) (*CustomStore, error) {
	templ, err := template.New("entry").Parse(cfg.Custom.Template)
	if err != nil {
		return nil, fmt.Errorf("parse entry template: %w", err)
	}

	return &CustomStore{
		path:        cfg.Path,
		templ:       templ,
		contextName: contextName,
	}, nil
}

func (cs *CustomStore) Commit(text string, opts ...LogOption) error {
	path, err := resolvePath(cs.path)
	log.Println("using path", path)
	if err != nil {
		return fmt.Errorf("failed to resolve filepath: %w", err)
	}

	var buf bytes.Buffer

	data := newEntryTemplateData(cs.contextName, text)
	log.Println("templ data", data)

	if err := cs.templ.Execute(&buf, data); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Println("cur string", buf.String())

	if _, err := f.WriteString(buf.String()); err != nil {
		return err
	}

	return nil
}

func (cs *CustomStore) List() ([]LogEntry, error) {
	return nil, errors.New("list is not supported for custom store")
}
