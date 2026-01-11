## Context

The ioc2query CLI tool needs a web interface to enable browser-based usage without installation. This is particularly valuable for security analysts who need quick access during incident response. The design follows the proven pattern from jakewarren/cvrf-review which successfully combines CLI and WASM-based web interfaces.

**Constraints:**
- Must maintain CLI as primary interface - web is supplementary
- No backend server - fully client-side execution via WASM
- Single-page application with minimal dependencies
- Must work offline once page is loaded

**Stakeholders:**
- Security analysts needing quick IOC transformations
- Users without Go installation or CLI access
- Mobile users (responsive design required)

## Goals / Non-Goals

**Goals:**
- Compile existing CLI to WASM without code modifications
- Create intuitive web UI matching cvrf-review design patterns
- Support all CLI flags through web controls
- Deploy automatically via GitHub Pages
- Provide same output quality as CLI

**Non-Goals:**
- Not building a backend API or server
- Not storing user data or IOCs on server
- Not replacing or diminishing CLI experience
- Not adding features exclusive to web (feature parity)
- Not supporting file uploads in v1 (use paste/textarea only)

## Decisions

### Decision: Single-Page Application with WASM
Compile the existing Go CLI to WebAssembly and run it directly in the browser.

**Rationale:**
- Reuses 100% of existing extraction and query generation logic
- No backend infrastructure or hosting costs beyond static files
- Works offline after initial page load
- Proven pattern from cvrf-review reference implementation
- Maintains exact CLI behavior consistency

**Alternatives considered:**
- JavaScript rewrite → Rejected: duplicate code maintenance, behavior drift risk
- Backend API → Rejected: adds infrastructure, hosting costs, latency
- Electron app → Rejected: still requires download/install

### Decision: GitHub Pages Deployment
Use GitHub Pages with Actions workflow to build and deploy.

**Rationale:**
- Zero-cost hosting for public repos
- Automatic SSL/HTTPS
- Simple workflow integration
- Same pattern as cvrf-review

**Alternatives considered:**
- Netlify/Vercel → Rejected: unnecessary complexity for static site
- Self-hosted → Rejected: maintenance burden

### Decision: Match cvrf-review UI Design
Follow cvrf-review's design system including theme support and layout patterns.

**Rationale:**
- Proven, polished design
- Same author/project family
- Reduces design decisions
- Familiar to users of other jakewarren tools

**Alternatives considered:**
- Custom design → Rejected: more work, potentially inconsistent
- Minimal/unstyled → Rejected: poor user experience

### Decision: No Backend Processing
All IOC extraction and query generation happens client-side in WASM.

**Rationale:**
- IOCs may be sensitive security data - keep local
- Simpler architecture and deployment
- Better privacy story
- Lower latency (no network round trip)

**Alternatives considered:**
- Server-side processing → Rejected: security/privacy concerns, infrastructure cost

## Architecture

### Component Overview

```
┌─────────────────────────────────────────┐
│           Browser (Client)              │
│  ┌───────────────────────────────────┐  │
│  │         index.html                │  │
│  │  ┌────────────┐   ┌────────────┐ │  │
│  │  │   Input    │   │  Controls  │ │  │
│  │  │  Textarea  │   │   (flags)  │ │  │
│  │  └────────────┘   └────────────┘ │  │
│  │  ┌────────────────────────────┐  │  │
│  │  │   Output Terminal          │  │  │
│  │  │   (ANSI-rendered results)  │  │  │
│  │  └────────────────────────────┘  │  │
│  └───────────────────────────────────┘  │
│              ↓ main.js                  │
│  ┌───────────────────────────────────┐  │
│  │      Go WASM Runtime              │  │
│  │  (wasm_exec.js + main.wasm)       │  │
│  │                                   │  │
│  │  ┌────────────┐  ┌─────────────┐ │  │
│  │  │ Extraction │→│  S1QL Gen   │ │  │
│  │  │   (IOCs)   │  │  (Queries)  │ │  │
│  │  └────────────┘  └─────────────┘ │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

### Data Flow

1. User pastes IOC text into textarea
2. User selects options (backend, separate, verbose)
3. User clicks "Run" button
4. `main.js` packages input as CLI arguments
5. Go WASM runtime executes with captured stdout/stderr
6. Output captured and returned to JavaScript
7. `ansiToHtml()` converts ANSI codes to styled HTML
8. Results displayed in output terminal

### File Structure

```
web/
├── index.html       # Main page structure
├── main.js          # WASM loader and UI logic
├── styles.css       # Theme and component styles
├── main.wasm        # Compiled Go binary (built by CI)
└── wasm_exec.js     # Go WASM runtime (copied from toolchain)

.github/workflows/
└── pages.yml        # Build and deploy workflow
```

## Technical Details

### WASM Compilation
```bash
GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/ioc2query
```

**No code changes needed** - existing CLI structure works as-is because:
- Go's `flag` package works in WASM
- `os.Stdin` can be mocked via Go runtime
- `os.Stdout`/`os.Stderr` captured by wasm_exec.js
- All pkg/ code is platform-independent

### Browser Compatibility
- **Minimum:** Chrome 88+, Firefox 89+, Safari 14+ (WebAssembly support)
- **WASM file size:** ~2-5MB (depends on dependencies)
- **Loading time:** 1-3 seconds on first load, cached thereafter

### ANSI Rendering
The `ansiToHtml()` function converts ANSI escape codes to styled spans:
- Maps codes 30-37, 90-97 to CSS custom properties
- Preserves table formatting with `white-space: pre`
- Links URLs automatically (pattern: fortiguard.fortinet.com → clickable)

### Theme System
```css
:root {
  --bg: light colors...
}
@media (prefers-color-scheme: dark) {
  :root { --bg: dark colors... }
}
body.theme-light { --bg: override... }
body.theme-dark { --bg: override... }
```

User preference stored in `localStorage` and takes precedence over system preference.

## Risks / Trade-offs

| Risk                             | Impact                  | Mitigation                                                  |
| -------------------------------- | ----------------------- | ----------------------------------------------------------- |
| Large WASM file size             | Slow initial load       | Use cache-busting + CDN, accept 2-5MB as normal for Go WASM |
| Browser compatibility            | Some users can't access | Document minimum versions, CLI remains primary              |
| WASM security sandbox limits     | No file system access   | Design doesn't need it - use textarea input only            |
| stdout/stderr capture complexity | Output truncation       | Test thoroughly, increase buffer size if needed             |
| ANSI rendering edge cases        | Malformed output        | Escape HTML before rendering, test with known outputs       |

**Key Trade-off:** Larger bundle size vs. zero maintenance code duplication  
→ **Decision:** Accept larger bundle. Reusing CLI logic worth the download cost.

## Migration Plan

### Phase 1: Build Infrastructure
1. Add WASM build to Makefile/scripts
2. Create GitHub Actions workflow
3. Test WASM compilation locally

### Phase 2: Web Interface
1. Create HTML/CSS/JS files matching cvrf-review structure
2. Implement WASM loader and capture logic
3. Wire up UI controls to CLI flags
4. Test locally with `python3 -m http.server`

### Phase 3: Deployment
1. Enable GitHub Pages in repo settings
2. Push to main → workflow builds and deploys
3. Verify hosted site works
4. Update README with link

### Rollback
If critical issues found post-deployment:
1. Disable GitHub Pages in settings (site goes offline)
2. Revert workflow changes
3. CLI unaffected - remains fully functional

No database migrations or API changes needed.

## Open Questions

1. **Should we support file upload in addition to textarea?**
   - Recommendation: Not in v1. Adds complexity. Users can paste file contents.
   
2. **Do we need to limit input size to prevent browser hangs?**
   - Recommendation: Yes. Add 1MB textarea limit with warning. Test with large inputs.
   
3. **Should output be downloadable?**
   - Recommendation: Nice to have but not v1. Users can copy/paste from textarea.

4. **Do we need analytics to see usage?**
   - Recommendation: No. Privacy-first approach. No tracking.
