package goac

import (
	"testing"
)

var A = VertexRef(0)
var B = VertexRef(1)
var f = VertexRef(2)
var C = VertexRef(3)
var n = VertexRef(4)

func testRefTable() vertexRefTable {
	vt := vertexRefTable{}
	vt.Reset(5)

	vt.Set(A, B)
	vt.Set(B, f)
	vt.Set(A, C)
	vt.Set(n, C)

	return vt
}
func testRefTablePathes() [][]bool {
	return [][]bool{
		{true, true, true, true, false},
		{false, true, true, false, false},
		{false, false, true, false, false},
		{false, false, false, true, false},
		{false, false, false, true, true},
	}
}

func evaluateTestRefTable(vt vertexRefTable, havePathes [][]bool, t *testing.T) {
	for elevated, list := range havePathes {
		for over, shouldHave := range list {
			result := vt.HavePath(VertexRef(elevated), VertexRef(over))
			if shouldHave != result {
				t.Error(elevated, "->", over, " HavePath: ", result)
			}
		}
	}
}

func TestSetHavePath(t *testing.T) {
	evaluateTestRefTable(testRefTable(), testRefTablePathes(), t)
}

func TestProposeHavePath(t *testing.T) {
	vt := testRefTable()
	if vt.Propose(n, C, f) {
		t.Fatal("n should not beable to do so")
	}
	if !vt.Propose(B, C, f) {
		t.Fatal("B should beable to do so")
	}
	havePathes := testRefTablePathes()
	havePathes[C][f] = true
	havePathes[n][f] = true
	evaluateTestRefTable(vt, havePathes, t)
	if !vt.Propose(n, C, f) {
		t.Fatal("n should beable to do so now, although it is a no-op")
	}
	evaluateTestRefTable(vt, havePathes, t)
}

func TestCirculairePath(t *testing.T) {
	vt := testRefTable()
	vt.Set(C, A)
	havePathes := testRefTablePathes()
	havePathes[C] = []bool{true, true, true, true, false}
	havePathes[n] = []bool{true, true, true, true, true}
	evaluateTestRefTable(vt, havePathes, t)

}
