# Change: Add WASM-based Web Interface

## Why
Security analysts and threat hunters currently need to install and run ioc2query locally. A browser-based interface would enable immediate access without installation, particularly useful for quick IOC transformations during incident response when installing tools may not be feasible or convenient. This mirrors the successful pattern from cvrf-review which provides both CLI and web interfaces.

## What Changes
- Compile ioc2query to WebAssembly (WASM) for browser execution
- Create single-page web interface with input/output areas and controls
- Implement WASM loader and Go runtime integration (wasm_exec.js)
- Add theme support (light/dark) matching cvrf-review design patterns
- Support all core CLI flags via web UI controls (backend selection, separate queries, verbose mode)
- Deploy via GitHub Pages for hosted access

## Impact
- **Affected specs:** Creates new `web-interface` capability
- **Affected code:**
  - `cmd/ioc2query/main.go` - May need WASM-friendly adjustments (no changes to core logic expected)
  - New `web/` directory: `index.html`, `main.js`, `styles.css`, `docs/` for data files
  - `.github/workflows/` - Add GitHub Pages build workflow
- **Dependencies:** 
  - Go WASM toolchain (GOOS=js GOARCH=wasm)
  - Static web server for development
- **Compatibility:** No breaking changes to CLI interface
- **Testing:** Manual browser testing (Chrome, Firefox, Safari), validate WASM output matches CLI
