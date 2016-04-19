package main

import(
	"github.com/daedaluz/agar"
	"github.com/gorilla/websocket"
	"net/http"
	_"math"
	"log"
)


var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


func IsDanger(x *agar.Node, node *agar.Node) bool {
	if node.Radius-5 > x.Radius {
		return true
	}
	return false
}

type Cb struct {
	t *websocket.Conn
	cli *agar.Client
	me uint32
	nodes map[uint32]*agar.Node
}

func (c *Cb) ResetAllCells() {
	c.nodes = make(map[uint32]*agar.Node)
	log.Println("cleared!")
}

func (c *Cb) GetClosestNode(x *agar.Node) *agar.Node {
	var closest *agar.Node
	var sdist float64
	var biggest = int16(0)
	for id, node := range c.nodes {
		if id == x.NodeID {
			continue
		}
		if !node.Virus() {
			dist := x.Distance(node)
			if (node.Radius > biggest) {
			//if node.Radius > sdist{
				closest = node
				biggest = node.Radius
				sdist = dist
				if sdist > 10{
				}
			}
		}
	}
	if closest != nil {
		log.Printf("Closest: %s %dx%d, %d", closest.Name, closest.Position.X, closest.Position.Y, closest.Radius)
	}
	return closest
}

func (c *Cb) OnDisconnect() {
	log.Println("Disconnected!")
	c.t.Close()
	c.cli.Close()
}

func (c *Cb) AddNode(id uint32) {
	c.me = id
}



func (c *Cb) Eaten(id uint32) {

}

func (c *Cb) UpdateNode(x *agar.Node) {
	var me *agar.Node
	if node, exist := c.nodes[x.NodeID]; exist {
		node.Update(x)
		if c.me == node.NodeID {
			me = node
		}
	} else {
		c.nodes[x.NodeID] = x
		if c.me == x.NodeID {
			me = x
		}
	}
	if me == nil {
		me = c.nodes[c.me]
	}
	if me != nil {
		closest := c.GetClosestNode(me)
		if closest != nil {
			posx, posy := me.Tangent(closest, 200)
			c.cli.Move(me.NodeID, posx, posy)
		} else {
			c.cli.Move(me.NodeID, 0, 0)
		}
	}
}

func (c *Cb) RemoveNode(id uint32) {
	delete(c.nodes, id)
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

func main() {
	http.HandleFunc("/", visual_feed)
	http.ListenAndServe(":9090", nil)
}
