#version 300 es
uniform mat4 u_mvp;
layout (location = 0) in vec2 a_position;
out vec2 v_texture;
void main() {
  vec4 position = u_mvp * vec4(a_position, 0.0, 1.0);
  v_texture = position.xy * 0.5 + 0.5;
  gl_Position = position;
}
===========================================================
#version 300 es
precision mediump float;
uniform sampler2D u_texture0;
in vec2 v_texture;
layout (location = 0) out vec4 color;
void main() {
  color = texture(u_texture0, v_texture);
}
