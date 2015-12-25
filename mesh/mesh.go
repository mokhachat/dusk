package mesh

import "github.com/go-gl/gl/v3.3-core/gl"

type Mesh struct {
	Vertices []float32
	Normals  []float32
	UVs      []float32
}

// newMesh is constructor for Mesh
func NewMesh(vertices []float32, uvs []float32) *Mesh {
	mesh := &Mesh{
		Vertices: vertices,
		Normals:  CalcNormals(vertices),
		UVs:      uvs,
	}
	return mesh
}

func (self *Mesh) GetVerticesNum() int32 {
	return (int32)(len(self.Vertices) / 3)
}

func (self *Mesh) StructVAO(program uint32) uint32 {
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo [3]uint32
	gl.GenBuffers(3, &vbo[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.UVs)*4, gl.Ptr(self.UVs), gl.STATIC_DRAW)
	uvAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vUv\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointer(uvAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Normals)*4, gl.Ptr(self.Normals), gl.STATIC_DRAW)
	norAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vNor\x00")))
	gl.EnableVertexAttribArray(norAttrib)
	gl.VertexAttribPointer(norAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	return vao
}

func (self *Mesh) Draw() {
	gl.DrawArrays(gl.TRIANGLES, 0, self.GetVerticesNum())
}

func CalcNormals(vertices []float32) []float32 {
	normals := make([]float32, len(vertices))
	for i := 0; i < len(vertices); i += 9 {
		v0 := [3]float32{vertices[i+0], vertices[i+1], vertices[i+2]}
		v1 := [3]float32{vertices[i+3], vertices[i+4], vertices[i+5]}
		v2 := [3]float32{vertices[i+6], vertices[i+7], vertices[i+8]}
		rv1 := [3]float32{v1[0] - v0[0], v1[1] - v0[1], v1[2] - v0[2]}
		rv2 := [3]float32{v2[0] - v0[0], v2[1] - v0[1], v2[2] - v0[2]}
		for k := 0; k < 3; k++ {
			normals[i+k*3+0] = rv1[1]*rv2[2] - rv1[2]*rv2[1]
			normals[i+k*3+1] = rv1[2]*rv2[0] - rv1[0]*rv2[2]
			normals[i+k*3+2] = rv1[0]*rv2[1] - rv1[1]*rv2[0]
		}
	}
	return normals
}
