package intern

import(
	"io"
	_"log"
	"bytes"
	"encoding/binary"
	"unicode/utf16"
	"unicode/utf8"
)

func ReadString(i io.Reader) string {
	tmp := uint16(0)
	//str := make([]uint16, 0, 20)
	res := make([]byte, 0, 1024)
	p := make([]byte, 4)
	for {
		x := binary.Read(i, binary.LittleEndian, &tmp)
		if x != nil {
			return ""
		}
		if tmp == 0 {
			return string(res)
		}
		r := utf16.Decode([]uint16{tmp})
		utf8.EncodeRune(p, r[0])
		res = append(res, p...)
	}
}

func WriteString(o io.Writer, text string) {
	buf := new(bytes.Buffer)
	codes := []rune(text)
	toprint := utf16.Encode(codes)
	toprint = append(toprint, uint16(0))
	binary.Write(buf, binary.LittleEndian, toprint)
	o.Write(buf.Bytes())
}

//func WriteString(o io.Writer, text string) {
//	buf := new(bytes.Buffer)
//	runes := make([]rune, 0, 10)
//	b := []byte(text)
//
//	for len(b) > 0 {
//		r, size := utf8.DecodeRune(b)
//		runes = append(runes, r)
//		b = b[size:]
//	}
//
//	encoded := utf16.Encode(runes)
//	encoded = append(encoded, uint16(0))
//
//	binary.Write(buf, binary.LittleEndian, encoded)
//
//	o.Write(buf.Bytes())
//}

