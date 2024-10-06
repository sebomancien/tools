package zhe

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type model struct {
	quitting  bool
	progress  progress.Model
	stopwatch stopwatch.Model
}

func newModel() *model {
	return &model{
		quitting:  false,
		progress:  progress.New(progress.WithDefaultGradient()),
		stopwatch: stopwatch.NewWithInterval(100 * time.Millisecond),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.progress.Init(), m.stopwatch.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	quit := false

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			quit = true
			m.stopwatch.Stop()
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width

	case Progress:
		cmds = append(cmds, m.progress.SetPercent(float64(msg.Counter)/float64(msg.Total)))
		if msg.Counter == msg.Total {
			quit = true
		}
	}

	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	cmds = append(cmds, cmd)

	progressModel, cmd := m.progress.Update(msg)
	m.progress = progressModel.(progress.Model)
	cmds = append(cmds, cmd)

	cmd = tea.Batch(cmds...)
	if quit && !m.quitting {
		cmd = tea.Sequence(cmd, tea.Quit)
		m.quitting = true
	}
	return m, cmd
}

func (m model) View() string {
	str := fmt.Sprintln()
	if m.quitting {
		str += fmt.Sprintln(m.progress.ViewAs(1.0))
	} else {
		str += fmt.Sprintln(m.progress.View())
	}
	str += fmt.Sprintln("Ellapsed time:", m.stopwatch.View())
	str += fmt.Sprintln()
	if !m.quitting {
		str += fmt.Sprintln(helpStyle("Press q to quit"))
	}
	return str
}

func NewTui() *tea.Program {
	m := newModel()
	return tea.NewProgram(m)
}
