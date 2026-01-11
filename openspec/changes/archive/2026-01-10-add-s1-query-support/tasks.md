# Implementation Tasks: Add S1 Query Support

## Pre-Implementation
- [x] Review proposal.md and design.md
- [x] Set up Go module structure with go.mod
- [x] Add `github.com/vertoforce/go-ioc` dependency

## IOC Extraction Module
- [x] Create `pkg/extraction/` package
- [x] Implement IOC extraction functions for each type:
  - [x] MD5 hash extraction
  - [x] SHA1 hash extraction
  - [x] SHA256 hash extraction
  - [x] Domain name extraction
  - [x] IPv4 address extraction
- [x] Implement deduplication logic
- [x] Write unit tests for extraction (table-driven)
- [x] Handle edge cases (malformed input, empty input)

## S1QL Backend Module
- [x] Create `pkg/backends/s1ql/` package
- [x] Define backend interface in `pkg/backends/interface.go`
- [x] Implement S1QLv2 query generator:
  - [x] Hash-based queries (process and file events)
  - [x] Domain-based queries (network connections)
  - [x] IPv4-based queries (network events)
  - [x] Combined query logic (OR conditions)
- [x] Implement query formatting and escaping
- [x] Write unit tests for query generation
- [x] Create test fixtures with expected S1QL output

## CLI Interface
- [x] Create `cmd/ioc2query/` package structure
- [x] Implement CLI with cobra or flag package:
  - [x] Input handling (stdin and file options)
  - [x] Output formatting options
  - [x] Backend selection flag (--backend s1)
  - [x] Help and usage documentation
- [x] Implement error handling and exit codes
- [x] Add verbose/debug logging options
- [x] Write CLI unit tests

## Integration Testing
- [x] Create `test/integration/` directory
- [x] Create sample IOC input files:
  - [x] Mixed IOC types
  - [x] Hash-only samples
  - [x] Network indicator samples
  - [x] Edge case samples
- [x] Create expected S1QL query outputs
- [x] Implement end-to-end integration tests
- [x] Verify test coverage with `go test -cover`

## Documentation
- [x] Write README.md with:
  - [x] Installation instructions
  - [x] Usage examples
  - [x] Supported IOC types
  - [x] Query output examples
- [x] Add Godoc comments to all exported functions
- [x] Document S1QL query patterns used

## Build & Distribution
- [x] Configure golangci-lint
- [ ] Set up GitHub Actions CI pipeline: (out of scope for initial implementation)
  - [ ] Run tests on push/PR
  - [ ] Run linting
  - [ ] Build for multiple platforms
- [ ] Create release build scripts (out of scope for initial implementation)
- [ ] Test cross-platform compilation (out of scope for initial implementation)

## Final Validation
- [x] Run full test suite (`go test ./...`)
- [x] Run linter (`golangci-lint run`)
- [x] Manual CLI testing with sample data
- [x] Validate all success criteria from proposal.md
- [x] Update all task checkboxes to [x]
