package record

import (
	"bytes"
	"encoding/binary"
	"os"
)

type File711 struct {
	file *os.File
	buf  []byte
	mode string //pcma  pcmu
}

func CreateFile711(file string, mode string) *File711 {
	f, _ := os.Create(file)
	return &File711{file: f, buf: make([]byte, 0), mode: mode}
}

func (f *File711) Close() {
	datalen := int2hexbyte(int32(len(f.buf)))
	rifflen := int2hexbyte(int32(len(f.buf) + 50))
	header := []byte{'R', 'I', 'F', 'F'}
	header = append(header, rifflen...)
	if f.mode == "pcma" {
		header = append(header, 'W', 'A', 'V', 'E', 'f', 'm', 't', ' ', 18, 0, 0, 0, 6, 0, 1, 0, 64, 31, 0, 0, 64, 31, 0, 0, 1, 0, 8, 0, 0, 0, 'f', 'a', 'c', 't', 4, 0, 0, 0)
	} else if f.mode == "pcmu" {
		header = append(header, 'W', 'A', 'V', 'E', 'f', 'm', 't', ' ', 18, 0, 0, 0, 7, 0, 1, 0, 64, 31, 0, 0, 64, 31, 0, 0, 1, 0, 8, 0, 0, 0, 'f', 'a', 'c', 't', 4, 0, 0, 0)
	}
	header = append(header, datalen...)
	header = append(header, 'd', 'a', 't', 'a')
	header = append(header, datalen...)
	f.buf = append(header, f.buf...)
	f.file.Write(f.buf)
	println(len(f.buf))
	f.file.Close()
}

func (f *File711) Write(b []byte) (n int, err error) {
	f.buf = append(f.buf, b[12:]...)
	return len(b), nil
}

//int转byte数组，在高地位反转，用于wav头部的riff长度，data长度和fact采样数
func int2hexbyte(l int32) []byte {
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, l)
	res := buffer.Bytes()
	res[0], res[1], res[2], res[3] = res[3], res[2], res[1], res[0]
	return res
}
