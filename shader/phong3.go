package shader

// Phong Shader
const VSPhong3 = `
#version 330
uniform mat4 mProj;
uniform mat4 mView;
uniform mat4 mModel;
uniform mat4 mNormal; // trans(inv(MV))
in vec3 vPos;
in vec3 vNor;
in vec4 vCol;
in vec2 vUv;
out vec2 fUv;
out vec4 fCol;
out vec3 fPos;
out vec3 fNor;
void main() {
	fNor = vec3(normalize(mNormal * vec4(vNor, 1.0)));
	mat4 MV = mView * mModel;
	fPos = vec3(MV * vec4(vPos, 1.0));
    fCol = vCol;
    fUv = vUv;
    gl_Position = mProj * MV * vec4(vPos, 1);
}
` + "\x00"

// PSPhong Phong Shader
const PSPhong3 = `
#version 330
uniform sampler2D tex;
uniform vec3 LP;
uniform vec3 LI;
uniform vec3 Ka;
uniform vec3 Kd;
uniform vec3 Ks;
uniform float Sh;
in vec4 fCol;
in vec2 fUv;
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
		vec4 tc = texture(tex, fUv);
        //if(tc.a < 0.01 || fCol.a < 0.01) discard;
        //outColor = vec4(ads(), 1.0) * fCol;
        //outColor = vec4(ads(), 1.0) * fCol * tc;
        outColor = tc;
        //outColor = vec4(fUv.x, fUv.y, 0.0, 1.0);
}
` + "\x00"
