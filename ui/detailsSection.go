package ui

import (
	ops "automp3tagger/file_ops"
	"math"
	"slices"
	"strconv"
	"strings"

	discogs "automp3tagger/discogs"

	table "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DetailsSectionModel struct {
	Mp3InfoTable 		table.Model
	DiscogsTable 		table.Model
	
	IsFocused 			bool
	IsUnset 			bool
	ColumnWidth 	    int

	File 				*ops.FileInfo
	DiscogsResponses 	[]discogs.DiscogsSearchResult
	Fetching 			bool
	IsInResDetails 		bool
	LastCursorIndex 	int
	// CursorIndex 		int
}

var focusedStyleDetails = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderForeground(lipgloss.Color("#FFE2AA"))

var baseStyleDetails = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func (m DetailsSectionModel) GetFileRows() []table.Row {
	return []table.Row{
		{"Path: " + m.File.Path},
		{"Name: " + m.File.FileName},
		{"Extension: " + m.File.Extension},
		{"Artist: " + m.File.Id3Info.Artist()},
		{"Album: " + m.File.Id3Info.Album()},
		{"Title: " + m.File.Id3Info.Title()},
		{"Genre: " + m.File.Id3Info.Genre()},
		{"Year: " + m.File.Id3Info.Year()},
		{"Query: " + m.File.Query},
	}
}

func (m DetailsSectionModel) GetDiscogsResponseRows() []table.Row {
	if(len(m.DiscogsResponses) == 0) {
		return []table.Row{}
	}

	var rows = []table.Row{}

	for _, response := range m.DiscogsResponses {
		rows = append(rows, table.Row{
			response.Title + " - " + response.Catno,
		})
	}

	return rows
}

func (m DetailsSectionModel) UpdateDiscogsResponses(responses []discogs.DiscogsSearchResult) DetailsSectionModel {
	if(m.IsUnset) {
		return m
	}
	
	var rows = []table.Row{}

	if(len(responses) == 0) {
		rows = slices.Concat(rows, []table.Row{
			{"No results found"},
		})
	}else{
		m.DiscogsResponses = responses
		m.DiscogsTable.SetColumns([]table.Column{
			{Title: "Responses | " + strconv.Itoa(len(m.DiscogsResponses)) + " results found", Width: m.ColumnWidth},
		})
		rows = slices.Concat(rows, m.GetDiscogsResponseRows())
	}

	m.DiscogsTable.SetRows(rows)

	return m
}

func (m DetailsSectionModel) UpdateFetching(fetching bool) DetailsSectionModel {
	m.Fetching = fetching
	var isFetching = "Idle"

	if(fetching) {
		isFetching = "Fetching..."
	}

	if(!m.IsUnset){
		m.Mp3InfoTable.SetColumns([]table.Column{
			{Title:  m.File.FileName + " - " + isFetching, Width: m.ColumnWidth},
		})
	}
	return m
}

func (m DetailsSectionModel) SetFile(file *ops.FileInfo, discogsResponses *[]discogs.DiscogsSearchResult) DetailsSectionModel {	
	if(file == nil) {
		m.File = nil
		m.IsUnset = true

		m.Mp3InfoTable.SetColumns([]table.Column{
			{Title: "No file selected", Width: m.ColumnWidth},
		})
	
		m.Mp3InfoTable.SetRows([]table.Row{
			{"NIL"},
		})

		return m
	}

	if(file.Extension != "mp3") {
		m.File = nil
		m.IsUnset = false

		m.Mp3InfoTable.SetColumns([]table.Column{
			{Title: "Not a mp3 file...", Width: m.ColumnWidth},
		})
	
		m.Mp3InfoTable.SetRows([]table.Row{
			{"Not a mp3 file... No Info to show"},
		})

		return m
	}

	m.IsUnset = false
	m.File = file
	var isFetching = "Idle"

	if(m.Fetching) {
		isFetching = "Fetching..."
	}

	m.Mp3InfoTable.SetColumns([]table.Column{
		{Title: file.FileName + " - " + isFetching, Width: m.ColumnWidth},
	})

	var rows = m.GetFileRows()
	var discogsRows = []table.Row{}
	
	m.DiscogsResponses = *discogsResponses


	if(len(m.DiscogsResponses) == 0) {
		m.DiscogsTable.SetColumns([]table.Column{
			{Title: "Responses | No results found", Width: m.ColumnWidth},
		})
	}else{
		m.DiscogsTable.SetColumns([]table.Column{
			{Title: "Responses | " + strconv.Itoa(len(m.DiscogsResponses)) + " results found", Width: m.ColumnWidth},
		})
		discogsRows = slices.Concat(discogsRows, m.GetDiscogsResponseRows())
	}

	m.Mp3InfoTable.SetRows(rows)
	m.DiscogsTable.SetRows(discogsRows)

	return m
}

func (m DetailsSectionModel) Resize(windowWidth int, windowHeight int) DetailsSectionModel {
	var tableWidth float64 = math.Round(float64((windowWidth / 2.0) - 10.0))
	var tableHeight float64 = math.Round(float64(((windowHeight) - 4) / 2))

	m.Mp3InfoTable.SetWidth(int(tableWidth))
	m.Mp3InfoTable.SetHeight(int(tableHeight))

	m.DiscogsTable.SetWidth(int(tableWidth))
	m.DiscogsTable.SetHeight(int(math.Ceil(tableHeight)) - 2)

	var columns = m.Mp3InfoTable.Columns()
	for i := 0; i < len(columns); i++ {
		columns[i].Width = int(tableWidth)
	}
	
	var discogsColumns = m.DiscogsTable.Columns()
	for i := 0; i < len(discogsColumns); i++ {
		discogsColumns[i].Width = int(tableWidth)
	}

	m.ColumnWidth = int(tableWidth)
	m.Mp3InfoTable.SetColumns(columns)
	m.DiscogsTable.SetColumns(discogsColumns)

	return m
}


func InitInfoStyles() table.Styles {
	var styles table.Styles = table.DefaultStyles()

	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	styles.Selected = styles.Selected.
		Foreground(lipgloss.Color("229")).
		Bold(true)


	return styles
}

func InitInfoTable() DetailsSectionModel {
	var styles = InitInfoStyles()
	tableMp3Info := table.New(
		table.WithColumns([]table.Column{
			{Title: "No file selected", Width: 70},
		}),
		table.WithRows([]table.Row{
			{"NIL"},
		}),
		table.WithFocused(false),
		table.WithHeight(15),

	)

	tableMp3Info.SetStyles(styles)

	tableDiscogs := table.New(
		table.WithColumns([]table.Column{
			{Title: "Responses", Width: 70},
		}),
		table.WithRows([]table.Row{
			{"No results."},
		}),
		table.WithFocused(false),
		table.WithHeight(15),
	)

	tableDiscogs.SetStyles(styles)

	return DetailsSectionModel{
		Mp3InfoTable: tableMp3Info,
		DiscogsTable: tableDiscogs,
		IsFocused: false,
		IsUnset: true,
		File: nil,
		DiscogsResponses: []discogs.DiscogsSearchResult{},
		Fetching: false,
		ColumnWidth: 70,
	}
}

func (m DetailsSectionModel) Init() tea.Cmd {
	return nil
}

func (m DetailsSectionModel) ResetState() DetailsSectionModel {
	m.IsInResDetails = false
	m.LastCursorIndex = -1
	return m
}

func (m DetailsSectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return m, nil

		case "up":
			m.DiscogsTable.MoveUp(1)
			
			return m, cmd

		case "down":
			m.DiscogsTable.MoveDown(1)

			return m, cmd

		case "shift+up":
			m.DiscogsTable.MoveUp(10)
			
			return m, cmd

		case "shift+down":
			m.DiscogsTable.MoveDown(10)

			return m, cmd

		case " ":
			if(len(m.DiscogsResponses) == 0 || m.IsInResDetails) {
				return m, cmd
			}

			m.IsInResDetails = true
			
			var cursorPosition = m.DiscogsTable.Cursor()
			var response = m.DiscogsResponses[cursorPosition]

			m.LastCursorIndex = cursorPosition

			var newRows = []table.Row{
				{"Title: " + response.Title},
				{"Catno: " + response.Catno},
				{"Format: " + strings.Join(response.Format, ", ")},
				{"Label: " + strings.Join(response.Label, ", ")},
				{"Year: " + response.Year},
				{"Genre: " + strings.Join(response.Genre, ", ")},
				{"Style: " + strings.Join(response.Style, ", ")},
				{"Country: " + response.Country},
				{"Thumb: " + response.Thumb},
				{"ResourceURL: " + response.ResourceURL},
				{"URI: " + response.URI},
			}

			m.DiscogsTable.SetRows(newRows)
			m.DiscogsTable.SetColumns([]table.Column{
				{Title: "Response Details for " + response.Title, Width: m.ColumnWidth},
			})
			m.DiscogsTable.SetCursor(0)

			return m, cmd


			
		case "backspace":
			if(m.IsInResDetails) {
				m.IsInResDetails = false
				m.DiscogsTable.SetColumns([]table.Column{
					{Title: "Responses | " + strconv.Itoa(len(m.DiscogsResponses)) + " results found", Width: m.ColumnWidth},
				})
				m.DiscogsTable.SetRows(m.GetDiscogsResponseRows())
				m.DiscogsTable.SetCursor(m.LastCursorIndex)

				m.LastCursorIndex = -1
			}
			return m, cmd
		}
	}
	return m, nil
}

func (m DetailsSectionModel) View() string {
	if(m.IsUnset) {
		return ""
	}

	if(m.IsFocused) {
		return lipgloss.JoinVertical(lipgloss.Left, baseStyleDetails.Render(m.Mp3InfoTable.View()), focusedStyleDetails.Render(m.DiscogsTable.View()))
	}

	return lipgloss.JoinVertical(lipgloss.Left, baseStyleDetails.Render(m.Mp3InfoTable.View()), baseStyleDetails.Render(m.DiscogsTable.View()))
}

