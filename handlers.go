package agar

import(
	"log"
	"io"
	"io/ioutil"
	"encoding/binary"
	"github.com/daedaluz/agar/intern"
)

type handler func(c *Client, i io.Reader)

type EaterHandler interface {
	Eaten(eater, victim uint32)
}

type NodeUpdater interface {
	UpdateNode(node *Node)
}

type NodeRemover interface {
	RemoveNode(id uint32)
}

func updateNodes(c *Client, i io.Reader) {
	var nBlobEats uint16
	if err := binary.Read(i, binary.LittleEndian, &nBlobEats); err != nil {
		log.Println("error reading nblob eats", err)
	}
	for x := uint16(0); x < nBlobEats; x++ {
		eat := Eat{}
		if err := binary.Read(i, binary.LittleEndian, &eat); err != nil {
			log.Println("error reading eat index", x, err)
		}
		if z, ok := c.cfg.(EaterHandler); ok {
			z.Eaten(eat.Eater, eat.Victim)
		}
	}
	for {
		var NodeID uint32
		if err := binary.Read(i, binary.LittleEndian, &NodeID); err != nil {
			log.Println("Error reading NodeID", err)
		}
		if NodeID == 0 {
			break
		}
		blobUpdate := &Node{}
		blobUpdate.NodeID = NodeID
		if err := binary.Read(i, binary.LittleEndian, &blobUpdate.Position); err != nil {
			log.Println("error reading position", err)
		}
		if err := binary.Read(i, binary.LittleEndian, &blobUpdate.Radius); err != nil {
			log.Println("error reading radius", err)
		}
		if err := binary.Read(i, binary.LittleEndian, &blobUpdate.Color); err != nil {
			log.Println("error reading color", err)
		}
		if err := binary.Read(i, binary.LittleEndian, &blobUpdate.Flags); err != nil {
			log.Println("error reading flags", err)
		}
		if blobUpdate.Skip4() {
			var toskip uint32
			binary.Read(i, binary.LittleEndian, &toskip)
			log.Println("Toskip:", toskip)
			if toskip < 100 {
				io.CopyN(ioutil.Discard, i, int64(toskip))
			}
		}
		if blobUpdate.HasFace() {
			data := make([]byte, 0, 30)
			var c byte
			for binary.Read(i, binary.LittleEndian, &c); c != 0; binary.Read(i, binary.LittleEndian, &c) {
				data = append(data, c)
			}
			blobUpdate.Face = string(data)
		}
		blobUpdate.Name = intern.ReadString(i)
		if x, ok := c.cfg.(NodeUpdater); ok {
			x.UpdateNode(blobUpdate)
		}
	}
	var nRemovals uint32
	binary.Read(i, binary.LittleEndian, &nRemovals)
	for x := uint32(0); x < nRemovals; x++ {
		var del uint32
		binary.Read(i, binary.LittleEndian, &del)
		if z, ok := c.cfg.(NodeRemover); ok {
			z.RemoveNode(del)
		}
	}
}

type ViewUpdater interface {
	UpdateView(w *View)
}

func updateView(c *Client, i io.Reader) {
	view := &View{}
	binary.Read(i, binary.LittleEndian, view)
	if x, ok := c.cfg.(ViewUpdater); ok {
		x.UpdateView(view)
	}
}

type AllCellResetter interface {
	ResetAllCells()
}

func resetAllCells(c *Client, i io.Reader) {
	if x, ok := c.cfg.(AllCellResetter); ok {
		x.ResetAllCells()
	}
}

type OwnCellResetter interface {
	ResetOwnCells()
}

func resetOwnCells(c *Client, i io.Reader) {
	if x, ok := c.cfg.(OwnCellResetter); ok {
		x.ResetOwnCells()
	}
}

type LineDrawer interface {
	DrawLine(x, y int16)
}

func drawLine(c *Client, i io.Reader) {
	var x int16
	var y int16
	binary.Read(i, binary.LittleEndian, x)
	binary.Read(i, binary.LittleEndian, y)
	if z, ok := c.cfg.(LineDrawer); ok {
		z.DrawLine(x, y);
	}
}

type NodeAdder interface {
	AddNode(id uint32)
}

func addNode(c *Client, i io.Reader) {
	var nodeid uint32
	binary.Read(i, binary.LittleEndian, &nodeid)
	if x, ok := c.cfg.(NodeAdder); ok {
		x.AddNode(nodeid)
	}
}

type LeaderboardUpdater interface {
	UpdateLeaderboard(board []Leaderboard)
}

func updateLeaderBoard(c *Client, i io.Reader) {
	var nLeaders uint32
	if err := binary.Read(i, binary.LittleEndian, &nLeaders); err != nil {
		log.Println("error reading number of leaders", err)
	}
	leaderboard := make([]Leaderboard, nLeaders)
	for x := uint32(0); x < nLeaders; x++ {
		binary.Read(i, binary.LittleEndian, &leaderboard[x].Highlight)
		leaderboard[x].Name = intern.ReadString(i)
		if leaderboard[x].Name == "" {
			leaderboard[x].Name = "An unnamed cell"
		}
	}
	if x, ok := c.cfg.(LeaderboardUpdater); ok {
		x.UpdateLeaderboard(leaderboard)
	}
}

func updateLeaderBoardTeam(c *Client, i io.Reader) {
	log.Println("UpdateLeaderboard")
}

type BorderUpdater interface {
	UpdateBorder(b *Border)
}
func setBorder(c *Client, i io.Reader) {
	ev := &Border {}
	if err := binary.Read(i, binary.LittleEndian, &ev.Left); err != nil {
		log.Println("error reading border", err)
	}
	if err := binary.Read(i, binary.LittleEndian, &ev.Top); err != nil {
		log.Println("error reading border", err)
	}
	if err := binary.Read(i, binary.LittleEndian, &ev.Right); err != nil {
		log.Println("error reading border", err)
	}
	if err := binary.Read(i, binary.LittleEndian, &ev.Bottom); err != nil {
		log.Println("error reading border", err)
	}
	if err := binary.Read(i, binary.LittleEndian, &ev.Game); err != nil {
	}
	ev.Server = intern.ReadString(i)
	if x, ok := c.cfg.(BorderUpdater); ok {
		x.UpdateBorder(ev)
	}
}

var stdHandlers = map[uint8]handler {
	16: updateNodes,
	17: updateView,
	18: resetAllCells,
	20: resetOwnCells,
	21: drawLine,
	32: addNode,
	49: updateLeaderBoard,
	50: updateLeaderBoardTeam,
	64: setBorder,
}
