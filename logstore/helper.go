package logstore

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func resolvePath(path string) (string, error) {
	path = time.Now().Format(path)
	return expandPath(path)
}

func expandPath(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("expanding path: %w", err)
	}

	return filepath.Join(home, path[1:]), nil
}
