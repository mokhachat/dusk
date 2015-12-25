package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/rhakt/dusk/mesh"
	"github.com/rhakt/dusk/shader"
	"github.com/rhakt/dusk/texture"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
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

	program, err := shader.Program(shader.VSPhong, shader.PSPhong)

	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl.Perspective(mgl.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 1000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("mProj\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl.LookAtV(mgl.Vec3{3, 3, 3}, mgl.Vec3{0, 0, 0}, mgl.Vec3{0, 1, 0})
	model := mgl.Ident4()

	normal := camera.Mul4(model).Inv().Transpose()
	normalUniform := gl.GetUniformLocation(program, gl.Str("mNormal\x00"))
	gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])

	cameraUniform := gl.GetUniformLocation(program, gl.Str("mView\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	modelUniform := gl.GetUniformLocation(program, gl.Str("mModel\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	setUniform3f(program, "LP", -5.0, 0.0, 5.0)
	setUniform3f(program, "LI", 0.7, 0.7, 0.7)
	setUniform3f(program, "Ka", 0.3, 0.3, 0.3)
	setUniform3f(program, "Kd", 0.0, 1.0, 1.0)
	setUniform3f(program, "Ks", 0.0, 0.7, 0.9)
	setUniform1f(program, "Sh", 1.0)

	gl.BindFragDataLocation(program, 0, gl.Str("outColor\x00"))

	// Load the texture
	tex, err := texture.Load("data/texture.png")
	//texture, err := texture.Text("data/RictyDiminished-Regular.ttf", 128, "ポA1")
	if err != nil {
		panic(err)
	}

	// Configure the vertex data
	cube := mesh.NewMesh(mesh.CubeVertices, mesh.CubeUVs)
	vao := cube.StructVAO(program)

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

	mode := window.GetInputMode(glfw.CursorMode)
	prevCX, prevCY := window.GetCursorPos()
	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		//fmt.Printf("mouser: (%f, %f)\n", xpos, ypos)
		if mode == glfw.CursorDisabled {
			fmt.Printf("mouse: (%f, %f)\n", xpos-prevCX, ypos-prevCY)
			window.SetCursorPos(prevCX, prevCY)
		}
	})

	// window.GetKey(glfw.KeyEscape) != glfw.Press
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// action: Release, Press, Repeat
		// mod: ModShift, ModControl ModAlt, ModSuper
		if key == glfw.KeyEscape && action == glfw.Press {
			next := glfw.CursorDisabled
			if mode == glfw.CursorDisabled {
				next = glfw.CursorNormal
			} else {
				prevCX, prevCY = window.GetCursorPos()
			}
			window.SetInputMode(glfw.CursorMode, next)
			mode = next
		}
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		// button: MouseButtonLeft, MouseButtonRight, MouseButtonMiddle
		btn := map[glfw.MouseButton]string{
			glfw.MouseButtonLeft:   "left",
			glfw.MouseButtonRight:  "right",
			glfw.MouseButtonMiddle: "middle",
		}
		act := map[glfw.Action]string{
			glfw.Press:   "press",
			glfw.Release: "release",
			glfw.Repeat:  "repeat",
		}
		fmt.Println(btn[button] + act[action])
	})

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		model = mgl.HomogRotate3D(float32(angle), mgl.Vec3{0, 1, 0})
		normal := camera.Mul4(model).Inv().Transpose()

		// Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, tex)

		cube.Draw()

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
