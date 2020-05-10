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
uniform sampler2D u_texture1;
in vec2 v_texture;
layout (location = 0) out vec4 color;
uniform mat4 inverse_projection;
uniform mat4 current_to_previous_matrix;
void main() {
  float z = texture(u_texture1, v_texture).r * 2.0 - 1.0;
  vec4 current = inverse_projection * vec4(v_texture * 2.0 - 1.0, z, 1.0);
  vec4 previous = current_to_previous_matrix * current;
  previous.xyz /= previous.w;
  previous.xy = previous.xy * 0.5 + 0.5;
  vec2 blur_vec = previous.xy - v_texture;
  const int num_samples = 4;
  vec3 sum = texture(u_texture0, v_texture).rgb;
  for (int i = 1; i < num_samples; ++i) {
    vec2 offset = blur_vec * (float(i) / float(num_samples - 1) - 0.5);
    sum += texture(u_texture0, v_texture + offset).rgb;
  }
  sum /= float(num_samples);
  color = vec4(sum, 1.0);
  // TODO convert unsigned int to float for frameGeo
}
