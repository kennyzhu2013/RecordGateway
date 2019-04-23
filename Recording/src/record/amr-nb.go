package record

import (
	"os"
)

type FileAmrNB struct {
	file *os.File
	mode bool
}

func CreateFileAmrNB(file string, mode bool) *FileAmrNB {
	// mode:true 八位对齐  mode: false 非八位对齐
	f, _ := os.Create(file)
	//写入amrnb文件头
	f.Write([]byte{'#', '!', 'A', 'M', 'R', '\n'})
	return &FileAmrNB{file: f, mode: mode}
}

func (f *FileAmrNB) Close() {
	f.file.Close()
}

//amr-nb帧长度
var amr_nbFrameLen = [9]int{95, 103, 118, 134, 148, 159, 204, 244, 39}

func (f *FileAmrNB) Write(b []byte) (n int, err error) {
	b = b[12:]
	b = Bytes2Bits(b)
	res := make([]byte, 0)
	heads := make([]byte, 0)
	i := 0
	if f.mode == false {
		i = 4
		for {
			heads = append(heads, 0)
			heads = append(heads, b[i+1:i+6]...)
			heads = append(heads, 0, 0)
			if b[i] == 0 {
				i += 6
				break
			}
			i += 6
		}
	} else if f.mode == true {
		i = 8
		for {
			heads = append(heads, 0)
			heads = append(heads, b[i+1:i+6]...)
			heads = append(heads, 0, 0)
			if b[i] == 0 {
				i += 8
				break
			}
			i += 8
		}
	}
	flen := 0
	for j := 0; j < len(heads); j += 8 {
		res = append(res, heads[j:j+8]...)
		flen = amr_nbFrameLen[int((heads[j+1]<<3)+(heads[j+2]<<2)+(heads[j+3]<<1)+heads[j+4])]
		res = append(res, b[i:i+flen]...)
		for num := 8 - flen%8; num > 0; num-- {
			res = append(res, 0)
		}
		if f.mode == false {
			i += flen
		} else if f.mode == true {
			i += flen + 8 - flen%8
		}
	}
	res = Bits2Bytes(res)
	return f.file.Write(res)
}
