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
	for _, x := range toprint {
		binary.Write(buf, binary.LittleEndian, x)
	}
	o.Write(buf.Bytes())
}

