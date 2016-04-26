package goac

import (
	"encoding/json"
	"testing"
)

func TestGraph(t *testing.T) {
	g := NewGraph("admin")
	g.UseNegativeBuffer(true)
	if g.HavePath("admin", "A") {
		t.Error("A does not yet exist")
	}
	if g.HavePath("A", "A") {
		t.Error("A does not yet exist")
	}
	g.SetVertex(Vertex{Name: "A"})
	if !g.HavePath("admin", "A") {
		t.Error("admin should have power over all")
	}
	if !g.HavePath("A", "A") {
		t.Error("should have self control")
	}
	g.SetVertex(Vertex{Name: "B"})
	if g.HavePath("A", "B") {
		t.Error("nop")
	}
	g.SetVertex(Vertex{Name: "A", FullAssignments: []FullAssignment{
		FullAssignment{Elevate: "A", Over: "B"}}})
	if g.HavePath("A", "B") {
		t.Error("A can't just self assign")
	}
	g.SetVertex(Vertex{Name: "admin", FullAssignments: []FullAssignment{
		FullAssignment{Elevate: "A", Over: "B"}}})
	if !g.HavePath("A", "B") {
		t.Error("assigned")
	}
	g.SetVertex(Vertex{Name: "A", FullAssignments: []FullAssignment{
		FullAssignment{Elevate: "C", Over: "B"}}})
	if !g.HavePath("C", "B") {
		t.Error("A shoud be able to share B")
	}
	g.SetVertex(Vertex{Name: "admin", FullAssignments: []FullAssignment{
		FullAssignment{Elevate: "A", Over: ""},
		FullAssignment{Elevate: "", Over: "Z"}}})
	if g.HavePath("A", "Z") {
		t.Error("Empty string can't be assigned or used as an intermediary")
		//even without this check, as long as admin don't assign anything over "",
		//no other user could misconfigure and cause harm with it.
		//On the other hand, misspelling elevate and over...
	}

	g.SetVertex(Vertex{Name: "admin", FullAssignments: []FullAssignment{
		FullAssignment{Elevate: "A", Over: "I"},
		FullAssignment{Elevate: "I", Over: "J"}}})
	if !g.HavePath("A", "J") {
		t.Error("A -> I -> J")
	}
	bs, _ := json.Marshal(g.vs)
	t.Log(string(bs))

	g.rebuildTable()
	if !g.hasFullFrom.negativeBuffer {
		t.Error("negativeBuffer should not be reset")
	}

}

func TestSuperAdmin(t *testing.T) {
	g := NewGraph("admin")
	g.UseNegativeBuffer(true)
	g.SetVertex(Vertex{Name: "A"})
	if !g.HavePath("admin", "A") {
		t.Error("admin sould always have power")
	}
	g.SetVertex(Vertex{Name: "admin", FullAssignments: []FullAssignment{
		FullAssignment{Elevate: "super", Over: "admin"}}})
	if !g.HavePath("super", "A") {
		t.Log(g)
		t.Error("super should have power given by admin")
	}
}
