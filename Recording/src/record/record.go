package record

type RecordFile interface {
	Close()
	Write(b []byte) (n int, err error)
}

func CreateRecordFile(filepath string, media string) RecordFile {
	if media == "711" {
		return CreateFile711(filepath)
	}else if media == "amrnb-ocet"{
		return CreateFileAmrNB(filepath)
	}
}