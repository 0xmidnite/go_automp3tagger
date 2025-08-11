package file_ops

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	id3v2 "github.com/bogem/id3v2/v2"
)

type FileInfo struct {
	Path 		string
	FileName 	string
	Extension 	string
	Query 		string
	Id3Info 	*id3v2.Tag
}

func SanitizeFilename(filename string, removeDashes bool) string {
    // Pattern to match numbers and periods at the beginning
    patternNumbersAndPeriodsFirst := regexp.MustCompile(`^[0-9.]*[0-9][0-9.]*`)
    
    // Pattern for special characters (with or without dashes)
    var patternSpecialCharacters *regexp.Regexp
    if removeDashes {
        patternSpecialCharacters = regexp.MustCompile(`[-_<>]`)
    } else {
        patternSpecialCharacters = regexp.MustCompile(`[_<>]`)
    }
    
    // Get the filename without extension
    base := filepath.Base(filename)
    name := strings.TrimSuffix(base, filepath.Ext(base))
    
    // Remove numbers and periods from the beginning
    name = patternNumbersAndPeriodsFirst.ReplaceAllString(name, "")
    
    // Replace special characters with spaces
    name = patternSpecialCharacters.ReplaceAllString(name, " ")
    
    // Trim whitespace
    return strings.TrimSpace(name)
}

func DeduceFromFilename(filename string) (string, string, string, error) {
	var sanitized = SanitizeFilename(filename, false)
	var split = strings.Split(sanitized, " - ")

	if(len(split) == 1) {
		return "", "", split[0], nil
	}

	if(len(split) == 2) {
		return "", split[0], split[1], nil
	}

	if(len(split) == 3) {
		return split[0], split[1], split[2], nil
	}

	return "", "", "", errors.New("found nothing")
}


func GetQuery(filename string, tag *id3v2.Tag) string {
	var catno, artist, album, err = DeduceFromFilename(filename)

	if(tag != nil) {
		var artistTag = tag.Artist()

		if(artistTag != "") {
			artist = artistTag
		}

		var albumTag = tag.Title()

		if(albumTag != "") {
			album = albumTag
		}

		return artist + " " + album 
	}

	if(err != nil) {
		fmt.Println("Error deducing from filename:", err)
		return ""
	}

	return catno + " " + artist + " " + album
}

func PrepareFiles() []FileInfo {
	var files []FileInfo
	var filePath string = "./music/"
	entries, err := os.ReadDir("./music")

	if err != nil {
		fmt.Println("Error reading directory:", err)
		return []FileInfo{}
	}

	for _, entry := range entries {
		var pathName = filePath + entry.Name()
		var extension = entry.Name()[len(entry.Name())-3:]
		var fileName = entry.Name()[:len(entry.Name())-4]

		if(extension != "mp3") {
			files = append(files, FileInfo{
				Path: pathName,
				FileName: fileName,
				Extension: extension,
				Query: GetQuery(fileName, nil),
				Id3Info: nil,
			})
			continue
		}

		var tag, err = id3v2.Open(pathName, id3v2.Options{Parse: true})

		if err != nil {
			fmt.Printf("Error opening file %s: %s\n", pathName, err)
			continue
		}

		files = append(files, FileInfo{
			Path: pathName,
			FileName: fileName,
			Extension: extension,
			Query: GetQuery(fileName, tag),
			Id3Info: tag,
		})

		tag.Close()
	}

	return files
}
