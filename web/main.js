// main.js

const inputData = document.getElementById('input-data');
const backendSelect = document.getElementById('backend-select');
const separateCheckbox = document.getElementById('separate-checkbox');
const verboseCheckbox = document.getElementById('verbose-checkbox');
const runButton = document.getElementById('run-button');
const copyButton = document.getElementById('copy-button');
const outputTerminal = document.getElementById('output-terminal');
const loading = document.getElementById('loading');
const themeToggle = document.getElementById('theme-toggle');
const sunIcon = document.getElementById('sun-icon');
const moonIcon = document.getElementById('moon-icon');

let terminalBuffer = "";
let lastUpdate = 0;
let worker;

// Theme Management
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'dark';
    if (savedTheme === 'light') {
        document.body.classList.remove('dark-theme');
        sunIcon.style.display = 'none';
        moonIcon.style.display = 'block';
    } else {
        document.body.classList.add('dark-theme');
        sunIcon.style.display = 'block';
        moonIcon.style.display = 'none';
    }
}

themeToggle.addEventListener('click', () => {
    document.body.classList.toggle('dark-theme');
    const isDark = document.body.classList.contains('dark-theme');
    localStorage.setItem('theme', isDark ? 'dark' : 'light');

    if (isDark) {
        sunIcon.style.display = 'block';
        moonIcon.style.display = 'none';
    } else {
        sunIcon.style.display = 'none';
        moonIcon.style.display = 'block';
    }
});

function updateTerminal(force = false) {
    const now = Date.now();
    if (force || now - lastUpdate > 100) {
        if (outputTerminal) {
            outputTerminal.innerHTML = ansiToHtml(terminalBuffer);
            outputTerminal.scrollTop = outputTerminal.scrollHeight;
        }
        lastUpdate = now;
    }
}

function ansiToHtml(text) {
    const colors = {
        0: 'initial',
        31: '#ff5555', // red
        32: '#50fa7b', // green
        33: '#f1fa8c', // yellow
        34: '#bd93f9', // blue
        35: '#ff79c6', // magenta
        36: '#8be9fd', // cyan
        37: '#f8f8f2'  // white
    };

    return text
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/\x1b\[(\d+)m/g, (match, code) => {
            if (code === '0') return '</span>';
            const color = colors[code];
            if (color) return `<span style="color: ${color}">`;
            return '';
        });
}

function initWorker() {
    loading.classList.remove('hidden');
    worker = new Worker('worker.js');

    worker.onmessage = (e) => {
        const { type, data } = e.data;

        if (type === 'ready') {
            loading.classList.add('hidden');
            runButton.disabled = false;
        } else if (type === 'stdout') {
            terminalBuffer += data;
            updateTerminal();
        } else if (type === 'done') {
            updateTerminal(true);
            loading.classList.add('hidden');
            runButton.disabled = false;
        } else if (type === 'error') {
            terminalBuffer += `\nError: ${data}\n`;
            updateTerminal(true);
            loading.classList.add('hidden');
            runButton.disabled = false;
        }
    };

    worker.postMessage({ type: 'init' });
}

runButton.addEventListener('click', () => {
    const input = inputData.value.trim();
    if (!input) {
        outputTerminal.textContent = "Error: Input is empty";
        return;
    }

    terminalBuffer = "";
    outputTerminal.textContent = "";
    loading.classList.remove('hidden');
    runButton.disabled = true;

    const args = ["ioc2query"];
    args.push("-backend", backendSelect.value);
    if (separateCheckbox.checked) args.push("-separate");
    if (verboseCheckbox.checked) args.push("-verbose");
    args.push("-input", "/input.txt");

    worker.postMessage({
        type: 'run',
        data: { args, input }
    });
});

copyButton.addEventListener('click', () => {
    const text = outputTerminal.innerText;
    navigator.clipboard.writeText(text).then(() => {
        const originalText = copyButton.textContent;
        copyButton.textContent = "Copied!";
        setTimeout(() => copyButton.textContent = originalText, 2000);
    });
});

initTheme();
initWorker();
