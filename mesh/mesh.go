package mesh

import "github.com/go-gl/gl/v3.3-core/gl"
import "fmt"

type Mesh struct {
	Vertices []float32
	Normals  []float32
	UVs      []float32
	Indices  []uint32
    Colors   []float32
    BoneIndices []uint32
    BoneWeights []float32
	vao      uint32
	vbo      [6]uint32
	ibo      uint32
}

// newMesh is constructor for Mesh
func NewMesh(vertices []float32, uvs []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		Vertices: vertices,
		Normals:  CalcNormals(vertices, indices),
		UVs:      uvs,
		Indices:  indices,
	}
	return mesh
}

func NewMesh2(vertices []float32, normals []float32, colors []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		Vertices: vertices,
		Normals:  normals,
		Colors:   colors,
		Indices:  indices,
	}
	return mesh
}

func NewMesh3(vertices []float32, normals []float32, uvs []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		Vertices: vertices,
		Normals:  normals,
		UVs:      uvs,
		Indices:  indices,
	}
	return mesh
}

func NewMesh4(vertices []float32, colors []float32, indices []uint32) *Mesh {
	mesh := &Mesh{
		Vertices: vertices,
		Normals:  CalcNormals(vertices, indices),
		Colors:   colors,
		Indices:  indices,
	}
	return mesh
}

func NewMesh5(
    vertices []float32, 
    normals []float32, 
    uvs []float32, 
    colors []float32, 
    indices []uint32,
    boneindices []uint32,
    boneweights []float32) *Mesh {
	mesh := &Mesh{
		Vertices: vertices,
		//Normals:  CalcNormals(vertices, indices),
        Normals:  normals,
        UVs:      uvs,
		Colors:   colors,
		Indices:  indices,
        BoneIndices: boneindices,
        BoneWeights: boneweights,
	}
	return mesh
}

func (self *Mesh) GetVerticesNum() int32 {
	return int32(len(self.Vertices) / 3)
}

func (self *Mesh) StructVAO(program uint32) /*uint32*/ {
	//var vao uint32
	gl.GenVertexArrays(1, &self.vao)
	gl.BindVertexArray(self.vao)

	//var vbo [3]uint32
	gl.GenBuffers(int32(len(self.vbo)), &self.vbo[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	//var ibo uint32
	gl.GenBuffers(1, &self.ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, self.ibo)
    gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(self.Indices)*4, gl.Ptr(self.Indices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.UVs)*4, gl.Ptr(self.UVs), gl.STATIC_DRAW)
	uvAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vUv\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointer(uvAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Normals)*4, gl.Ptr(self.Normals), gl.STATIC_DRAW)
	norAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vNor\x00")))
	gl.EnableVertexAttribArray(norAttrib)
	gl.VertexAttribPointer(norAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	//return vao
}

func (self *Mesh) StructVAO2(program uint32) /*uint32*/ {
	//var vao uint32
	gl.GenVertexArrays(1, &self.vao)
	gl.BindVertexArray(self.vao)

	//var vbo [3]uint32
	gl.GenBuffers(int32(len(self.vbo)), &self.vbo[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	//var ibo uint32
	gl.GenBuffers(1, &self.ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, self.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(self.Indices)*4, gl.Ptr(self.Indices), gl.STATIC_DRAW)
    fmt.Printf("%d\n", len(self.Indices))
	
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
    
	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Colors)*4, gl.Ptr(self.Colors), gl.STATIC_DRAW)
	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vCol\x00")))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Normals)*4, gl.Ptr(self.Normals), gl.STATIC_DRAW)
	norAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vNor\x00")))
	gl.EnableVertexAttribArray(norAttrib)
	gl.VertexAttribPointer(norAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	//return vao
}

func (self *Mesh) StructVAO3(program uint32) {
	gl.GenVertexArrays(1, &self.vao)
	gl.BindVertexArray(self.vao)

	gl.GenBuffers(int32(len(self.vbo)), &self.vbo[0])
    
	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Vertices)*4, gl.Ptr(self.Vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &self.ibo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, self.ibo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(self.Indices)*4, gl.Ptr(self.Indices), gl.STATIC_DRAW)
    //fmt.Printf("%d\n", len(self.Indices))
	
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
    
    gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.UVs)*4, gl.Ptr(self.UVs), gl.STATIC_DRAW)
	uvAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vUv\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointer(uvAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))
    
	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Colors)*4, gl.Ptr(self.Colors), gl.STATIC_DRAW)
	colAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vCol\x00")))
	gl.EnableVertexAttribArray(colAttrib)
	gl.VertexAttribPointer(colAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[3])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.Normals)*4, gl.Ptr(self.Normals), gl.STATIC_DRAW)
	norAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vNor\x00")))
	gl.EnableVertexAttribArray(norAttrib)
	gl.VertexAttribPointer(norAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
    
    //fmt.Println("po");
    //fmt.Println(len(self.BoneIndices));
    //fmt.Println(len(self.BoneWeights));
    
    gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[4])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.BoneIndices)*4, gl.Ptr(self.BoneIndices), gl.STATIC_DRAW)
	biAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vBI\x00")))
	gl.EnableVertexAttribArray(biAttrib)
    gl.VertexAttribPointer(biAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
    
    gl.BindBuffer(gl.ARRAY_BUFFER, self.vbo[5])
	gl.BufferData(gl.ARRAY_BUFFER, len(self.BoneWeights)*4, gl.Ptr(self.BoneWeights), gl.STATIC_DRAW)
	bwAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vBW\x00")))
	gl.EnableVertexAttribArray(bwAttrib)
	gl.VertexAttribPointer(bwAttrib, 4, gl.FLOAT, false, 4*4, gl.PtrOffset(0))
    
}

func (self *Mesh) Draw() {
	gl.BindVertexArray(self.vao)
	//gl.DrawArrays(gl.TRIANGLES, 0, self.GetVerticesNum())
	gl.DrawElements(gl.TRIANGLES, int32(len(self.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}

func CalcNormals(vertices []float32, indices []uint32) []float32 {
	normals := make([]float32, len(vertices))
	// TODO: [improve] calc duplication
	for j := 0; j < len(indices); j += 3 {
		i1 := indices[j] * 3
		i2 := indices[j+1] * 3
		i3 := indices[j+2] * 3
		v0 := [3]float32{vertices[i1+0], vertices[i1+1], vertices[i1+2]}
		v1 := [3]float32{vertices[i2+0], vertices[i2+1], vertices[i2+2]}
		v2 := [3]float32{vertices[i3+0], vertices[i3+1], vertices[i3+2]}
		rv1 := [3]float32{v1[0] - v0[0], v1[1] - v0[1], v1[2] - v0[2]}
		rv2 := [3]float32{v2[0] - v0[0], v2[1] - v0[1], v2[2] - v0[2]}
		for _, k := range []uint32{i1, i2, i3} {
			normals[k+0] = rv1[1]*rv2[2] - rv1[2]*rv2[1]
			normals[k+1] = rv1[2]*rv2[0] - rv1[0]*rv2[2]
			normals[k+2] = rv1[0]*rv2[1] - rv1[1]*rv2[0]
		}
	}
	return normals
}
