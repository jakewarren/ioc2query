// FS Polyfill inside worker (must be defined before new Go())
self.files = {};
const fdPositions = {};
self.fs = {
    constants: { O_WRONLY: -1, O_RDONLY: 0, O_RDWR: 1, O_CREAT: 2, O_TRUNC: 4, O_APPEND: 8, O_EXCL: 16 },
    writeSync(fd, buf) {
        if (fd === 1 || fd === 2) {
            const str = new TextDecoder("utf-8").decode(buf);
            self.postMessage({ type: 'stdout', data: str });
            return buf.length;
        }
        return 0;
    },
    write(fd, buf, offset, length, position, callback) {
        if (fd === 1 || fd === 2) {
            const out = buf.slice(offset, offset + length);
            this.writeSync(fd, out);
            callback(null, length);
            return;
        }
        callback(new Error("not implemented"));
    },
    open(path, flags, mode, callback) {
        if (path === "/input.txt") {
            fdPositions[3] = 0;
            callback(null, 3);
            return;
        }
        const err = new Error("file not found");
        err.code = "ENOENT";
        callback(err);
    },
    read(fd, buffer, offset, length, position, callback) {
        if (fd === 3) {
            const fileContent = self.files["/input.txt"] || new Uint8Array(0);
            const explicitPosition = position !== null && position !== undefined;
            const currentPosition = explicitPosition ? position : (fdPositions[fd] || 0);
            const bytesRemaining = fileContent.length - currentPosition;
            if (bytesRemaining <= 0) {
                callback(null, 0);
                return;
            }
            const bytesToRead = Math.min(length, bytesRemaining);
            buffer.set(fileContent.slice(currentPosition, currentPosition + bytesToRead), offset);
            if (!explicitPosition) {
                fdPositions[fd] = currentPosition + bytesToRead;
            }
            callback(null, bytesToRead);
            return;
        }
        callback(new Error("invalid fd"));
    },
    fstat(fd, callback) {
        if (fd === 3) {
            const fileContent = self.files["/input.txt"] || new Uint8Array(0);
            callback(null, {
                size: fileContent.length,
                mode: 0o644,
                isDirectory: () => false,
                isCharacterDevice: () => false,
            });
            return;
        }
        callback(new Error("invalid fd"));
    },
    stat(path, callback) {
        if (path === "/input.txt") {
            const fileContent = self.files["/input.txt"] || new Uint8Array(0);
            callback(null, {
                size: fileContent.length,
                mode: 0o644,
                isDirectory: () => false,
                isCharacterDevice: () => false,
            });
            return;
        }
        const err = new Error("file not found");
        err.code = "ENOENT";
        callback(err);
    },
    close(fd, callback) {
        delete fdPositions[fd];
        callback(null);
    },
    fsync(fd, callback) { callback(null); },
    lstat(path, callback) { this.stat(path, callback); },
    mkdir(path, mode, callback) { callback(new Error("not implemented")); },
    readdir(path, callback) { callback(new Error("not implemented")); },
    unlink(path, callback) { callback(new Error("not implemented")); },
    rmdir(path, callback) { callback(new Error("not implemented")); },
    rename(from, to, callback) { callback(new Error("not implemented")); },
};

importScripts('wasm_exec.js');

let wasmModule;

async function initWasm() {
    try {
        const response = await fetch(`main.wasm?v=${Date.now()}`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const buffer = await response.arrayBuffer();
        const result = await WebAssembly.instantiate(buffer, new Go().importObject);
        wasmModule = result.module;
        self.postMessage({ type: 'ready' });
    } catch (err) {
        console.error("Worker init error:", err);
        self.postMessage({ type: 'error', data: `Failed to initialize: ${err.message}` });
    }
}

async function runWasm(data) {
    if (!wasmModule) {
        self.postMessage({ type: 'error', data: 'WASM module not initialized' });
        return;
    }

    const { args, input } = data;
    self.files["/input.txt"] = new TextEncoder().encode(input);

    // Create a new Go instance for each run - the Go runtime can only run once
    const go = new Go();
    go.argv = args;

    try {
        const instance = await WebAssembly.instantiate(wasmModule, go.importObject);
        await go.run(instance);
        self.postMessage({ type: 'done' });
    } catch (err) {
        console.error("Worker run error:", err);
        self.postMessage({ type: 'error', data: err.message });
    }
}

self.onmessage = async (e) => {
    const { type, data } = e.data;
    if (type === 'init') {
        await initWasm();
    } else if (type === 'run') {
        await runWasm(data);
    }
};
