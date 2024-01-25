package groups

import (
	"sort"

	m "roytracer/math"
	"roytracer/shape"
)

type (
	Group struct {
		Shapes []ShapeEntry
		Tf     m.Mat4
	}

	ShapeEntry struct {
		Shape    shape.Shape
		Parent   *ShapeEntry
		Children []ShapeEntry
	}
)

func NewGroup() Group {
	return Group{
		Tf:     m.Mat4Ident(),
		Shapes: []ShapeEntry{},
	}
}

func (g *Group) AddChild(s shape.Shape) {
	g.Shapes = append(g.Shapes, ShapeEntry{Shape: s})
}

func (g *Group) Contains(s shape.Shape) bool {
	for _, se := range g.Shapes {
		if se.Contains(s) {
			return true
		}
	}
	return false
}

func (g *Group) localIntersect(r shape.Ray) []shape.Intersection {
	xs := []shape.Intersection{}
	for _, se := range g.Shapes {
		xs = append(xs, se.Shape.Intersect(r)...)
	}
	sort.Slice(xs, func(i, j int) bool { return xs[i].T < xs[j].T })
	return xs
}

func (se *ShapeEntry) Contains(s shape.Shape) bool {
	if se.Shape == s {
		return true
	}

	for _, child := range se.Children {
		if child.Contains(s) {
			return true
		}
	}
	return false
}
