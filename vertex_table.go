package goac

type vertexRefTable struct {
	t              []map[VertexRef]bool
	negativeBuffer bool
}

func (v *vertexRefTable) Reset(size int) {
	if len(v.t) < size {
		v.t = append(v.t, make([]map[VertexRef]bool, size-len(v.t))...)
	} else {
		v.t = v.t[:size]
	}
	for e, mo := range v.t {
		if mo == nil || len(mo) != 0 {
			v.t[e] = make(map[VertexRef]bool)
		}
	}
}

func (v *vertexRefTable) Set(elevate, over VertexRef) {
	if v.negativeBuffer && !v.HavePath(elevate, over) {
		v.clearNegativeBuffer()
	}
	v.t[elevate][over] = true
}

func (v *vertexRefTable) HavePath(elevate, over VertexRef) bool {
	if elevate == over {
		return true
	}
	have := v.havePath(elevate, over, map[VertexRef]struct{}{})
	if have || v.negativeBuffer {
		v.t[elevate][over] = have //optional, can make same search faster
	}
	return have
}

func (v *vertexRefTable) havePath(elevate, over VertexRef, visted map[VertexRef]struct{}) bool {
	b, ok := v.t[elevate][over]
	if ok {
		return b
	}
	visted[elevate] = struct{}{}
	for child, b := range v.t[elevate] {
		if b {
			_, seen := visted[child]
			if !seen && v.havePath(child, over, visted) {
				return true
			}
		}
	}
	return false
}

func (v *vertexRefTable) Propose(proposer, elevate, over VertexRef) bool {
	if v.HavePath(proposer, over) {
		v.Set(elevate, over)
		return true
	}
	return false
}

func (v *vertexRefTable) UseNegativeBuffer(b bool) {
	if !b {
		v.clearNegativeBuffer()
	}
	v.negativeBuffer = b
}
func (v *vertexRefTable) clearNegativeBuffer() {
	if v.negativeBuffer {
		for _, mo := range v.t {
			for k, m := range mo {
				if !m {
					delete(mo, k)
				}
			}
		}
	}
}
