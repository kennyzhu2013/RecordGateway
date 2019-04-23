package record

//字节数组转二进制数组
func Bytes2Bits(data []byte) []byte {
	bin := make([]byte, 0)
	for _, v := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			bin = append(bin, uint8((v>>move)&1))
		}
	}
	return bin
}

//二进制数组转字节数组
func Bits2Bytes(data []byte) []byte {
	bytes := make([]byte, 0)
	for i := 0; i < len(data); i += 8 {
		bytes = append(bytes, (data[i]<<7)+(data[i+1]<<6)+(data[i+2]<<5)+
			(data[i+3]<<4)+(data[i+4]<<3)+(data[i+5]<<2)+(data[i+6]<<1)+(data[i+7]))
	}
	return bytes
}
