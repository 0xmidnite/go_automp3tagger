package ui

import (
	ops "automp3tagger/file_ops"

	discogs "automp3tagger/discogs"

	table "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DetailsSectionModel struct {
	Table 				table.Model
	IsFocused 			bool
	IsUnset 			bool

	File 				*ops.FileInfo
	DiscogsResponses 	[]discogs.DiscogsSearchResponse
	Fetching 			bool
}

func (m DetailsSectionModel) Init() tea.Cmd {
	return nil
}

func (m DetailsSectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return m, nil
		}
	}
	return m, nil
}

func (m DetailsSectionModel) View() string {
	if(m.IsUnset) {
		return ""
	}

	return baseStyleDetails.Render(m.Table.View())
}

func (m DetailsSectionModel) UpdateDiscogsResponses(responses []discogs.DiscogsSearchResponse) DetailsSectionModel {
	m.DiscogsResponses = responses
	return m
}

func (m DetailsSectionModel) UpdateFetching(fetching bool) DetailsSectionModel {
	m.Fetching = fetching
	return m
}

func (m DetailsSectionModel) SetFile(file *ops.FileInfo) DetailsSectionModel {
	if(file == nil) {
		m.File = nil
		m.IsUnset = true

		m.Table.SetColumns([]table.Column{
			{Title: "No file selected", Width: 50},
		})
	
		m.Table.SetRows([]table.Row{
			{"NIL"},
		})

		return m
	}

	if(file.Extension != "mp3") {
		m.File = nil
		m.IsUnset = false

		m.Table.SetColumns([]table.Column{
			{Title: "Not a mp3 file...", Width: 50},
		})
	
		m.Table.SetRows([]table.Row{
			{"Not a mp3 file... No Info to show"},
		})

		return m
	}

	m.IsUnset = false
	m.File = file

	m.Table.SetColumns([]table.Column{
		{Title: file.FileName, Width: 50},
	})

	m.Table.SetRows([]table.Row{
		{"Path: " + file.Path},
		{"Name: " + file.FileName},
		{"Extension: " + file.Extension},
		{"Artist: " + file.Id3Info.Artist()},
		{"Album: " + file.Id3Info.Album()},
		{"Title: " + file.Id3Info.Title()},
		{"Genre: " + file.Id3Info.Genre()},
		{"Year: " + file.Id3Info.Year()},
	})

	return m
}

func (m DetailsSectionModel) Resize(windowWidth int, windowHeight int) DetailsSectionModel {
	var tableWidth float32 = float32((windowWidth / 2.0) - 10.0)
	var tableHeight = (windowHeight) - 4

	m.Table.SetWidth(int(tableWidth))
	m.Table.SetHeight(tableHeight)

	m.Table.SetColumns([]table.Column{
		{Title: "No file selected", Width: int(tableWidth)},
	})

	return m
}

var baseStyleDetails = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))


func InitInfoStyles() table.Styles {
	var styles table.Styles = table.DefaultStyles()

	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	return styles
}

func InitInfoTable() DetailsSectionModel {
	var styles = InitInfoStyles()
	table := table.New(
		table.WithColumns([]table.Column{
			{Title: "No file selected", Width: 50},
		}),
		table.WithRows([]table.Row{
			{"NIL"},
		}),
		table.WithFocused(false),
		table.WithHeight(15),

	)

	table.SetStyles(styles)

	return DetailsSectionModel{
		Table: table,
		IsFocused: false,
		IsUnset: true,
		File: nil,
	}
}


