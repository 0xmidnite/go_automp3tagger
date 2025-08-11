package ui

import (
	ops "automp3tagger/file_ops"
	"strconv"

	discogs "automp3tagger/discogs"

	table "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	lipgloss "github.com/charmbracelet/lipgloss"
)

type MusicRow struct {
	Index 			int
	Extension 		string
	FileName 		string
	HasCompleteID3 	string
	Status 			FileStatus
	DiscogsResults 	[]discogs.DiscogsSearchResult
}

type MusicTableModel struct {
	IsFocused 		bool
	Table 			table.Model

	MusicDir 		string
	Files 			[]MusicRow
	FilesInfo 		[]ops.FileInfo
	CursorIndex 	int
	IndexSelected 	int
	IndexProcessing int
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m MusicTableModel) Init() tea.Cmd {
	return nil
}

type SetFileMsg struct {
	SelectedIndex int
}

type UpdateSpinnerMsg struct {
	Index int
}

func SetFileInfo(selectedIndex int) tea.Cmd {
	return func() tea.Msg {
		return SetFileMsg{
			SelectedIndex: selectedIndex,
		}
	}
}

func (m MusicTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
		case ActionMsg:
			switch msg.Action {
				case START_FETCHING_ALL:
					return m, cmd

				case START_FETCHING_SELECTED:
					var currentFile = m.Files[m.CursorIndex]
				
					if(currentFile.Status != STATUS_PENDING && currentFile.Status != STATUS_FETCH_ERROR) {
						return m, cmd
					}

					var rows = m.Table.Rows()
					rows[m.CursorIndex][4] = FileStatusToString(STATUS_FETCHING)
					m.Files[m.CursorIndex].Status = STATUS_FETCHING

					m.Table.SetRows(rows)
				
					cmd = discogs.DiscogsRequest(m.CursorIndex, &m.FilesInfo[m.CursorIndex])

					return m, cmd

				case STOP_FETCHING:
					return m, cmd
		}
		case tea.KeyMsg:
			switch msg.String() {

				case "up":
					m.Table.MoveUp(1)

					if m.CursorIndex > 0 {
						m.CursorIndex--
					}

					return m, cmd

				case "down":
					m.Table.MoveDown(1)

					if m.CursorIndex < len(m.Files) - 1 {
						m.CursorIndex++
					}

					return m, cmd

				case "shift+up":
					m.Table.MoveUp(10)

					if m.CursorIndex - 10 > 0 {
						m.CursorIndex -= 10
					} else {
						m.CursorIndex = 0
					}
					
					return m, cmd

				case "shift+down":
					m.Table.MoveDown(10)

					if m.CursorIndex + 10 < len(m.Files) {
						m.CursorIndex += 10
					} else {
						m.CursorIndex = len(m.Files) - 1
					}

					return m, cmd

				case " ":
					var rows = m.Table.Rows()

					if(m.IndexSelected == m.CursorIndex) {
						rows[m.CursorIndex][0] = strconv.Itoa(m.IndexSelected)
						m.IndexSelected = -1
					} else {
						if(m.IndexSelected != -1) {
							rows[m.IndexSelected][0] = strconv.Itoa(m.IndexSelected)
						}

						rows[m.CursorIndex][0] = ">>"
						m.IndexSelected = m.CursorIndex
					}

					m.Table.SetRows(rows)
					

					return m, SetFileInfo(m.IndexSelected)

			}
	}

	m.Table, cmd = m.Table.Update(msg)

	return m, cmd
}

func (m MusicTableModel) Resize(windowWidth int, windowHeight int) MusicTableModel {
	var tableWidth = float64((windowWidth / 2) - 4)
	var tableHeight = float64((windowHeight) - 4)

	m.Table.SetWidth(int(tableWidth))
	m.Table.SetHeight(int(tableHeight))

	var indexWidth, extensionWidth, nameWidth, id3Width, statusWidth float64 = 
	tableWidth * (1/12.0), 
	tableWidth * (1/12.0), 
	tableWidth * (6/12.0), 
	tableWidth * (2/12.0), 
	tableWidth * (2/12.0)

	if(indexWidth < 4) {
		indexWidth = 4
	}

	if(extensionWidth < 4) {
		extensionWidth = 4
	}

	if(id3Width < 7) {
		id3Width = 7
	}

	if(statusWidth < 10) {
		statusWidth = 10
	}

	if(indexWidth + extensionWidth + id3Width + statusWidth > tableWidth) {
		nameWidth = tableWidth - (indexWidth - extensionWidth - id3Width - statusWidth)
	}

	m.Table.SetColumns([]table.Column{
		{Title: "Index", Width: int(indexWidth)},
		{Title: "Ext.", Width: int(extensionWidth)},
		{Title: "Name", Width: int(nameWidth)},
		{Title: "Has ID3", Width: int(id3Width)},
		{Title: "Status", Width: int(statusWidth)},
	})

	return m
}

func (m MusicTableModel) View() string {
	return baseStyle.Render(m.Table.View())
}

func InitStyles() table.Styles {
	var styles table.Styles = table.DefaultStyles()

	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	return styles
}

func InitTable(files []ops.FileInfo) MusicTableModel {
	var musicRows []MusicRow
	var rows []table.Row
	var styles table.Styles = InitStyles()

	for index, file := range files {
		var status FileStatus

		var id3Check, complete = CheckID3(file)

		if(file.Extension == "mp3") {
			if(!complete) {
				status = STATUS_PENDING
			}else{
				status = STATUS_FETCH_ACCEPTED
			}
		} else {
			status = STATUS_NOT_MP3
		}

		musicRows = append(musicRows, MusicRow{
			Index: index,
			Extension: file.Extension,
			FileName: file.FileName,
			HasCompleteID3: id3Check,
			Status: status,
		})

		rows = append(rows, table.Row{
			strconv.Itoa(index), file.Extension, file.FileName, id3Check, FileStatusToString(status),
		})
	}

	table := table.New(
		table.WithColumns([]table.Column{
			{Title: "Index", Width: 6},
			{Title: "Ext.", Width: 6},
			{Title: "Name", Width: 40},
			{Title: "Has ID3", Width: 10},
			{Title: "Status", Width: 15},
		}),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(15),

	)
	
	table.SetStyles(styles)

	return MusicTableModel{
		Table: table,
		FilesInfo: files,
		Files: musicRows,
		MusicDir: "./music",
		IndexSelected: 0,
		IndexProcessing: -1,
		CursorIndex: 0,
	}
}