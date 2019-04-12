package record

import (
	"os"
)

type File711 struct {
	file *os.File
}

func CreateFile711(file string) *File711{
	f, _ := os.Create(file)
	return &File711{file: f}
}

func (f *File711) Close(){
	f.file.Close()
}

func (f *File711) Write(b []byte) (n int, err error) {
	return f.file.Write(b)
}