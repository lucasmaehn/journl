package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasmaehn/journl/logstore"
)

// UI State Constants
type state int

const (
	listView state = iota
)

var (
	docStyle       = lipgloss.NewStyle().Margin(1, 2)
	contextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	timestampStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	headerStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true).BorderForeground(lipgloss.Color("205"))
)

// LogEntry is your existing struct (re-declared or imported)
type LogEntry struct {
	Timestamp interface{} // Use your actual time.Time
	Text      string
	Context   string
}

type item struct {
	entry logstore.LogEntry
}

func (i item) Title() string       { return strings.Split(i.entry.Text, "\n")[0] }
func (i item) Description() string { return fmt.Sprintf("%s | %v", i.entry.Context, i.entry.Timestamp) }
func (i item) FilterValue() string { return i.entry.Text + " " + i.entry.Context }

type model struct {
	state    state
	list     list.Model
	viewport viewport.Model
	ready    bool
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc", "backspace":
			return m, tea.Quit
		case "enter":
			if m.state == listView {
				selected := m.list.SelectedItem().(item)
				openReadOnly(selected.entry.Text)
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		m.viewport.Width = msg.Width - h
		m.viewport.Height = msg.Height - v - 3 // Leave room for header
	}

	// Route updates based on current state
	if m.state == listView {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func openReadOnly(content string) error {
	tmpFile, err := os.CreateTemp("", "journl-view-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		return err
	}
	tmpFile.Close()

	// Use -R for readonly or -M for non-modifiable
	cmd := exec.Command("nvim", "-R", tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RenderJournal(entries []logstore.LogEntry) {
	items := make([]list.Item, len(entries))
	for i, e := range entries {
		items[i] = item{entry: e}
	}

	m := model{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		viewport: viewport.New(0, 0),
		state:    listView,
	}
	m.list.Title = "Journl Entries"

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
