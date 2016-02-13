package goac

type vertexRefTable struct {
	t []map[VertexRef]struct{}
}

func (v *vertexRefTable) Reset(size int) {
	if len(v.t) < size {
		v.t = append(v.t, make([]map[VertexRef]struct{}, size-len(v.t))...)
	} else {
		v.t = v.t[:size]
	}
	for e, mo := range v.t {
		if mo == nil || len(mo) != 0 {
			v.t[e] = make(map[VertexRef]struct{})
		}
	}
}

func (v *vertexRefTable) Set(elevate, over VertexRef) {
	v.t[elevate][over] = struct{}{}
}

func (v *vertexRefTable) HavePath(elevate, over VertexRef) bool {
	if elevate == over {
		return true
	}
	have := v.havePath(elevate, over, map[VertexRef]struct{}{})
	if have {
		v.t[elevate][over] = struct{}{} //optional, can make same search faster
	}
	return have
}

func (v *vertexRefTable) havePath(elevate, over VertexRef, visted map[VertexRef]struct{}) bool {
	_, ok := v.t[elevate][over]
	if ok {
		return true
	}
	visted[elevate] = struct{}{}
	for child, _ := range v.t[elevate] {
		_, seen := visted[child]
		if !seen && v.havePath(child, over, visted) {
			return true
		}
	}
	return false
}

func (v *vertexRefTable) Propose(proposer, elevate, over VertexRef) bool {
	_, ok := v.t[elevate][over]
	if ok {
		return true
	}
	if v.HavePath(proposer, over) {
		v.Set(elevate, over)
		return true
	}
	return false
}
