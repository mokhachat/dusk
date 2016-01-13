package shader

// Phong Shader
const VSPhong2 = `
#version 330
uniform mat4 mProj;
uniform mat4 mView;
uniform mat4 mModel;
uniform mat4 mNormal; // trans(inv(MV))
in vec3 vPos;
in vec3 vNor;
in vec4 vCol;
out vec4 fCol;
out vec3 fPos;
out vec3 fNor;
void main() {
		fNor = vec3(normalize(mNormal * vec4(vNor, 1.0)));
		mat4 MV = mView * mModel;
		fPos = vec3(MV * vec4(vPos, 1.0));
    fCol = vCol;
    gl_Position = mProj * MV * vec4(vPos, 1);
}
` + "\x00"

// PSPhong Phong Shader
const PSPhong2 = `
#version 330
uniform vec3 LP;
uniform vec3 LI;
uniform vec3 Ka;
uniform vec3 Kd;
uniform vec3 Ks;
uniform float Sh;
in vec4 fCol;
in vec3 fPos;
in vec3 fNor;
out vec4 outColor;
vec3 ads() {
	vec3 n = normalize(fNor);
	vec3 s = normalize(vec3(LP) - fPos);
	vec3 v = normalize(-fPos);
	vec3 h = normalize(v + s);
	return LI * (Ka + Kd * max(dot(s, n), 0.0) + Ks * pow(max(dot(h, n), 0.0), Sh) );
}
void main() {
		if(fCol.a == 0.0) discard;
        outColor = vec4(ads(), 1.0) * fCol;
}
` + "\x00"
