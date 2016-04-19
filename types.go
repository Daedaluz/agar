package agar

import(
	"math"
)

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

func (l *Location) Move(v *Vector) (nx, ny int32) {
	return l.X + int32(v.X), l.Y + int32(v.Y)
}

func (l *Location) NewVector(dest *Location) *Vector {
	return &Vector {
		float64(dest.X - l.X),
		float64(dest.Y - l.Y),
	}
}

type Vector struct {
	X, Y float64
}

func (v *Vector) Norm() {
	length := v.Length()
	v.X = v.X / length
	v.Y = v.Y / length
}

func (v *Vector) Add(v2 *Vector) {
	v.X += v2.X
	v.Y += v2.Y
}

func (v *Vector) Length() float64 {
	cx := float64(v.X)
	cy := float64(v.Y)
	return math.Sqrt(cx*cx + cy*cy)
}

func (v *Vector) Rotate(deg float64) {
	_deg := (math.Pi / 180) * deg
	xp := (v.X * math.Cos(_deg)) - (v.Y * math.Sin(_deg))
	yp := (v.X * math.Sin(_deg)) + (v.Y * math.Cos(_deg))
	v.X = xp
	v.Y = yp
}

func (v *Vector) Scale(mul float64) {
	v.X = v.X*mul
	v.Y = v.Y*mul
}

type Node struct {
	NodeID uint32
	Position Location
	Radius int16
	Color Color
	Flags uint8
	Face string
	Name string

	Direction Vector
}

func (u *Node) IsMoving() bool {
	if u.Direction.Length() > 0 {
		return true
	}
	return false
}

func (u *Node) Distance(n *Node) float64 {
	dx, dy := n.Position.X - u.Position.X, n.Position.Y - u.Position.Y
	dist := math.Sqrt(float64(dx*dx) + float64(dy*dy))
	return dist
}

func (u *Node) DistanceSafe(n *Node, safe float64) float64 {
	dx, dy := n.Position.X - u.Position.X, n.Position.Y - u.Position.Y
	dist := math.Sqrt(float64(dx*dx) + float64(dy*dy))
	dist = dist - safe - float64(u.Radius) - float64(n.Radius)
	return dist
}



func (u *Node) Tangent(n *Node, safe int32) (tx, ty int32) {
	dx, dy := n.Position.X - u.Position.X, n.Position.Y - u.Position.Y
	dist := math.Sqrt(float64(dx*dx) + float64(dy*dy))
	vecx, vecy := float64(dx)/dist, float64(dy)/dist

	dist = dist - float64(u.Radius) - float64(n.Radius)
	dist = dist - float64(safe)
	vecx, vecy = vecx*dist, vecy*dist

	tx, ty = u.Position.X + int32(vecx), u.Position.Y + int32(vecy)
	return
}

func (u *Node) Update(n *Node) {
	u.Direction.X, u.Direction.Y = float64(n.Position.X - u.Position.X), float64(n.Position.Y - u.Position.Y)
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

