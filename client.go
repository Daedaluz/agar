package agar

import (
	"log"
	"encoding/binary"
	"bytes"
	"github.com/gorilla/websocket"
	"github.com/daedaluz/agar/intern"
	"net/http"
	"sync"
)

type Ondc interface {
	OnDisconnect()
}

type Client struct {
	sync.Mutex
	c *websocket.Conn
	tee []*websocket.Conn
	name string
	handlers map[uint8]handler
	cfg interface{}
}

func NewClient(srv *Server, handlers interface{}) (*Client, error) {
	headers := http.Header{
		"Origin": []string{"http://agar.io"},
	}
	sock,_, e := websocket.DefaultDialer.Dial("ws://" + srv.Ip + "/", headers)
	if e != nil {
		return nil, e
	}

	if x, err := intern.ResetConnection1(); err != nil {
		sock.Close()
		return nil, err
	} else {
		if err := sock.WriteMessage(websocket.BinaryMessage, x); err != nil {
			sock.Close()
			return nil, err
		}
	}

	if x, err := intern.ResetConnection2(); err != nil {
		sock.Close()
		return nil, err
	} else {
		if err := sock.WriteMessage(websocket.BinaryMessage, x); err != nil {
			sock.Close()
			return nil, err
		}
	}

	if x, err := intern.SendToken([]byte(srv.Token)); err != nil {
		sock.Close()
		return nil, err
	} else {
		if err := sock.WriteMessage(websocket.BinaryMessage, x); err != nil {
			sock.Close()
			return nil, err
		}
	}

	cli := &Client {
		c: sock,
		handlers: stdHandlers,
		cfg: handlers,
		tee: make([]*websocket.Conn, 0, 10),
	}
	go cli.reader()
	return cli, nil
}

func (c *Client) reader() {
	for {
		t, data, err := c.c.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		for i, tee := range c.tee {
			if tee != nil {
				if err := tee.WriteMessage(t, data); err != nil {
					c.tee[i] = nil
				}
			}
		}
		reader := bytes.NewBuffer(data)
		var packetID uint8
		binary.Read(reader, binary.LittleEndian, &packetID)
		if decoder, exist := c.handlers[packetID]; exist {
			decoder(c, reader)
		} else {
			log.Printf("Undhandled packetID: %d", packetID)
		}
	}
	if x, ok := c.cfg.(Ondc); ok {
		x.OnDisconnect()
	}
}

func (c *Client) Tee(dst *websocket.Conn) {
	c.tee = append(c.tee, dst)
}

func (c *Client) SetName(name string) {
	c.name = name
}

func (c *Client) Spawn() {
	if c.name == "" {
		c.name = "Assguard"
	}
	data, _ := intern.SetNickname(c.name)
	c.Lock()
	c.c.WriteMessage(websocket.BinaryMessage, data)
	c.Unlock()
}

func (c *Client) Spectate() {
	data, _ := intern.Spectate()
	c.Lock()
	c.c.WriteMessage(websocket.BinaryMessage, data)
	c.Unlock()
}

func (c *Client) Move(nodeid uint32, x, y int32) {
	data, _ := intern.Move(nodeid, x, y)
	c.Lock()
	c.c.WriteMessage(websocket.BinaryMessage, data)
	c.Unlock()
}

func (c *Client) Split() {
	data, _ := intern.Split()
	c.Lock()
	c.c.WriteMessage(websocket.BinaryMessage, data)
	c.Unlock()
}

func (c *Client) Explode() {
	data, _ := intern.Explode()
	c.Lock()
	c.c.WriteMessage(websocket.BinaryMessage, data)
	c.Unlock()
}

func (c *Client) Eject() {
	data, _ := intern.Eject()
	c.Lock()
	c.c.WriteMessage(websocket.BinaryMessage, data)
	c.Unlock()
}

func (c *Client) Close() {
	c.c.Close()
}


