package editor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Open() (io.ReadCloser, error) {
	f, err := os.CreateTemp("", "*.jounrl")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	defer func() {
		// TODO: we could add a log here, letting the user know if the temporary file
		// could not be deleted
		_ = os.Remove(f.Name())
	}()

	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("error while closing file: %w", err)
	}

	editor := "nvim"
	if e, set := os.LookupEnv("EDITOR"); set {
		editor = e
	}

	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to open editor %q: %w", editor, err)
	}

	f, err = os.Open(f.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open updated file: %w", err)
	}

	return f, nil
}
