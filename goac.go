package goac

import (
	"fmt"
	"sync"
)

//Name denotes a string that has been normalized and unambiguous to use as graph
//vertex identifier. This cast should be made by user code that moves data into
//graphs.
type Name string

type Graph struct {
	admin         Name
	vs            []Vertex
	byName        map[Name]VertexRef
	hasFullFrom   vertexRefTable
	buildRefTable *sync.Once
	resetLock     sync.RWMutex
}

type VertexRef int

func NewGraph(admin Name) *Graph {
	g := &Graph{
		admin:  admin,
		byName: make(map[Name]VertexRef),
	}
	g.SetVertex(Vertex{Name: "Void Vertex"}) //Protect against uncaught errors that return 0.
	//This makes sure that all not found vertexes will point to the Void Vertex.
	//Which is also protected from being assigned power.
	g.SetVertex(Vertex{Name: admin})
	return g
}

func (g *Graph) resetBuildRefTable() {
	g.resetLock.Lock()
	g.buildRefTable = &sync.Once{}
	g.resetLock.Unlock()
}

func (g *Graph) GetVertex(name Name) *Vertex {
	v, ok := g.byName[name]
	if !ok {
		return nil
	}
	return &g.vs[v]
}

func (g *Graph) SetVertex(v Vertex) {
	if v.Name == "" {
		return
	}
	ref, ok := g.byName[v.Name]
	if !ok {
		ref = VertexRef(len(g.vs))
		g.vs = append(g.vs, v)
		g.byName[v.Name] = ref
	} else {
		g.vs[ref] = v
	}
	for _, a := range v.FullAssignments {
		_, ok := g.byName[a.Elevate]
		if !ok {
			g.SetVertex(Vertex{Name: a.Elevate})
		}
		_, ok = g.byName[a.Over]
		if !ok {
			g.SetVertex(Vertex{Name: a.Over})
		}
	}
	g.resetBuildRefTable()
}

func (g *Graph) rebuildTable() {
	g.hasFullFrom.Reset(len(g.byName))
	adminVertex := g.vs[g.byName[g.admin]]
	for _, a := range adminVertex.FullAssignments {
		e := g.byName[a.Elevate]
		if e == VertexRef(0) {
			//protect against giving the void vertex any power
		} else {
			g.hasFullFrom.Set(g.byName[a.Elevate], g.byName[a.Over])
		}
	}
	triples := [][3]VertexRef{}
	adminRef := g.byName[g.admin]
	for _, v := range g.vs {
		if v.Name != g.admin {
			p := g.byName[v.Name]
			for _, a := range v.FullAssignments {
				e := g.byName[a.Elevate]
				if e == VertexRef(0) {
					//protect against giving the void vertex any power
				} else {
					triples = append(triples, [3]VertexRef{p, e, g.byName[a.Over]})
				}
			}
			//all refs are also controlled by admin
			g.hasFullFrom.Set(adminRef, p)
		}
	}
	if len(triples) != 0 {
		lastApplied := 0
		i := 0
		applied := true
		for applied || i != lastApplied {
			triple := triples[i]
			applied = g.hasFullFrom.Propose(triple[0], triple[1], triple[2])
			if applied {
				if len(triples) == 1 {
					break
				}
				triples[i] = triples[len(triples)-1]
				triples = triples[:len(triples)-1]
				lastApplied = i + 1
			}
			lastApplied = lastApplied % len(triples)
			i = (i + 1) % len(triples)
		}
	}
}

//HavePath tests if elevate node is above the over node.
func (g *Graph) HavePath(elevate, over Name) bool {
	b, _ := g.HavePathDebug(elevate, over)
	return b
}

//HavePathDebug is HavePath with extra error information
func (g *Graph) HavePathDebug(elevate, over Name) (bool, error) {
	g.resetLock.RLock()
	g.buildRefTable.Do(g.rebuildTable)

	e, ok := g.byName[elevate]
	if !ok {
		g.resetLock.RUnlock()
		return false, fmt.Errorf("%v is not in graph", elevate)
	}
	o, ok2 := g.byName[over]
	if !ok2 {
		g.resetLock.RUnlock()
		return false, fmt.Errorf("%v is not in graph", over)
	}
	havePath := g.hasFullFrom.HavePath(e, o)
	g.resetLock.RUnlock()
	return havePath, nil
}

func (g *Graph) UseNegativeBuffer(b bool) {
	g.resetLock.Lock()
	g.hasFullFrom.UseNegativeBuffer(b)
	g.resetLock.Unlock()
}
