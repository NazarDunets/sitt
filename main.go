package main

import (
	"fmt"
	"os"
	"time"

	"github.com/NazarDunets/sitt/command"
	sdl "github.com/NazarDunets/sitt/schedule"
	"github.com/NazarDunets/sitt/storage"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	schedule   *sdl.Schedule
	input      textinput.Model
	date       time.Time
	dateString string
	err        error
}

func main() {
	time := time.Now()
	schedule, err := storage.LoadOrCreateSchedule(time)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(schedule))
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.err = nil
			entry, err := command.GenerateEntryFromCommand(m.input.Value(), m.schedule.FilledUpTo())

			if err == nil {
				m.schedule.Insert(*entry)
				m.err = storage.Save(m.date, m.schedule)
			} else {
				m.err = err
			}

			m.input.SetValue("")
			return m, nil
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	s := "SITT v0\n"
	s += m.dateString
	s += "\n\n"

	if m.err != nil {
		s += "ERROR: "
		s += m.err.Error()
		s += "\n\n"
	}

	s += m.schedule.View()
	s += "\n"
	s += m.input.View()

	return s
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func initialModel(schedule *sdl.Schedule) model {
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 100

	date := time.Now()
	dateString := date.Format("2006-01-02")

	return model{
		schedule:   schedule,
		input:      ti,
		date:       date,
		dateString: dateString,
	}
}
