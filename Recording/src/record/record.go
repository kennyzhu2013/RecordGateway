package record

import (
	"fmt"
)

type RecordFile interface {
	Close()
	Write(b []byte) (n int, err error)
}

func CreateRecordFile(filepath string, media string) RecordFile {
	//amrnb-0 amrnb-1 amrwb-0 amrwb-1 pcma(8) pcmu(0)
	if media == "pcma" {
		return CreateFile711(filepath, "pcma")
	} else if media == "pcmu" {
		return CreateFile711(filepath, "pcmu")
	} else if media == "amrnb-0" {
		return CreateFileAmrNB(filepath, false)
	} else if media == "amrnb-1" {
		return CreateFileAmrNB(filepath, true)
	} else if media == "amrwb-0" {
		return CreateFileAmrWB(filepath, false)
	} else if media == "amrwb-1" {
		return CreateFileAmrWB(filepath, true)
	} else {
		panic(fmt.Sprintf("unkonwn codec!! %s", media))
	}
}
