import Veda from 'vedajs';

// Create canvas
const canvas = document.createElement('canvas');
canvas.style.position = 'fixed';
canvas.style.top = '0';
canvas.style.left = '0';
canvas.style.width = '100%';
canvas.style.height = '100%';
canvas.style.zIndex = '-1';
document.body.appendChild(canvas);

// Make everything transparent
const style = document.createElement('style');
style.innerHTML = `
body {
    background-color: transparent !important;
}
#workbench\.parts\.sidebar, #workbench\.main\.container {
    background-color: transparent !important;
}
#workbench\.parts\.activitybar, #workbench\.parts\.titlebar {
    opacity: 0.5 !important;
}
.part.sidebar, .monaco-editor .margin, .tab,
.monaco-workbench, .content, .monaco-editor, .monaco-editor-background, .editor-container, .title.tabs {
    background-color: transparent !important;
}
.minimap canvas {
    opacity: 0.5;
}
.part.titlebar, .part.activitybar {
    background-color: rgba(0, 0, 0, 0.5) !important;
}
.view-line span {
    background: black;
}
`;
document.body.appendChild(style);

// Init VEDA
const veda = new Veda({});
veda.setCanvas(canvas);
window.addEventListener('resize', () => {
    veda.resize(window.innerWidth, window.innerHeight);
});

const DEFAULT_CODE = `
precision mediump float;
uniform float time;
uniform vec2 mouse;
uniform vec2 resolution;
void main() {
    vec2 uv = gl_FragCoord.xy / resolution.xy;
    gl_FragColor = vec4(uv,0.5+0.5*sin(time),1.0);
}
`.trim();

veda.loadFragmentShader(DEFAULT_CODE);
veda.play();

// Bind keys
document.onkeydown = (e) => {
    if (e.ctrlKey && e.code === 'Enter') {
        const text: string = (document.querySelector('.monaco-editor-background') as any).innerText.trim();
        const decoded = decodeURIComponent(encodeURIComponent(text).replace(/%C2%A0/g, '%20'));

        veda.loadFragmentShader(decoded);
    }
};
