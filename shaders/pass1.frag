#version 460 core

in vec2 TexCoord;

out vec4 fragColor;

uniform sampler2D tex;

void main() {
  fragColor = texture(tex, TexCoord);
}