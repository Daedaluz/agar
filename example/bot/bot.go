package main

import(
	"github.com/daedaluz/agar"
	"github.com/gorilla/websocket"
	"math"
	"sort"
	"log"
)


type ByNodeID []*agar.Node

func (a ByNodeID) Len() int { return len(a) }
func (a ByNodeID) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByNodeID) Less(i, j int) bool { return a[i].NodeID > a[j].NodeID }

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
	var biggest = int16(0)
	for id, node := range c.nodes {
		if id == x.NodeID {
			continue
		}
		if !node.Virus() {
			if node.Radius > x.Radius {
				if (node.Radius < biggest || biggest == 0) {
					closest = node
					biggest = node.Radius
				}
			}
		}
	}
	if closest != nil {
	}
	return closest
}

func wmass(me, target float64) float64 {
	if target > (me-25) {
		return -500
	}
	//weight := -(math.Log((me - target) - 10)*10) + 20
	weight := 500 - target
	return weight
}


func wdist(dist float64) float64 {
	lvl := -267*math.Log2(dist)+820
	lvl = lvl / 0.3
	return math.Max(0, lvl)
}

func (c *Cb) GetNext(me *agar.Node) (*agar.Node, []*agar.Node) {
	var tpellets []*agar.Node
	var tmoving *agar.Node
	danger := make([]*agar.Node, 0, len(c.nodes))

	moving_closest := float64(0)
	pellet_closest := float64(0)

	for id, node := range c.nodes {
		if id == me.NodeID {
			continue
		}
		if node.Virus() {
			log.Println("Virus radius:", node.Radius)
			if node.Radius < (me.Radius - 20) {
				danger = append(danger, node)
			}
		} else {
			if node.IsMoving() {
				if node.Radius > me.Radius {
					danger = append(danger, node)
				} else {
					if node.Radius < (me.Radius - 35) {
						if dist := me.Distance(node); dist < moving_closest || moving_closest == 0 {
							tmoving = node
							moving_closest = dist
						}
					}
				}
			} else {
				log.Println("pellet radius:", node.Radius, me.Radius, me.Radius - node.Radius)
				if dist := me.Distance(node); dist < pellet_closest || pellet_closest == 0 {
					if node.Radius < (me.Radius - 20) {
						pellet_closest = dist
						tpellets = []*agar.Node{node}
					}
				} else if dist == pellet_closest {
					tpellets = append(tpellets, node)
				}
			}
		}
	}

	if tmoving != nil {
		return tmoving, danger
	}

	switch len(tpellets) {
	case 0:
		return nil, danger
	case 1:
	default:
		sort.Sort(ByNodeID(tpellets))
	}
	return tpellets[0], danger
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
		node := c.nodes[x.NodeID]
		if node.Name != "" && node.Name != "companionball<3" {
			log.Println("Distance:", me.Distance(x), node.Name)
		}
		target, danger := c.GetNext(me)
		if target != nil {
			influence := me.Position.NewVector(&target.Position)
			//influence.Norm()
			//influence.Scale(float64(me.Radius))
			if target.Distance(me) < float64(me.Radius * 3) && target.Radius < (me.Radius/2){
				if target.Direction.Length() != 0 {
					//c.cli.Eject()
				}
			}

			vDanger := agar.Vector{}
			for _, node := range danger {
				if node.NodeID != c.me && node.NodeID != target.NodeID {
					if node.Virus() {
						if node.DistanceSafe(me, 50) < 0 {
							vtmp := node.Position.NewVector(&me.Position)
							vtmp.Norm()
							vtmp.Scale(float64(me.Radius)*10)
							vDanger.Add(vtmp)
						}
					}else {
						if node.DistanceSafe(me, 400) < 0 {
							vtmp := node.Position.NewVector(&me.Position)
							vtmp.Norm()
							vtmp.Scale(wdist(vtmp.Length()))
							vDanger.Add(vtmp)
						}
					}
				}
			}

			influence.Add(&vDanger)

			posx, posy := me.Position.Move(influence)

			c.cli.Move(me.NodeID, posx, posy)
		} else {
			c.cli.Move(me.NodeID, 0, 0)
		}
	}
}

func (c *Cb) RemoveNode(id uint32) {
	delete(c.nodes, id)
}


