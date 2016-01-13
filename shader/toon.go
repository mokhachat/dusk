package shader

// VSToon Toon Shader
const VSToon = `
#version 330
uniform mat4 mProj;
uniform mat4 mView;
uniform mat4 mModel;
uniform mat4 mNormal; // trans(inv(MV))
in vec3 vPos;
in vec3 vNor;
in vec2 vUv;
out vec2 fUv;
out vec3 fPos;
out vec3 fNor;
void main() {
		fNor = vec3(normalize(mNormal * vec4(vNor, 1.0)));
		mat4 MV = mView * mModel;
		fPos = vec3(MV * vec4(vPos, 1.0));
    fUv = vUv;
    gl_Position = mProj * MV * vec4(vPos, 1);
}
` + "\x00"

// PSToon Toon Shader
const PSToon = `
#version 330
uniform sampler2D tex;
uniform vec3 LP;
uniform vec3 LI;
uniform vec3 Ka;
uniform vec3 Kd;
const int levels = 64;
const float sf = 1.0 / levels;
in vec2 fUv;
in vec3 fPos;
in vec3 fNor;
out vec4 outColor;
vec3 toon() {
	vec3 s = normalize(vec3(LP) - fPos);
	float cosine = max(0.0, dot(s, fNor));
	vec3 diffuse = Kd * floor(cosine * levels) * sf;
	return LI * (Ka + diffuse);
}
void main() {
		outColor = vec4(toon(), 1.0) * texture(tex, fUv);
}
` + "\x00"
