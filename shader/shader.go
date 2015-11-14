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
uniform mat4 mProj;
uniform mat4 mView;
uniform mat4 mModel;
uniform mat4 mNormal; // trans(inv(MV))
in vec3 vPos;
in vec3 vNor;
in vec2 vUv;
out vec2 fUv;
out vec3 LI;

struct Light {
	vec4 pos;
	vec3 La; // ambient
	vec3 Ld; // diffuse
	vec3 Ls; // specular
};
struct Material {
	vec3 Ka;
	vec3 Kd;
	vec3 Ks;
	float sh;
};
uniform Light light;
uniform Material mat;

void main() {
		vec4 tnorm = normalize(mNormal * vec4(vNor, 1.0));
		mat4 MV = mView * mModel;
		vec4 eye = MV * vec4(vPos ,1.0);
		vec3 s = normalize(vec3(light.pos - eye));
		vec3 v = normalize(-eye.xyz);
		vec3 r = reflect(-s, tnorm.xyz);

		vec3 amb = light.La * mat.Ka;

		float sn = max(dot(s, tnorm.xyz), 0.0);
		vec3 diff = light.Ld * mat.Kd * sn;

		vec3 spec = vec3(1.0);
		if(sn > 0.0){
			spec = light.Ls * mat.Ks * pow(max(dot(r, v) ,0.0), mat.sh);
		}

		LI = amb + diff + spec;
    fUv = vUv;
    gl_Position = mProj * MV * vec4(vPos, 1);
}
` + "\x00"

// FragmentShader po
const FragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fUv;
in vec3 LI;
out vec4 outColor;
void main() {
    outColor = vec4(LI, 1.0) * texture(tex, fUv);
}
` + "\x00"
