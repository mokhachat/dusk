package main

import (
	"fmt"
	"runtime"
    //"strings"
	"time"
    "io/ioutil"
    
    "github.com/rhakt/dusk/mesh"
	"github.com/rhakt/dusk/shader"
	"github.com/rhakt/dusk/texture"
    schema "github.com/rhakt/dusk/model"
    
    "github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	mgl "github.com/go-gl/mathgl/mgl32"
    
    "github.com/rhakt/lz4"
    
    //flatbuffers "github.com/google/flatbuffers/go"
)

const windowWidth = 800
const windowHeight = 600

func init() {
	runtime.LockOSThread()
}


type MeshInfo struct {
    MeshMatrices []mgl.Mat4
    BoneMatrices [][]mgl.Mat4
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

	program, err := shader.Program(shader.VSPhong3, shader.PSPhong3)

	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl.Perspective(mgl.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 1000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("mProj\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl.LookAtV(mgl.Vec3{160, 160, 160}, mgl.Vec3{0, 80, 0}, mgl.Vec3{0, 1, 0})
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
    
    mmUniform := gl.GetUniformLocation(program, gl.Str("meshMatrix\x00"))
    
    /*ubi := gl.GetUniformBlockIndex(program, gl.Str("frameInfo\x00"));
    var bmUniform uint32
    gl.GenBuffers(1, &bmUniform)
    */
    bmUniform := gl.GetUniformLocation(program, gl.Str("boneMatrices\x00"))
    fmt.Println("bmUniform", bmUniform, mmUniform)

	setUniform3f(program, "LP", -300.0, 150.0, 300.0)
	setUniform3f(program, "LI", 0.7, 0.7, 0.7)
	setUniform3f(program, "Ka", 0.3, 0.3, 0.3)
	setUniform3f(program, "Kd", 1.0, 1.0, 1.0)
	setUniform3f(program, "Ks", 0.8, 0.7, 0.9)
	setUniform1f(program, "Sh", 1.0)

	gl.BindFragDataLocation(program, 0, gl.Str("outColor\x00"))

	// Load the texture
	//tex, err := texture.Load("data/texture.png")
    //tex, err := texture.Text("data/RictyDiminished-Regular.ttf", 128, "ポA1")
	/*if err != nil {
		panic(err)
	}*/
    
    src, err := ioutil.ReadFile("data/unitychan.rchr")
    if err != nil {
        panic(err)
    }
    
    dst := make([]byte, len(src) * 10)
    dstSize, err := lz4.Decompress(src, dst)
    if err != nil {
        panic(err)
    }
    dst = dst[:dstSize]
    
    scene := schema.GetRootAsScene(dst, 0)
    
    var meshes []*mesh.Mesh
    var texes []uint32
    for i := 0; i < scene.MeshesLength(); i++ {
        var m schema.Mesh
        ok := scene.Meshes(&m, i)
        if !ok {
            break
        }
        ver := make([]float32, m.VerticesLength())
        for k := 0; k < m.VerticesLength(); k++ {
            ver[k] = m.Vertices(k)
        }
        nor := make([]float32, m.NormalsLength())
        for k := 0; k < m.NormalsLength(); k++ {
            nor[k] = m.Normals(k)
        }
        uv := make([]float32, m.UvsLength())
        for k := 0; k < m.UvsLength() / 2; k++ {
            uv[k*2] = m.Uvs(k*2)
            uv[k*2+1] = 1.0 - m.Uvs(k*2+1)
        }
        col := make([]float32, m.ColorsLength())
        for k := 0; k < m.ColorsLength(); k++ {
            col[k] = m.Colors(k)
        }
        ind := make([]uint32, m.IndicesLength())
        for k := 0; k < m.IndicesLength(); k++ {
            ind[k] = uint32(m.Indices(k))
            /*if(ind[k] != uint32(k)) {
                fmt.Println("mesh:",  i, "  ", ind[k], " != ", k);
            }*/
        }
        
        texname := string(m.Texture())
        //texname = strings.Replace(texname, ".tga", ".png", 1)
        //fmt.Println(texname);
        tex, err := texture.Load("./data/texture/" + texname)
        if err != nil {
		  panic(err)
        }
        
        //fmt.Println(m.VerticesLength() * 4 / 3, m.BoneIndicesLength(), m.BoneWeightsLength());
        
        bi := make([]uint32, m.VerticesLength() * 4 / 3)
        if m.BoneIndicesLength() == 0 {
            for k, _ := range bi {
                bi[k] = 0
            }
        } else {
            for k, _ := range bi {
                bi[k] = uint32(m.BoneIndices(k))
                if(bi[k] > 127) {
                    fmt.Println(i, k, bi[k])
                } 
            }
        }
        
        
        bw := make([]float32, m.VerticesLength() * 4 / 3)
        if m.BoneWeightsLength() == 0 {
            for k, _ := range bw {
                bw[k] = 0.25
            }
        } else {
            for k := 0; k < m.BoneWeightsLength(); k++ {
                bw[k] = m.BoneWeights(k)
                if(bw[k] < 0 || bw[k] > 1) {
                    fmt.Println(i, k, bw[k])
                }
            }    
        }
        
        
        //newmesh := mesh.NewMesh2(ver, nor, col, ind)
        //newmesh := mesh.NewMesh4(ver, col, ind)
        newmesh := mesh.NewMesh5(ver, nor, uv, col, ind, bi, bw)
        newmesh.StructVAO3(program)
        
        /*if(strings.Contains(texname, "cheek")){
            meshes2 = append(meshes2, newmesh)
            texes2 = append(texes2, tex)
            fmt.Println("alpha ->", texname);
        }*/
        meshes = append(meshes, newmesh)
        texes = append(texes, tex)    
    }
    
    var a schema.Anim
    ok := scene.Animes(&a, 1)
    if !ok {
        fmt.Println("baka")
    }
    frameInfo := make([]MeshInfo, a.MeshesLength())
    for k := 0; k < a.MeshesLength(); k++ {
        var af schema.AnimFrame
        ok := a.Meshes(&af, k)
        if !ok {
            break
        }
        mml := af.MeshMatricesLength()
        bml := af.BoneMatricesLength()
        var vmm []mgl.Mat4
        for j := 0; j < mml; j++ {
            var f schema.Frame
            ok := af.MeshMatrices(&f, j)
            if !ok {
                break
            }
            var matrix [16]float32
            for l := 0; l < 16; l++ {
                matrix[l] = f.Data(l)
            }
            vmm = append(vmm, mgl.Mat4(matrix))
        }
        var vbm [][]mgl.Mat4
        for j := 0; j < bml; j++ {
            var f schema.Frame
            ok := af.BoneMatrices(&f, j)
            if !ok {
                break
            }
            var matrices []mgl.Mat4
            fmt.Println(k, j, f.DataLength())
            for l := 0; l < f.DataLength() ; l+=16 {
                var matrix [16]float32
                for ll := 0; ll < 16 ; ll++ {
                    matrix[ll] = f.Data(l + ll)
                }
                matrices = append(matrices, mgl.Mat4(matrix))
            }
            vbm = append(vbm, matrices)
        }
        //fmt.Println("debug", k)
        frameInfo[k] = MeshInfo{vmm, vbm}
    }
    

	//cube := mesh.NewMesh(mesh.CubeVertices, mesh.CubeUVs, mesh.CubeIndices)
    //cube.StructVAO(program)
    
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
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
				//fmt.Println(counter)
				counter = 0
			}
		}
	}()

	t := time.NewTicker(16666 * time.Microsecond)
	defer t.Stop()
	go func() {
		for range t.C {
			angle += 0.00 // po
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
    
    frametime := 0
    
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
        if(act[action] == "press") { 
            if(btn[button] == "left") { frametime++ }
            if(btn[button] == "right" && frametime > 0) { frametime-- }
        }
        
	})
    
    //gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
    gl.Enable(gl.BLEND)
    //gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
    gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA);
    //gl.Disable(gl.CULL_FACE);
    for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		model = mgl.HomogRotate3D(float32(angle), mgl.Vec3{0, 1, 0})
		normal := camera.Mul4(model).Inv().Transpose()
        
        // Render
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
		gl.UniformMatrix4fv(normalUniform, 1, false, &normal[0])
        
       //gl.ActiveTexture(gl.TEXTURE0)  
       //gl.BindTexture(gl.TEXTURE_2D, tex)
        
        //cube.Draw()
        //gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
        for k, m := range(meshes) {
            //if(k != 9) { continue }
            mm := frameInfo[k].MeshMatrices
            mmmat := mgl.Ident4()
            if len(mm) > 0 {
                mmmat = mm[frametime % len(mm)]
            }
            gl.UniformMatrix4fv(mmUniform, 1, false, &mmmat[0])
            
            bm := frameInfo[k].BoneMatrices
            var bmmat [128]mgl.Mat4
            for j, _ := range(bmmat) {
                bmmat[j] = mgl.Ident4()
                if len(bm) > 0 && j < len(bm[frametime % (len(bm))]) {
                    bmmat[j] = bm[frametime % (len(bm))][j]
                }
            }
            //gl.BindBuffer(gl.UNIFORM_BUFFER, bmUniform)
            //gl.BufferData(gl.UNIFORM_BUFFER, 4*16*128, gl.Ptr(&bmmat[0][0]), gl.DYNAMIC_DRAW)
            //gl.BindBufferBase(gl.UNIFORM_BUFFER, ubi, bmUniform)
            /*for j, _ := range(bmmat) {
                gl.UniformMatrix4fv(bmUniform + int32(j), 1, false, &bmmat[j][0])
            }*/
            gl.UniformMatrix4fv(bmUniform, 128, false, &bmmat[0][0])
            
            gl.ActiveTexture(gl.TEXTURE0)
            gl.BindTexture(gl.TEXTURE_2D, texes[k])
            m.Draw()
        }
        //gl.DepthMask(false);
        /*gl.BlendFunc(gl.ONE, gl.ONE_MINUS_SRC_ALPHA);
        for k, m := range(meshes2) {
            gl.ActiveTexture(gl.TEXTURE0)
            gl.BindTexture(gl.TEXTURE_2D, texes2[k])
            m.Draw()
        }*/
        //gl.DepthMask(true);
        
		window.SwapBuffers()
		glfw.PollEvents()
        
        frametime++
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
