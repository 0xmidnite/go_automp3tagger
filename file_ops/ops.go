package file_ops

import (
	"fmt"
	"os"

	id3v2 "github.com/bogem/id3v2/v2"
)

type FileInfo struct {
	Path 		string
	FileName 	string
	Extension 	string
	Query 		string
	Id3Info 	*id3v2.Tag
}

func DeduceFromFilename(filename string) (string, string, string) {
	// var stringSplit = strings.Split(filename, " - ")
	var artist = ""
	var album = ""
	var extra = ""

	return artist, album, extra
}


func GetQuery(filename string) string {
	// var artist, album, catno = deduceFromFilename(filename)
	return ""
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
				Query: GetQuery(fileName),
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
			Query: GetQuery(fileName),
			Id3Info: tag,
		})

		tag.Close()
	}

	return files
}
