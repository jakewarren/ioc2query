## 1. WASM Build Infrastructure
- [x] 1.1 Verify Go code is WASM-compatible (no os.Exit in library code, handle stdin/stdout gracefully)
- [x] 1.2 Create Makefile with `wasm` target for compilation
- [x] 1.3 Update Makefile to copy wasm_exec.js from Go toolchain to web directory
- [ ] 1.4 Test WASM binary loads and runs in browser console

## 2. Web Interface - HTML Structure
- [x] 2.1 Create `web/index.html` with semantic structure
- [x] 2.2 Add page header with title and theme toggle button
- [x] 2.3 Create input textarea for IOC data entry
- [x] 2.4 Add controls grid: backend selector, separate queries checkbox, verbose checkbox
- [x] 2.5 Create output section with pre.terminal for displaying results
- [x] 2.6 Include wasm_exec.js and main.js script tags

## 3. Web Interface - Styling
- [x] 3.1 Create `web/styles.css` with CSS custom properties for theming
- [x] 3.2 Define light/dark theme color schemes with ANSI color mappings
- [x] 3.3 Style input controls (textarea, buttons, selects) with modern design
- [x] 3.4 Style output terminal with monospace font and ANSI support
- [x] 3.5 Add responsive layout with mobile breakpoints
- [x] 3.6 Implement card-based design matching cvrf-review aesthetic

## 4. Web Interface - JavaScript Logic
- [x] 4.1 Create `web/main.js` with runCommand() function to execute WASM
- [x] 4.2 Implement Go WASM runtime initialization with proper stdout/stderr capture
- [x] 4.3 Add ansiToHtml() function to convert ANSI escape codes to styled HTML
- [x] 4.4 Implement theme toggle with localStorage persistence
- [x] 4.5 Wire up Run button to collect inputs and call WASM with flags
- [x] 4.6 Add input validation and error display
- [x] 4.7 Handle WASM loading errors gracefully

## 5. Build Automation
- [x] 5.1 Create `.github/workflows/pages.yml` for GitHub Pages deployment
- [x] 5.2 Add workflow step to build WASM using `make wasm`
- [x] 5.4 Configure GitHub Pages to deploy from workflow artifact
- [x] 5.5 Add cache-busting query parameters for WASM file

## 6. Documentation
- [x] 6.1 Update README.md with "Web Interface" section
- [x] 6.2 Add local development instructions (how to build and serve)
- [x] 6.3 Add link to hosted GitHub Pages instance
- [x] 6.4 Document browser compatibility requirements

## 7. Testing and Validation
- [ ] 7.1 Test in Chrome, Firefox, and Safari
- [ ] 7.2 Verify WASM output matches CLI output for same inputs
- [ ] 7.3 Test all flag combinations (backend, separate, verbose)
- [ ] 7.4 Verify theme toggle works and persists
- [ ] 7.5 Test with empty input and no IOCs found scenarios
- [ ] 7.6 Validate mobile responsive design
- [ ] 7.7 Test ANSI color rendering in output

## 8. Deployment
- [ ] 8.1 Enable GitHub Pages in repository settings
- [ ] 8.2 Trigger Pages workflow and verify successful build
- [ ] 8.3 Verify hosted site is accessible
- [ ] 8.4 Update repository description with Pages URL
