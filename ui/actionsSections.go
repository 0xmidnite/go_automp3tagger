package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileAction int
type DetailsAction int

const (
	START_FETCHING_ALL 			FileAction = iota
	START_FETCHING_SELECTED 	
	STOP_FETCHING 			
)

const (
	APPLY_DISCOGS_RESPONSE_ALL 		DetailsAction = iota
	EDIT_FIELD_ASSOCIATION
	EDIT_FIELD_VALUE
)

var MapActionToOption = map[FileAction]string{
	START_FETCHING_ALL: "Start Fetching All",
	START_FETCHING_SELECTED: "Start Fetching Selected",
	STOP_FETCHING: "Stop Fetching",
}

var MapDetailsActionToOption = map[DetailsAction]string{	
	APPLY_DISCOGS_RESPONSE_ALL: "Apply Discogs Response",
	EDIT_FIELD_ASSOCIATION: "Edit Field Association",
	EDIT_FIELD_VALUE: "Edit Field Value",
}

type ActionSectionModel struct {
	Cursor 					int
	IsFocusedOnMusicTable 	bool
	IsFocusedOnDetailsTable 	bool
	Log 					string
}

var SelectedStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")). // White text
	Background(lipgloss.Color("#7D56F4")). // Purple background
	Bold(true).
	Padding(0, 1)

var NormalStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA")). // White text
	Padding(0, 1)

var LogStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#C60101")). // Red text
	Padding(0, 1)

type ActionMsg struct {
	Action FileAction
}

type DetailsActionMsg struct {
	Action DetailsAction
}

func SendAction(action FileAction) tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{
			Action: action,
		}
	}
}

func SendDetailsAction(action DetailsAction) tea.Cmd {
	return func() tea.Msg {
		return DetailsActionMsg{
			Action: action,
		}
	}
}

func (m ActionSectionModel) SetLog(log string) ActionSectionModel {
	m.Log = log
	return m
}

func InitActionSection() ActionSectionModel {
	return ActionSectionModel{
		Cursor: 0,
		IsFocusedOnMusicTable: true,
		IsFocusedOnDetailsTable: false,
	}
}

func (m ActionSectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var optionsLength = 0;
	
	if(m.IsFocusedOnMusicTable) {
		optionsLength = len(MapActionToOption)
	} else {
		optionsLength = len(MapDetailsActionToOption)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			if m.Cursor > optionsLength {
				m.Cursor = optionsLength - 1
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
			
			if(m.IsFocusedOnMusicTable) {
				var option = FileAction(m.Cursor)

				switch option {
				case START_FETCHING_ALL:
					return m, SendAction(START_FETCHING_ALL)
				case START_FETCHING_SELECTED:
					return m, SendAction(START_FETCHING_SELECTED)
				case STOP_FETCHING:
					return m, SendAction(STOP_FETCHING)	
				}
			} 

			if(m.IsFocusedOnDetailsTable) {
				var option = DetailsAction(m.Cursor)

				switch option {
				case APPLY_DISCOGS_RESPONSE_ALL:
					return m, SendDetailsAction(APPLY_DISCOGS_RESPONSE_ALL)
				case EDIT_FIELD_ASSOCIATION:
					return m, SendDetailsAction(EDIT_FIELD_ASSOCIATION)
				case EDIT_FIELD_VALUE:
					return m, SendDetailsAction(EDIT_FIELD_VALUE)
				}
			}

			return m, nil
		}
	}

	return m, nil
}


func (m ActionSectionModel) View() string {
	var s string = ""

	if(m.IsFocusedOnMusicTable) {
		var actions []FileAction = []FileAction{
			START_FETCHING_ALL,
			START_FETCHING_SELECTED,
			STOP_FETCHING,
		}

		for i, option := range actions {
			if i == m.Cursor {
				s += SelectedStyle.Render(MapActionToOption[option])
			} else {
				s += NormalStyle.Render(MapActionToOption[option])
			}

			i++
		}
	}

	if(m.IsFocusedOnDetailsTable) {
		var actions []DetailsAction = []DetailsAction{
			APPLY_DISCOGS_RESPONSE_ALL,
			EDIT_FIELD_ASSOCIATION,
			EDIT_FIELD_VALUE,
		}

		for i, option := range actions {
			if i == m.Cursor {
				s += SelectedStyle.Render(MapDetailsActionToOption[option])
			} else {
				s += NormalStyle.Render(MapDetailsActionToOption[option])
			}
		}
	}

	return s + "\n" + LogStyle.Render(m.Log)
}

func (m ActionSectionModel) Init() tea.Cmd {
	return nil
}
