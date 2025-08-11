package ui

import (
	"strconv"

	ops "automp3tagger/file_ops"

	id3v2 "github.com/bogem/id3v2/v2"
)

type FileStatus int

const (
	STATUS_PENDING 	FileStatus = iota
	STATUS_FETCH_OK
	STATUS_FETCH_ERROR 
	STATUS_FETCH_ACCEPTED
	STATUS_FETCH_REJECTED
	STATUS_FETCHING
	STATUS_NOT_MP3
)


func FileStatusToString(status FileStatus) string {
	switch status {
		case STATUS_FETCH_OK:
			return "☑️\t"
		case STATUS_FETCH_ACCEPTED:
			return "✅\t"
		case STATUS_FETCH_REJECTED:
			return "❌\t"
		case STATUS_PENDING:
			return "💬\t"
		case STATUS_FETCH_ERROR:
			return "‼️\t"
		case STATUS_FETCHING:
			return "🔎\t"
		default:
			return "⛔️\t"
	}
}


func CheckID3(file ops.FileInfo) (string, bool) {
	var complete bool = true

	if(file.Id3Info != nil) {
		var id3TagFlag = ""
		var img = file.Id3Info.GetFrames(id3v2.V23CommonIDs["Attached picture"])

		if(file.Id3Info.Artist() != ""){
			id3TagFlag += "A"
		}else{
			id3TagFlag += "-"
			complete = false
		}

		if(file.Id3Info.Album() != ""){
			id3TagFlag += "a"
		}else{
			id3TagFlag += "-"
			complete = false
		}

		if(file.Id3Info.Title() != ""){
			id3TagFlag += "T"
		}else{
			id3TagFlag += "-"
			complete = false
		}

		if(file.Id3Info.Genre() != ""){
			id3TagFlag += "G"
		}else{
			id3TagFlag += "-"
			complete = false
		}
		if(file.Id3Info.Year() != ""){
			id3TagFlag += "Y"
		}else{
			id3TagFlag += "-"
			complete = false
		}

		if(len(img) > 0){
			id3TagFlag += "I" + strconv.Itoa(len(img))
		}else{
			id3TagFlag += "i-"
			complete = false
		}

		if(len(id3TagFlag) == 0){
			return "Empty", false
		}

		return id3TagFlag, complete
	}
	
	return "No Tag", false
}

