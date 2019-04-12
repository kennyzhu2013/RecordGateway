package record

import (
	"os"
)

type FileAmrNB struct {
	file *os.File
}

func CreateFileAmrNB(file string) *File711{
	f, _ := os.Create(file)
	return &File711{file: f}
}

func (f *FileAmrNB) Close(){
	f.file.Close()
}

func (f *FileAmrNB) Write(b []byte) (n int, err error) {
	return f.file.Write(b)
}