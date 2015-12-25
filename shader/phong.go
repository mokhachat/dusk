package shader

// Phong Shader
const VSPhong = `
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

// PSPhong Phong Shader
const PSPhong = `
#version 330
uniform sampler2D tex;
uniform vec3 LP;
uniform vec3 LI;
uniform vec3 Ka;
uniform vec3 Kd;
uniform vec3 Ks;
uniform float Sh;
in vec2 fUv;
in vec3 fPos;
in vec3 fNor;
out vec4 outColor;
vec3 ads() {
	vec3 n = normalize(fNor);
	vec3 s = normalize(vec3(LP) - fPos);
	vec3 v = normalize(-fPos);
	vec3 h = normalize(v + s);
	// vec3 r = reflect(-s, n);
	// dot(r, v) -> dot(h, n)
	return LI * (Ka + Kd * max(dot(s, n), 0.0) + Ks * pow(max(dot(h, n), 0.0), Sh) );
}
void main() {
		outColor = vec4(ads(), 1.0) * texture(tex, fUv);
}
` + "\x00"
