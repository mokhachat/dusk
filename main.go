package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rhakt/dusk/shader"
	"github.com/rhakt/dusk/texture"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "dusk", nil, nil)

	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		panic(err)
	}
	//version := gl.GoStr(gl.GetString(gl.VERSION))
	//fmt.Println("OpenGL version", version)

	program, err := shader.Program(shader.VertexShader, shader.FragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 1000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("mProj\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	model := mgl32.Ident4()

	normal := camera.Mul4(model).Inv().Transpose()
	normalUniform := gl.GetUniformLocation(program, gl.Str("mNormal\x00"))
	gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])

	cameraUniform := gl.GetUniformLocation(program, gl.Str("mView\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	modelUniform := gl.GetUniformLocation(program, gl.Str("mModel\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	setUniform3f(program, "light.pos", -5.0, 0.0, 5.0)
	setUniform3f(program, "light.La", 0.3, 0.3, 0.3)
	setUniform3f(program, "light.Ld", 0.0, 1.0, 1.0)
	setUniform3f(program, "light.Ls", 0.0, 0.7, 0.9)
	setUniform3f(program, "mat.Ka", 0.3, 0.3, 0.3)
	setUniform3f(program, "mat.Kd", 0.0, 1.0, 1.0)
	setUniform3f(program, "mat.Ks", 0.0, 0.7, 0.9)
	setUniform1f(program, "mat.sh", 1.0)

	gl.BindFragDataLocation(program, 0, gl.Str("outColor\x00"))

	// Load the texture
	texture, err := texture.Load("data/texture.png")
	//texture, err := texture.Text("data/RictyDiminished-Regular.ttf", 128, "ポA1")
	if err != nil {
		panic(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo [3]uint32
	gl.GenBuffers(3, &vbo[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)
	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vPos\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeUVs)*4, gl.Ptr(cubeUVs), gl.STATIC_DRAW)
	uvAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vUv\x00")))
	gl.EnableVertexAttribArray(uvAttrib)
	gl.VertexAttribPointer(uvAttrib, 2, gl.FLOAT, false, 2*4, gl.PtrOffset(0))

	cubeNormals := make([]float32, len(cubeVertices))
	for i := 0; i < len(cubeVertices); i += 9 {
		v0 := [3]float32{cubeVertices[i+0], cubeVertices[i+1], cubeVertices[i+2]}
		v1 := [3]float32{cubeVertices[i+3], cubeVertices[i+4], cubeVertices[i+5]}
		v2 := [3]float32{cubeVertices[i+6], cubeVertices[i+7], cubeVertices[i+8]}
		rv1 := [3]float32{v1[0] - v0[0], v1[1] - v0[1], v1[2] - v0[2]}
		rv2 := [3]float32{v2[0] - v0[0], v2[1] - v0[1], v2[2] - v0[2]}
		for k := 0; k < 3; k++ {
			cubeNormals[i+k*3+0] = rv1[1]*rv2[2] - rv1[2]*rv2[1]
			cubeNormals[i+k*3+1] = rv1[2]*rv2[0] - rv1[0]*rv2[2]
			cubeNormals[i+k*3+2] = rv1[0]*rv2[1] - rv1[1]*rv2[0]
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[2])
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeNormals)*4, gl.Ptr(cubeNormals), gl.STATIC_DRAW)
	norAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vNor\x00")))
	gl.EnableVertexAttribArray(norAttrib)
	gl.VertexAttribPointer(norAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	angle := 0.0
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	frame := make(chan bool)

	tf := time.NewTicker(time.Second)
	defer tf.Stop()
	go func() {
		counter := 0
		for {
			select {
			case <-frame:
				counter++
			case <-tf.C:
				fmt.Println(counter)
				counter = 0
			}
		}
	}()

	t := time.NewTicker(16666 * time.Microsecond)
	defer t.Stop()
	go func() {
		for range t.C {
			angle += 0.01
		}
	}()

	for !window.ShouldClose() && window.GetKey(glfw.KeyEscape) != glfw.Press {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})
		normal := camera.Mul4(model).Inv().Transpose()

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, 6*2*3)

		window.SwapBuffers()
		glfw.PollEvents()

		frame <- true
	}
}

func setUniform1f(program uint32, name string, val float32) {
	uniform := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform1f(uniform, val)
}

func setUniform3f(program uint32, name string, v1, v2, v3 float32) {
	uniform := gl.GetUniformLocation(program, gl.Str(name+"\x00"))
	gl.Uniform3f(uniform, v1, v2, v3)
}

var cubeVertices = []float32{
	//  X, Y, Z
	// Bottom
	-1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, 1.0,
	-1.0, -1.0, 1.0,

	// Top
	-1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,
	-1.0, 1.0, 1.0,
	1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0,
	1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, 1.0,
	-1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	1.0, -1.0, -1.0,
	1.0, -1.0, -1.0,
	-1.0, 1.0, -1.0,
	1.0, 1.0, -1.0,

	// Left
	-1.0, -1.0, 1.0,
	-1.0, 1.0, -1.0,
	-1.0, -1.0, -1.0,
	-1.0, -1.0, 1.0,
	-1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0,

	// Right
	1.0, -1.0, 1.0,
	1.0, -1.0, -1.0,
	1.0, 1.0, -1.0,
	1.0, -1.0, 1.0,
	1.0, 1.0, -1.0,
	1.0, 1.0, 1.0,
}

var cubeUVs = []float32{
	// U, V
	// Bottom
	0.0, 0.0,
	1.0, 0.0,
	0.0, 1.0,
	1.0, 0.0,
	1.0, 1.0,
	0.0, 1.0,

	// Top
	0.0, 0.0,
	0.0, 1.0,
	1.0, 0.0,
	1.0, 0.0,
	0.0, 1.0,
	1.0, 1.0,

	// Front
	1.0, 0.0,
	0.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
	1.0, 1.0,

	// Back
	0.0, 0.0,
	0.0, 1.0,
	1.0, 0.0,
	1.0, 0.0,
	0.0, 1.0,
	1.0, 1.0,

	// Left
	0.0, 1.0,
	1.0, 0.0,
	0.0, 0.0,
	0.0, 1.0,
	1.0, 1.0,
	1.0, 0.0,

	// Right
	1.0, 1.0,
	1.0, 0.0,
	0.0, 0.0,
	1.0, 1.0,
	0.0, 0.0,
	0.0, 1.0,
}
