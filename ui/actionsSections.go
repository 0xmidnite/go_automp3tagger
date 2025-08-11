package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileAction int

const (
	START_FETCHING_ALL 			FileAction = iota
	START_FETCHING_SELECTED 	
	STOP_FETCHING 			
)

var MapActionToOption = map[FileAction]string{
	START_FETCHING_ALL: "Start Fetching All",
	START_FETCHING_SELECTED: "Start Fetching Selected",
	STOP_FETCHING: "Stop Fetching",
}

type ActionSectionModel struct {
	Cursor 					int
	Options 				[]FileAction
	IsFocusedOnMusicTable 	bool
	IsFocusedOnFileTable 	bool
}


func (m ActionSectionModel) Init() tea.Cmd {
	return nil
}

var SelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")). // White text
	Background(lipgloss.Color("#7D56F4")). // Purple background
	Bold(true).
	Padding(0, 1)

var NormalStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")). // White text
	Padding(0, 1)

type ActionMsg struct {
	Action FileAction
}

func SendAction(action FileAction) tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{
			Action: action,
		}
	}
}

func (m ActionSectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			if m.Cursor + 1 >= len(m.Options) {
				m.Cursor = len(m.Options) - 1
			}else{
				m.Cursor++
			}

			return m, nil

		case "left":
			if m.Cursor - 1 < 0 {
				m.Cursor = 0
			}else{
				m.Cursor--
			}

			return m, nil

		case "enter":
			var option = m.Options[m.Cursor]

			switch option {
			case START_FETCHING_ALL:
				return m, SendAction(START_FETCHING_ALL)
			case START_FETCHING_SELECTED:
				return m, SendAction(START_FETCHING_SELECTED)
			case STOP_FETCHING:
				return m, SendAction(STOP_FETCHING)	
			}

			return m, nil
		}
	}

	return m, nil
}


func (m ActionSectionModel) View() string {
	var s = ""

	for i, option := range m.Options {
		if i == m.Cursor {
			s += SelectedStyle.Render(MapActionToOption[option])
		} else {
			s += NormalStyle.Render(MapActionToOption[option])
		}
	}

	return s
}

func InitActionSection() ActionSectionModel {
	var actions = []FileAction{
		START_FETCHING_ALL,
		START_FETCHING_SELECTED,
		STOP_FETCHING,
	}

	return ActionSectionModel{
		Cursor: 0,
		Options: actions,
		IsFocusedOnMusicTable: true,
		IsFocusedOnFileTable: false,
	}
}