// automatically generated, do not modify

package model

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Scene struct {
	_tab flatbuffers.Table
}

func GetRootAsScene(buf []byte, offset flatbuffers.UOffsetT) *Scene {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Scene{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Scene) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Scene) Meshes(obj *Mesh, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 4
		x = rcv._tab.Indirect(x)
	if obj == nil {
		obj = new(Mesh)
	}
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *Scene) MeshesLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func SceneStart(builder *flatbuffers.Builder) { builder.StartObject(1) }
func SceneAddMeshes(builder *flatbuffers.Builder, meshes flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(meshes), 0) }
func SceneStartMeshesVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT { return builder.StartVector(4, numElems, 4)
}
func SceneEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
