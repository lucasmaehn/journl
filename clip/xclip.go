//go:build linux

package clip

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var ErrEmptyClipboard = errors.New("empty clipboard")

type ClipboardType byte

var (
	bin      string
	listArgs []string
	readArgs []string
)

const (
	TypeNone ClipboardType = iota
	TypeImage
	TypeText
	TypeURIList
)

func (ctype ClipboardType) Mimetype() string {
	switch ctype {
	case TypeURIList:
		return "text/uri-list"
	case TypeImage:
		return "image/png"
	case TypeText:
		return "text/plain"
	default:
		return ""
	}
}

var initBackend = sync.OnceFunc(func() {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		bin = "wl-paste"
		listArgs = []string{"--list-types"}
		readArgs = []string{"--type"}
	} else {
		bin = "xclip"
		listArgs = []string{"-selection", "clipboard", "-t", "TARGETS", "-o"}
		readArgs = []string{"-selection", "clipboard", "-t"}
	}

	if _, err := exec.LookPath(bin); err != nil {
		panic(fmt.Sprintf("Error: %s not found. Please install it.\n", bin))
	}
})

// Paste stores the content of the clipboard to one or more temp files and returns the created files
func Paste() ([]string, error) {
	initBackend()

	ctype, err := detectType()
	if err != nil {
		return nil, err
	}

	if ctype == TypeNone {
		return nil, ErrEmptyClipboard
	}

	data, err := readClipboard(ctype)
	if err != nil {
		return nil, err
	}

	if ctype == TypeURIList {
		return copyFiles(string(data))
	} else {
		var ext string
		if ctype == TypeImage {
			ext = ".png"
		}

		f, err := os.CreateTemp("/tmp", "*"+ext)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if _, err := f.Write(data); err != nil {
			return nil, err
		}
		return []string{f.Name()}, nil
	}
}

func detectType() (ClipboardType, error) {
	out, err := exec.Command(bin, listArgs...).Output()
	if err != nil {
		return TypeNone, err
	}
	targets := string(out)

	switch {
	case strings.Contains(targets, "text/uri-list"):
		return TypeURIList, nil
	case strings.Contains(targets, "image/png"):
		return TypeImage, nil
	case strings.Contains(targets, "text/plain") || strings.Contains(targets, "UTF8_STRING"):
		return TypeText, nil
	default:
		return TypeNone, nil

	}
}

func copyFiles(_ string) ([]string, error) {
	return nil, errors.New("copying files is currently unimplemented")
}

func readClipboard(ctype ClipboardType) ([]byte, error) {
	args := append(readArgs, ctype.Mimetype())
	if bin == "xclip" {
		args = append(args, "-o")
	}

	data, err := exec.Command(bin, args...).Output()
	if err != nil {
		return nil, err
	}

	return data, nil
}
