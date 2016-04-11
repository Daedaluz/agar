package intern

import(
	"bytes"
	"encoding/binary"
)

func bwrite(buf *bytes.Buffer, towrite ...interface{}) error {
	for _, v := range towrite {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	return nil
}



func ResetConnection1() ([]byte, error) {
	buf := new(bytes.Buffer)
	packetID := uint8(254)
	if ok := bwrite(buf, packetID, uint32(5)); ok != nil {
		return nil, ok
	}
	return buf.Bytes(), nil
}

func ResetConnection2() ([]byte, error) {
	buf := new(bytes.Buffer)
	packetID := uint8(255)
	if ok := bwrite(buf, packetID, uint32(154669603)); ok != nil {
		return nil, ok
	}
	return buf.Bytes(), nil
}

func SendToken(token []byte) ([]byte, error) {
	buff := make([]byte,0, 10)
	buff = append(buff, uint8(80))
	buff = append(buff, token...)
	return buff, nil
}

func SetNickname(name string) ([]byte, error){
	buff := new(bytes.Buffer)
	buff.Write([]byte{0x00})
	WriteString(buff, name)
	return buff.Bytes(), nil
}

func Spectate() ([]byte, error){
	return []byte{1}, nil
}

func Move(nid uint32, x, y int32) ([]byte, error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, uint8(16))
	binary.Write(buff, binary.LittleEndian, x)
	binary.Write(buff, binary.LittleEndian, y)
	binary.Write(buff, binary.LittleEndian, nid)
	return buff.Bytes(), nil
}

func Split() ([]byte, error) {
	return []byte{17}, nil
}

func Explode() ([]byte, error) {
	return []byte{20}, nil
}

func Eject() ([]byte, error) {
	return []byte{21}, nil
}

