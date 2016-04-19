package main

import(
	"net/http"
	"github.com/daedaluz/agar"
	"github.com/gorilla/websocket"
	"log"
)



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func visual_feed(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("unable to upgrade:", err)
		return
	}
	if srv, err := agar.FindServer(); err != nil {
		log.Println("Unable to find server", err)
		return
	} else {
		cbs := &Cb{
			t: conn,
			nodes: make(map[uint32]*agar.Node),
		}

		if cli, err := agar.NewClient(srv, cbs); err != nil {
			log.Println("couldn't establish connection to agar server", err)
		} else {
			cbs.cli = cli
			cbs.cli.Tee(conn)
			cli.SetName("companionball<3")
			for _, data, err := conn.ReadMessage(); err == nil; _, data, err = conn.ReadMessage() {
				if data[0] == 0 {
					cli.Spawn()
				}
			}
		}
	}
}

func init() {
	http.HandleFunc("/", visual_feed)
}
