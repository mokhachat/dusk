package shader

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// Program po
func Program(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compile(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compile(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compile(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

// VertexShader po
const VertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
in vec3 vPos;
in vec2 vUv;
out vec2 fUv;
void main() {
    fUv = vUv;
    gl_Position = projection * view * model * vec4(vPos, 1);
}
` + "\x00"

// FragmentShader po
const FragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fUv;
out vec4 outColor;
void main() {
    outColor = texture(tex, fUv);
}
` + "\x00"
