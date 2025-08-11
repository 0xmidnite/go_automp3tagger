package discogs

import (
	ops "automp3tagger/file_ops"
	"encoding/json"
	"fmt"
	"io"
	net "net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// DiscogsSearchResult represents a single result from the Discogs search API
type DiscogsSearchResult struct {
	Country        string   `json:"country"`
	Year           string   `json:"year"`
	Format         []string `json:"format"`
	Label          []string `json:"label"`
	Type           string   `json:"type"`
	Genre          []string `json:"genre"`
	Style          []string `json:"style"`
	ID             int      `json:"id"`
	Barcode        []string `json:"barcode"`
	UserData       struct {
		InWantlist    bool `json:"in_wantlist"`
		InCollection  bool `json:"in_collection"`
	} `json:"user_data"`
	MasterID       int      `json:"master_id"`
	MasterURL      string   `json:"master_url"`
	URI            string   `json:"uri"`
	Catno          string   `json:"catno"`
	Title          string   `json:"title"`
	Thumb          string   `json:"thumb"`
	CoverImage     string   `json:"cover_image"`
	ResourceURL    string   `json:"resource_url"`
	Community      struct {
		Want int `json:"want"`
		Have int `json:"have"`
	} `json:"community"`
	FormatQuantity int         `json:"format_quantity"`
	Formats        interface{} `json:"formats"`
}

// DiscogsSearchResponse represents the structure of a Discogs search API response
type DiscogsSearchResponse struct {
	Results []DiscogsSearchResult `json:"results"`
	Pagination struct {
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
		Items   int `json:"items"`
	} `json:"pagination"`
}



func GetRequest(file *ops.FileInfo) (*net.Request, error) {
	var req, err = net.NewRequest("GET", "https://api.discogs.com/database/search?q=" + file.Query + "&per_page=100", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "MP3AutoTagger/0.1 +https://github.com/0xmidnite")
	req.Header.Add("Authorization", "Discogs token=qHYwFaeDMopPKQKLUbnbgsLpepLPhaISPsKmAsKx")

	return req, nil
}

func GetDiscogsResponse(file *ops.FileInfo) (*DiscogsSearchResponse, error) {
	var req, requestErr = GetRequest(file)

	if requestErr != nil {
		return nil, fmt.Errorf("error getting request: %w", requestErr)
	}

	var client = &net.Client{}
	var resp, queryErr = client.Do(req)

	if queryErr != nil {
		return nil, fmt.Errorf("error doing request: %w", queryErr)
	}

	defer resp.Body.Close()

	var body, err = io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var response DiscogsSearchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &response, nil
}

type DiscogsRequestMsg struct {
	Index int
	Response *DiscogsSearchResponse
	Error error
}

func DiscogsRequest(index int, file *ops.FileInfo) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)

		var response, err = GetDiscogsResponse(file)

		return DiscogsRequestMsg{
			Index: index + 1,
			Response: response,
			Error: err,
		}
	}
}
