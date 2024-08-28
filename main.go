package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
	todo status = iota
	inProgress
	done
)

// Styling
var (
	columnStyle  = lipgloss.NewStyle().Padding(1, 2)
	focusedStyle = lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type Task struct {
	status      status
	title       string
	description string
}

func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

type Model struct {
	focused  status
	lists    []list.Model
	err      error
	loaded   bool
	quitting bool
}

func New() *Model {
	log.Println("Initializing Bubble Tea model")
	// Initialize list with some default dimensions (can be adjusted as needed)
	m := &Model{}
	m.initLists(90, 32) // Default dimensions
	return m
}

// Go to previous list
func (m *Model) Previous() {
	if m.focused == todo {
		m.focused = done
	} else {
		m.focused--
	}
}

// Go to next list
func (m *Model) Next() {
	if m.focused == done {
		m.focused = todo
	} else {
		m.focused++
	}
}

func (m *Model) initLists(width, height int) {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width/divisor, height)
	defaultList.SetShowHelp(false)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	// Init to do
	m.lists[todo].Title = "To Do"
	m.lists[todo].SetItems([]list.Item{
		Task{status: todo, title: "Buy oats", description: "organic oats"},
		Task{status: todo, title: "Buy mate", description: "tony mate"},
		Task{status: todo, title: "do laundry", description: "30 degrees C and dry outdoors"},
	})
	// Init in progress
	m.lists[inProgress].Title = "In progress"
	m.lists[inProgress].SetItems([]list.Item{
		Task{status: inProgress, title: "learn bubble tea", description: "building a TUI"},
	})

	// Init done
	m.lists[done].Title = "Done"
	m.lists[done].SetItems([]list.Item{
		Task{status: done, title: "Install Go", description: "New machine no problems."},
	})
}

func (m Model) Init() tea.Cmd {
	// Optionally, you can return some initialization command here
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			columnStyle.Width(msg.Width / divisor)
			focusedStyle.Width(msg.Width / divisor)
			columnStyle.Height(msg.Height - divisor)
			focusedStyle.Height(msg.Height - divisor)
			smaller := .5
			m.initLists(msg.Width, int(float64(msg.Height)*smaller))
			m.loaded = true
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Previous()
		case "right", "l":
			m.Next()
		}

	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.loaded {
		todoView := m.lists[todo].View()
		inProgressView := m.lists[inProgress].View()
		doneView := m.lists[done].View()
		switch m.focused {
		case inProgress:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				focusedStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		case done:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				columnStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				focusedStyle.Render(doneView),
			)
		default:
			return lipgloss.JoinHorizontal(
				lipgloss.Left,
				focusedStyle.Render(todoView),
				columnStyle.Render(inProgressView),
				columnStyle.Render(doneView),
			)
		}

	} else {
		return "loading"
	}
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := New()
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
