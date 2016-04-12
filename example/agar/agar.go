package main

import(
	"log"
	_"time"
	"github.com/daedaluz/agar"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const (
	WIDTH = 1280
	HEIGHT = 720
)
var client *agar.Client

var lock sync.Mutex
var mynodes = map[uint32]bool{}

var nodes = map[uint32]*agar.Node{}
var names = map[uint32]string{}

var CamNode uint32

func pickCamNode() {
	for id, _ := range mynodes {
		originX = nodes[id].Position.X
		originY = nodes[id].Position.Y
		CamNode = id
		return
	}
}

func isMe(id uint32) bool {
	if me, ok := mynodes[id]; ok && me {
		return true
	}
	return false
}

var border agar.Border

var originX, originY int32

func getNodeRect(n *agar.Node) *sdl.Rect {

	x := n.Position.X - originX
	y := n.Position.Y - originY

	x = x + WIDTH / 2
	y = y + HEIGHT / 2

	left := x - int32(n.Radius)
	top := y - int32(n.Radius)



	return &sdl.Rect {
		X: left,
		Y: top,
		W: int32(n.Radius) * 2,
		H: int32(n.Radius) * 2,
	}
}

type Callbacks struct {
}

func (c *Callbacks) AddNode(id uint32) {
	lock.Lock()
	mynodes[id] = true
	CamNode = id
	lock.Unlock()
}

func (c *Callbacks) Eaten(eater, victim uint32) {
	if isMe(victim) {
	}
}

func (c *Callbacks) UpdateNode(x *agar.Node) {
	lock.Lock()
	if prev, exist := nodes[x.NodeID]; exist {
		prev.Update(x)
	} else {
		nodes[x.NodeID] = x
	}
	lock.Unlock()
}

func (c *Callbacks) RemoveNode(id uint32) {
	lock.Lock()
	delete(nodes, id)
	delete(mynodes, id)
	lock.Unlock()
	if id == CamNode {
		pickCamNode()
	}
	if len(mynodes) == 0 {
		log.Println("R.I.P")
	}
}

func (c *Callbacks) UpdateBorder(x *agar.Border) {
	border.Left = x.Left
	border.Top = x.Top
	border.Right = x.Right
	border.Bottom = x.Bottom
	if x.Server != "" {
		log.Println("Agar server version:", x.Server, x.Game)
	}
}

func (c *Callbacks) OnDisconnect() {
	log.Println("Disconnected!")
}

func render(renderer *sdl.Renderer) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()

	lock.Lock()
	for _, node := range nodes {
		if CamNode == node.NodeID {
			originX = node.Position.X
			originY = node.Position.Y
		}
		rect := getNodeRect(node)
		renderer.SetDrawColor(0, node.Color.R, node.Color.G, node.Color.B)
		if node.Virus() {
			renderer.DrawRect(rect)
		} else {
			renderer.FillRect(rect)
		}
	}
	lock.Unlock()
	renderer.Present()
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	var err error
	srv, _ := agar.FindServer()
	client, err = agar.NewClient(srv, &Callbacks{})
	if err != nil {
		log.Fatal(err)
	}
	window, e := sdl.CreateWindow("Agar.IO", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if e != nil {
		log.Fatal(e)
	}
	renderer, e := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if e != nil {
		log.Fatal(e)
	}


	client.SetName("Hello World!")
	client.Spawn()
//	client.Spectate()
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.MouseMotionEvent:
				for id, _ := range mynodes {
					vector_x := int32(-(WIDTH / 2))
					vector_y := int32(-(HEIGHT / 2))
					vector_x += t.X
					vector_y += t.Y
//					log.Println("Vector: ",vector_x, vector_y)
					client.Move(id, originX + vector_x, originY + vector_y)
				}
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case 'w':
					client.Eject()
				case 'q':
					client.Spawn()
				case ' ':
					client.Split()
				}
			}
		}
		render(renderer)
	}
}
