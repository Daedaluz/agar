package agar

import(
	"encoding/json"
	"strings"
	"strconv"
	"net/http"
)

const (
	AGAR_VERSION = uint32(154669603)
	LOCATION = "EU-London"
)

type Server struct {
	Ip string `json:"ip"`
	Token string `json:"token"`
}

func FindServer() (*Server, error) {
	data := strings.NewReader(LOCATION + "\n" + strconv.Itoa(int(AGAR_VERSION)))
	res, e := http.Post("http://m.agar.io/findServer", "plain/text", data)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	result := Server{}
	dec := json.NewDecoder(res.Body)
	e = dec.Decode(&result)
	if e != nil {
		return nil, e
	}
	return &result, nil
}

