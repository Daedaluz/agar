package agar

type Leaderboard struct {
	Name string
	Highlight uint32
}

type View struct {
	X float32
	Y float32
	Zoom float32
}

type Border struct {
	Left float64
	Top float64
	Right float64
	Bottom float64
	Game uint32
	Server string
}

func (b *Border) Height() float64 {
	return b.Bottom - b.Top
}

func (b *Border) Width() float64 {
	return b.Right - b.Left
}

type Eat struct {
	Eater uint32
	Victim uint32
}

type Color struct {
	R, G, B uint8
}

type Location struct {
	X int32
	Y int32
}

type Node struct {
	NodeID uint32
	Position Location
	Radius int16
	Color Color
	Flags uint8
	Face string
	Name string
}

func (u *Node) Update(n *Node) {
	u.Position = n.Position
	u.Radius = n.Radius
	u.Color = n.Color
	u.Flags = n.Flags
	if n.Face != "" {
		u.Face = n.Face
	}

	if n.Name != "" {
		u.Name = n.Name
	}
}

func (u *Node) Virus() bool {
	return (u.Flags & 1) > 0
}

func (u *Node) Skip4() bool {
	return (u.Flags & 2) > 0
}

func (u *Node) HasFace() bool {
	return (u.Flags & 4) > 0
}

func (u *Node) Agitated() bool {
	return (u.Flags & 16) > 0
}

