package boilang

import (
	"bytes"
)

type BoiBuf struct {
	buf *bytes.Buffer
}

func (BoiBuf) log(x ...interface{}) {
	x = append([]interface{}{"[buf]"}, x...)
	//fmt.Println(x...)
}

func (bb *BoiBuf) u8(x uint8) {
	bb.log("u8", x)
	bb.buf.WriteByte(x)
}

func (bb *BoiBuf) u16(x uint16) {
	bb.log("u16", x)
	bb.u8(uint8(x))
	bb.u8(uint8(x >> 8))
}

func (bb *BoiBuf) ru8() uint8 {
	x, err := bb.buf.ReadByte()
	if err != nil {
		panic(err)
	}
	bb.log("ru8", x)
	return x
}

func (bb *BoiBuf) ru16() uint16 {
	x := uint16(bb.ru8())
	x |= uint16(bb.ru8()) << 8
	bb.log("ru16", x)
	return x
}
