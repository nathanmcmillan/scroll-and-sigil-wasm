#version 300 es
uniform mat4 u_mvp;
layout(location = 0) in vec2 a_position;
layout(location = 1) in vec3 a_color;
out vec3 v_color;
void main() {
  v_color = a_color;
  gl_Position = u_mvp * vec4(a_position, 0.0, 1.0);
}
===========================================================
#version 300 es
precision mediump float;         
in vec3 v_color;
out vec4 pixel;
void main() {
  pixel = vec4(v_color, 1.0);
}
