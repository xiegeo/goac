package goac

type Graph struct {
	vs          []Vertex
	byName      map[string]VertexRef
	hasFullFrom vertexRefTable
}

type VertexRef int
