# web-interface Specification

## Purpose
TBD - created by archiving change add-wasm-web-interface. Update Purpose after archive.
## Requirements
### Requirement: System SHALL compile to WebAssembly
The system MUST build a WebAssembly binary from the existing CLI that runs in modern browsers.

**Acceptance Criteria:**
- Builds with `GOOS=js GOARCH=wasm` compilation
- Produces `main.wasm` artifact under 10MB
- Includes Go WASM runtime (`wasm_exec.js`)
- No modifications to core CLI logic required
- Works in Chrome 88+, Firefox 89+, Safari 14+

#### Scenario: Build WASM binary
**GIVEN** the ioc2query CLI codebase exists  
**WHEN** developer runs `GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/ioc2query`  
**THEN** a valid `main.wasm` file is created that can execute in browsers

#### Scenario: Load WASM in browser
**GIVEN** `main.wasm` and `wasm_exec.js` are served  
**WHEN** user loads the web page  
**THEN** the WASM module initializes without errors and is ready to execute

---

### Requirement: System SHALL provide web-based input interface
The system MUST provide a textarea for users to paste IOC data directly in the browser.

**Acceptance Criteria:**
- Multi-line textarea for IOC input
- Accepts minimum 10KB of text
- Placeholder text with example IOC format
- Clear visual focus state
- Responsive on mobile devices (min 320px width)

#### Scenario: Paste IOCs into textarea
**GIVEN** the web page is loaded  
**WHEN** user pastes "5d41402abc4b2a76b9719d911017c592" into the input textarea  
**THEN** the text is displayed correctly and ready for processing

#### Scenario: Handle large input
**GIVEN** user has 500 lines of IOC data  
**WHEN** user pastes the data into textarea  
**THEN** all text is accepted without truncation or browser hang

---

### Requirement: System SHALL provide control options
The system MUST expose CLI flags through web UI controls for backend selection, separate queries mode, and verbose output.

**Acceptance Criteria:**
- Backend selector (dropdown or radio buttons) showing "S1QL"
- Checkbox for "Generate separate queries per IOC type"
- Checkbox for "Verbose output (show extraction details)"
- Run button to trigger processing
- Clear button to reset form (optional v2)

#### Scenario: Select backend and run
**GIVEN** IOCs are entered in textarea  
**WHEN** user selects "S1QL" backend and clicks "Run"  
**THEN** the system generates S1QL queries and displays results

#### Scenario: Enable separate queries mode
**GIVEN** IOCs are entered in textarea  
**WHEN** user checks "Generate separate queries" and clicks "Run"  
**THEN** the system outputs individual queries per IOC type, not combined

#### Scenario: Enable verbose mode
**GIVEN** IOCs are entered in textarea  
**WHEN** user checks "Verbose output" and clicks "Run"  
**THEN** the system shows extraction statistics (MD5: X, SHA1: Y, etc.) before query output

---

### Requirement: System SHALL display query results
The system MUST display generated queries in a styled output area that preserves formatting and supports ANSI colors.

**Acceptance Criteria:**
- Output displayed in monospace font
- ANSI color codes converted to HTML/CSS styling
- Preserves whitespace and line breaks
- Output area has minimum height of 100px
- Scrollable when output exceeds viewport
- Clear visual distinction from input area

#### Scenario: Display simple query output
**GIVEN** user runs query generation  
**WHEN** WASM process returns `src.process.image.md5 IN ("abc123")`  
**THEN** output area displays the query with proper formatting

#### Scenario: Render ANSI colors
**GIVEN** verbose mode is enabled  
**WHEN** CLI outputs "Extracted: 5 MD5" with ANSI color codes  
**THEN** the output renders with styled colors matching the ANSI codes

#### Scenario: Handle long output
**GIVEN** user generates queries for 100+ IOCs  
**WHEN** output exceeds visible area  
**THEN** output area becomes scrollable and maintains readability

---

### Requirement: System SHALL support theme switching
The system MUST provide light and dark theme options that affect the entire UI including output terminal colors.

**Acceptance Criteria:**
- Defaults to system preference (prefers-color-scheme)
- Provides manual theme toggle button
- Persists user choice in localStorage
- Switches entire UI including input, controls, and output
- ANSI colors adapt to theme for readability

#### Scenario: Auto-detect dark mode
**GIVEN** user's system is set to dark mode  
**WHEN** user loads the page  
**THEN** the UI displays in dark theme with appropriate colors

#### Scenario: Toggle to light theme
**GIVEN** page is in dark theme  
**WHEN** user clicks the theme toggle button  
**THEN** UI switches to light theme and choice is saved to localStorage

#### Scenario: Remember theme preference
**GIVEN** user previously selected light theme  
**WHEN** user returns to page in a new session  
**THEN** light theme is applied automatically on page load

---

### Requirement: System SHALL execute CLI via WASM
The system MUST run the compiled Go CLI in the browser by constructing proper command-line arguments and capturing stdout/stderr.

**Acceptance Criteria:**
- Constructs argv array from UI controls
- Passes input text to Go program via stdin or arguments
- Captures stdout for query output
- Captures stderr for errors and verbose logging
- Returns output to JavaScript for display
- Handles execution errors gracefully

#### Scenario: Execute with basic flags
**GIVEN** user enters IOCs and selects S1QL backend  
**WHEN** user clicks "Run"  
**THEN** JavaScript invokes WASM with `['ioc2query', '--backend', 's1']` arguments and displays output

#### Scenario: Capture stdout and stderr separately
**GIVEN** verbose mode is enabled  
**WHEN** WASM executes  
**THEN** verbose logging (stderr) and query output (stdout) are both captured and displayed

#### Scenario: Handle WASM execution error
**GIVEN** WASM fails to initialize  
**WHEN** user clicks "Run"  
**THEN** an error message is displayed: "Failed to initialize WASM. Please refresh the page."

---

### Requirement: System SHALL deploy via GitHub Pages
The system MUST automatically build and deploy the web interface using GitHub Actions to GitHub Pages.

**Acceptance Criteria:**
- GitHub Actions workflow builds WASM on every push to main
- Workflow copies wasm_exec.js from Go toolchain
- Workflow deploys HTML/CSS/JS/WASM to GitHub Pages
- Site accessible at `https://jakewarren.github.io/ioc2query/` (or repo owner's domain)
- Build failures prevent deployment
- Cache-busting for WASM file to prevent stale versions

#### Scenario: Deploy on push to main
**GIVEN** changes are merged to main branch  
**WHEN** GitHub Actions workflow runs  
**THEN** WASM is built, files are deployed, and site is updated within 5 minutes

#### Scenario: Build failure prevents deployment
**GIVEN** WASM build fails due to code error  
**WHEN** workflow runs  
**THEN** deployment is skipped and workflow shows failed status

#### Scenario: Access hosted site
**GIVEN** deployment succeeded  
**WHEN** user navigates to GitHub Pages URL  
**THEN** the web interface loads and is fully functional

---

### Requirement: System SHALL validate user input
The system MUST provide client-side validation to prevent invalid operations and give clear feedback.

**Acceptance Criteria:**
- Disables "Run" button when textarea is empty
- Shows warning if no IOCs detected (optional: pre-scan)
- Validates backend is selected before running
- Displays loading state during WASM execution
- Prevents concurrent executions

#### Scenario: Disable run with empty input
**GIVEN** textarea is empty  
**WHEN** user attempts to click "Run"  
**THEN** button is disabled or shows error: "Please enter IOC data"

#### Scenario: Show loading state
**GIVEN** user clicks "Run"  
**WHEN** WASM is executing  
**THEN** button text changes to "Running..." and is disabled until completion

#### Scenario: Prevent concurrent runs
**GIVEN** WASM execution is in progress  
**WHEN** user clicks "Run" again  
**THEN** the click is ignored and no second execution starts

---

### Requirement: System SHALL match CLI output format
The web interface MUST produce identical output to the CLI for the same inputs and flags.

**Acceptance Criteria:**
- Query format matches CLI exactly (byte-for-byte when possible)
- Verbose output shows same statistics as CLI
- Error messages match CLI error messages
- Exit codes mapped to success/error states in UI
- No web-specific output variations

#### Scenario: Verify output parity
**GIVEN** CLI command `echo "abc123" | ioc2query --backend s1` produces output X  
**WHEN** same IOC and flags used in web interface  
**THEN** web interface displays identical output X

#### Scenario: Match verbose output
**GIVEN** CLI with `--verbose` shows "Extracted: 1 MD5"  
**WHEN** web interface runs with verbose enabled  
**THEN** same "Extracted: 1 MD5" message appears in output

---

### Requirement: System SHALL be mobile responsive
The web interface MUST be usable on mobile devices with touch-friendly controls and readable text.

**Acceptance Criteria:**
- Minimum viewport width: 320px
- Touch-friendly buttons (minimum 44x44px tap targets)
- Readable font sizes (minimum 14px)
- No horizontal scrolling required
- Stacks controls vertically on narrow screens
- Works in mobile Safari, Chrome, Firefox

#### Scenario: View on mobile phone
**GIVEN** user visits site on iPhone SE (375px width)  
**WHEN** page loads  
**THEN** all controls are visible, usable, and no horizontal scroll is present

#### Scenario: Use textarea on mobile
**GIVEN** user taps input textarea on mobile  
**WHEN** keyboard appears  
**THEN** textarea remains visible and usable without layout breaking

---

### Requirement: System SHALL handle errors gracefully
The system MUST catch and display user-friendly error messages for all failure scenarios.

**Acceptance Criteria:**
- WASM loading errors show actionable message
- Execution errors don't crash page
- Network errors (if WASM not loaded) are explained
- No IOCs found shows clear message: "No IOCs detected in input"
- JavaScript errors logged to console for debugging

#### Scenario: Handle no IOCs found
**GIVEN** user enters "hello world" with no IOCs  
**WHEN** user clicks "Run"  
**THEN** output shows: "No IOCs found in input"

#### Scenario: Handle WASM load failure
**GIVEN** main.wasm fails to download (network error)  
**WHEN** user clicks "Run"  
**THEN** error message shows: "Failed to load WASM. Check connection and refresh."

#### Scenario: Handle malformed input
**GIVEN** input contains special characters that might break processing  
**WHEN** WASM executes  
**THEN** error is caught and displayed without crashing the page

