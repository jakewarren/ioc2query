# Project Context

## Purpose
`ioc2query` is a Golang-based CLI utility that ingests raw text (from stdin or file), extracts key indicators of compromise (IOCs) — MD5/SHA1/SHA256 hashes, domain names, and IPv4 addresses — and transforms them into vendor-specific detection queries.

**Goals:**
- Provide security analysts with a simple CLI to convert IOCs into detection queries
- Support SentinelOne S1QLv2 and Rapid7 InsightIDR query formats
- Deliver reliable, reproducible integration tests to validate end-to-end output

**Target Users:**
- Incident Responders needing rapid transformation of indicators into detection queries
- Threat Hunters profiling malicious IOC patterns across endpoints via vendor query language

## Tech Stack
- **Language:** Go (Golang)
- **CLI Framework:** cobra or standard library flag package
- **IOC Extraction:** `github.com/vertoforce/go-ioc` - robust library for finding/extracting IOCs from text
- **Testing:** Go standard library testing package with table-driven tests
- **Linting:** golangci-lint
- **CI/CD:** GitHub Actions

## Project Conventions

### Code Style
- Follow standard Go conventions (gofmt, golint)
- Use structured logging with JSON output for machine readability
- Clear error messages with proper exit codes (0 on success, non-zero on errors)
- Idiomatic Go patterns: error handling with explicit returns, interfaces where appropriate
- Godoc comments for all exported functions, types, and packages

### Architecture Patterns
- CLI-first design: stdin/stdout by default with optional file I/O
- Modular extraction: separate IOC extraction logic from query generation
- Backend abstraction: clean interface for supporting multiple vendor backends
- Deduplication: always deduplicate extracted IOCs before processing

### Testing Strategy
- **Unit Tests:** Table-driven tests for IOC extraction, deduplication, CLI flag logic
- **Integration Tests:** End-to-end CLI validation with known IOC files and expected query outputs
- **Test Coverage:** Aim for 100% extraction accuracy for supported IOC types, use `go test -cover`
- **CI Pipeline:** GitHub Actions workflow to run tests on push and PR
- Test files should include realistic examples for each backend

### Git Workflow
- Feature branches for new development
- Pull requests required for merging to main
- CI must pass before merge
- Semantic commit messages preferred

## Domain Context
**Indicators of Compromise (IOCs):**
- MD5, SHA1, SHA256 hashes (file integrity indicators)
- Domain names (C2 infrastructure, malicious sites)
- IPv4 addresses (network indicators)
- Future consideration: IPv6, URLs, YARA rules (out of scope for v1)

**Supported Backends (v1):**
- **SentinelOne S1QLv2:** Query language for SentinelOne endpoint detection
- **Rapid7 InsightIDR:** Query syntax for InsightIDR SIEM platform

**Query Generation:**
- Custom backend implementations for each vendor query language
- Template-based query construction
- Extensible design for adding new backends

## Important Constraints
- v1 scope limited to MD5/SHA1/SHA256 hashes, domains, and IPv4 addresses only
- CLI-focused: no GUI or web interface in v1
- Single query output per run (or configurable individual queries per IOC)
- Cross-platform compatibility (Linux, macOS, Windows)

## External Dependencies
- **github.com/vertoforce/go-ioc:** Third-party library for extracting IOCs (MD5, SHA1, SHA256, domains, IPv4, IPv6, URLs, etc.)
  - Actively maintained with 8 stars, last updated Feb 2022
  - Supports defanging/fanging of IOCs
  - Regex-based extraction with comprehensive test coverage
  - MIT licensed
- **Go standard library:** Primary dependency for core functionality
- **Backend-specific libraries:** May require vendor-specific SDK or API clients for SentinelOne and Rapid7
- **Distribution:** Compiled binaries for multiple platforms, optionally via Homebrew or go install
