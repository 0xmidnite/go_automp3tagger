package main

import (
	ops "automp3tagger/file_ops"
	ui "automp3tagger/ui"

	"fmt"
	"os"

	discogs "automp3tagger/discogs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainModel struct {
	files 			[]ops.FileInfo
	musicTable 		ui.MusicTableModel
	detailsSection 	ui.DetailsSectionModel
	actionSection 	ui.ActionSectionModel
}

func (m mainModel) Init() tea.Cmd {
	return nil
}



func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

    switch msg := msg.(type) {
		case discogs.DiscogsRequestMsg:
			var rows = m.musicTable.Table.Rows()
			var index = msg.Index - 1

			if(msg.Error != nil) {

				rows[index][4] = ui.FileStatusToString(ui.STATUS_FETCH_ERROR)

				m.musicTable.Table.SetRows(rows)

				return m, cmd
			}
			

			rows[index][4] = ui.FileStatusToString(ui.STATUS_FETCH_OK)

			m.musicTable.Table.SetRows(rows)

			return m, cmd

		case ui.SetFileMsg:
			if(msg.SelectedIndex == -1) {
				m.detailsSection = m.detailsSection.SetFile(nil)
			}else{
				m.detailsSection = m.detailsSection.SetFile(&m.files[msg.SelectedIndex])
			}

			return m, cmd

		case ui.ActionMsg:
			var updatedModel tea.Model
			updatedModel, cmd = m.musicTable.Update(msg)
			m.musicTable = updatedModel.(ui.MusicTableModel)
			
			return m, cmd

		case tea.WindowSizeMsg:
			m.musicTable = m.musicTable.Resize(msg.Width, msg.Height)
			m.detailsSection = m.detailsSection.Resize(msg.Width, msg.Height)

			return m, cmd


    	case tea.KeyMsg:
    		switch msg.String() {
				case "tab":

					return m, nil

				case "right", "left", "enter":
					var updatedModel tea.Model

					updatedModel, cmd = m.actionSection.Update(msg)
					m.actionSection = updatedModel.(ui.ActionSectionModel)

					return m, cmd


				case "up", "down", "shift+up", "shift+down", " ":
					if(m.musicTable.Table.Focused()) {
						var updatedModel tea.Model

						updatedModel, cmd = m.musicTable.Update(msg)
						m.musicTable = updatedModel.(ui.MusicTableModel)
					}

					return m, cmd

        		case "ctrl+c", "q":
        	    	return m, tea.Quit

    	}
	}

    return m, cmd
}

func (m mainModel) View() string { 
	if(!m.detailsSection.IsUnset) {
		var tablesView = lipgloss.JoinHorizontal(lipgloss.Left, m.musicTable.View(), m.detailsSection.View())
		
		return lipgloss.JoinVertical(lipgloss.Left, tablesView,  m.actionSection.View())
	} else {
		return lipgloss.JoinVertical(lipgloss.Left, m.musicTable.View(), m.actionSection.View())
	}
}

func main() {
	var files []ops.FileInfo = ops.PrepareFiles()

	// fmt.Printf(
	// 	"length: %d\nAt index 2: %s\nExtension: %s\nPath: %s\nTitle: %s\nArtist: %s\nAlbum: %s\n",
	// 	 len(files), files[2].FileName, files[2].Extension, files[2].Path, files[2].Id3Info.Title(), files[2].Id3Info.Artist(), files[2].Id3Info.Album(),
	// )

	musicTable := ui.InitTable(files)
	actionSection := ui.InitActionSection()
	detailsSection := ui.InitInfoTable();

	p := tea.NewProgram(mainModel{musicTable: musicTable, files: files, actionSection: actionSection, detailsSection: detailsSection})

	if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }

	os.Exit(0)
}