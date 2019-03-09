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
document.onkeydown = async (e) => {
    if (e.ctrlKey && e.code === 'Enter') {
        // Send file URI to index.js
        const uri = document.querySelector('.monaco-editor')!.getAttribute('data-uri')!;
        chrome.runtime.sendMessage({ fileUri: uri }, (res: string) => {
            console.log(`>> Sent fileUri: ${uri}`);
            console.log('>> Response: ', res);
        });
    }
};

window.chrome.runtime.onMessage.addListener((msg: any) => {
    if (!msg.shader) { return; }
    console.log(msg.shader);

    veda.loadFragmentShader(msg.shader);
});
