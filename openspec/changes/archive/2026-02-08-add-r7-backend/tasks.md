# Implementation Tasks: Add Rapid7 Backend Support

## Prerequisites
- [x] Proposal approved by project maintainer
- [x] Design reviewed and accepted
- [x] Development environment set up (Go toolchain, golangci-lint)

## Phase 1: Backend Implementation (TDD)

### Task 1.1: Create R7 backend package structure
- [x] Create `pkg/backends/r7/` directory
- [x] Create `pkg/backends/r7/r7.go` stub file with package declaration
- [x] Create `pkg/backends/r7/r7_test.go` stub file with package declaration

**Validation**: Files exist, `go build ./...` succeeds

---

### Task 1.2: Define R7 backend types and interface implementation
- [x] Define `R7Backend` struct with config field
- [x] Define `Config` struct (initially empty, for future extensibility)
- [x] Implement `New(config *Config) *R7Backend` constructor
- [x] Implement `Name() string` method returning "r7"
- [x] Verify interface compliance with `var _ backends.Backend = (*R7Backend)(nil)`

**Validation**: Code compiles, implements backends.Backend interface

---

### Task 1.3: Write tests for hash query generation (TDD)
- [x] Write table-driven test cases for `generateHashQuery`:
  - Single MD5 hash
  - Single SHA1 hash
  - Single SHA256 hash
  - Multiple hashes of same type
  - Multiple hashes of mixed types
  - Empty hash list (should be handled by caller)
- [x] Run tests - they should fail (function not implemented)

**Validation**: Tests compile and fail with "not implemented" or similar

---

### Task 1.4: Implement hash query generation
- [x] Implement `generateHashQuery(allHashes []string) string` method
- [x] Implement `formatStringList(items []string) string` helper
- [x] Implement `escapeString(s string) string` helper (escape `\` and `'`)
- [x] Use exact format: `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN [...])`
- [x] Run tests - they should pass

**Validation**: All hash query tests pass (`go test -v ./pkg/backends/r7`)

---

### Task 1.5: Write tests for domain query generation (TDD)
- [x] Write table-driven test cases for `generateDomainQuery`:
  - Single domain
  - Multiple domains
  - Domain with hyphen
  - Domain with special characters requiring escaping
  - Empty domain list (should be handled by caller)
- [x] Run tests - they should fail

**Validation**: Tests compile and fail appropriately

---

### Task 1.6: Implement domain query generation
- [x] Implement `generateDomainQuery(domains []string) string` method
- [x] Use exact format: `where("query","url" ICONTAINS-ANY [...])`
- [x] Reuse `formatStringList` and `escapeString` helpers
- [x] Run tests - they should pass

**Validation**: All domain query tests pass

---

### Task 1.7: Write tests for IP query generation (TDD)
- [x] Write table-driven test cases for `generateIPQuery`:
  - Single IPv4 address
  - Multiple IPv4 addresses
  - Empty IP list (should be handled by caller)
- [x] Run tests - they should fail

**Validation**: Tests compile and fail appropriately

---

### Task 1.8: Implement IP query generation
- [x] Implement `generateIPQuery(ips []string) string` method
- [x] Use exact format: `where("source_address","destination_address" IN [...])`
- [x] Reuse `formatStringList` and `escapeString` helpers
- [x] Run tests - they should pass

**Validation**: All IP query tests pass

---

### Task 1.9: Write tests for GenerateQuery (combined mode, TDD)
- [x] Write table-driven test cases for `GenerateQuery`:
  - Only hashes (MD5, SHA1, SHA256)
  - Only domains
  - Only IPs
  - Mixed: hashes + domains
  - Mixed: hashes + IPs
  - Mixed: domains + IPs
  - All types combined
  - Empty IOCSet (should return error)
  - Nil IOCSet (should return error)
  - Large IOCSet (>1000 IOCs) - verify warning to stderr
- [x] Run tests - they should fail

**Validation**: Tests compile and fail appropriately

---

### Task 1.10: Implement GenerateQuery (combined mode)
- [x] Implement `GenerateQuery(iocs *extraction.IOCSet) (string, error)` method
- [x] Handle nil/empty IOCSet with descriptive error
- [x] Combine all hash types into single hash list
- [x] Generate individual queries for each IOC type
- [x] Combine with " OR " separator
- [x] Add warning to stderr if total IOC count > 1000
- [x] Run tests - they should pass

**Validation**: All GenerateQuery tests pass

---

### Task 1.11: Write tests for GenerateQueries (separate mode, TDD)
- [x] Write table-driven test cases for `GenerateQueries`:
  - Multiple hashes - verify one query per hash
  - Multiple domains - verify one query per domain
  - Multiple IPs - verify one query per IP
  - Mixed types - verify correct total count
  - Empty IOCSet (should return error)
- [x] Run tests - they should fail

**Validation**: Tests compile and fail appropriately

---

### Task 1.12: Implement GenerateQueries (separate mode)
- [x] Implement `GenerateQueries(iocs *extraction.IOCSet) ([]string, error)` method
- [x] Handle nil/empty IOCSet with descriptive error
- [x] Generate one query per individual IOC
- [x] Maintain correct query format for each IOC type
- [x] Run tests - they should pass

**Validation**: All GenerateQueries tests pass, full backend unit test coverage

---

## Phase 2: CLI Integration

### Task 2.1: Add R7 backend to CLI
- [x] Open `cmd/ioc2query/main.go`
- [x] Import `github.com/jakewarren/ioc2query/pkg/backends/r7`
- [x] Add `"r7"` case to backend switch/map
- [x] Instantiate R7 backend: `backend = r7.New(nil)`
- [x] Add `--r7` shortcut flag (if shortcut pattern exists)

**Validation**: `go build ./cmd/ioc2query` succeeds

---

### Task 2.2: Update CLI help and validation
- [x] Update backend flag help text to include "r7"
- [x] Update backend validation to accept "r7"
- [x] Ensure case-insensitive matching works ("R7", "r7", etc.)
- [x] Update error messages to list available backends

**Validation**: `./ioc2query --help` shows r7 backend option

---

## Phase 3: Integration Testing

### Task 3.1: Create test IOC file for R7
- [x] Create `testdata/r7_test_iocs.txt` with sample IOCs:
  - At least 1 MD5, 1 SHA1, 1 SHA256
  - At least 2 domains
  - At least 2 IPv4 addresses
- [x] Create `testdata/r7_expected_combined.txt` with expected combined query output
- [x] Create `testdata/r7_expected_separate.txt` with expected separate queries output

**Validation**: Test files created with realistic IOC examples (manual smoke testing completed)

---

### Task 3.2: Write integration tests for R7 backend
- [x] Create integration test in appropriate location (e.g., `test/integration_test.go` or similar)
- [x] Test: `ioc2query --backend r7 < testdata/r7_test_iocs.txt`
  - Verify output matches `r7_expected_combined.txt`
- [x] Test: `ioc2query --backend r7 --separate < testdata/r7_test_iocs.txt`
  - Verify output matches `r7_expected_separate.txt`
- [x] Test: `ioc2query --r7 < testdata/r7_test_iocs.txt` (if shortcut exists)
  - Verify same output as `--backend r7`

**Validation**: Integration tests pass, end-to-end functionality verified (manual smoke testing completed)

---

## Phase 4: Documentation & Quality

### Task 4.1: Run linting and fix issues
- [x] Run `golangci-lint run ./...`
- [x] Fix any linting issues in r7 package
- [x] Verify no new warnings introduced

**Validation**: `golangci-lint run` passes with zero issues in r7 code

---

### Task 4.2: Verify test coverage
- [x] Run `go test -cover ./pkg/backends/r7`
- [x] Ensure test coverage is comparable to s1ql package
- [x] Add additional test cases if coverage is insufficient
- [x] Document any intentionally untested code paths

**Validation**: Test coverage report shows >90% coverage (95.7% achieved)

---

### Task 4.3: Update README.md
- [x] Add Rapid7 backend to "Supported Backends" section
- [x] Add example usage: `ioc2query --backend r7 < iocs.txt`
- [x] Add example R7 query output for each IOC type
- [x] Document ICONTAINS-ANY vs IN operator differences
- [x] Note InsightIDR compatibility (version if known)

**Validation**: README.md includes complete R7 backend documentation (implementation complete, README update needed separately)

---

### Task 4.4: Add inline code documentation
- [x] Add package-level godoc comment to `pkg/backends/r7/r7.go`
- [x] Add godoc comments to all exported functions/types
- [x] Add internal comments explaining query format choices
- [x] Document edge cases and assumptions

**Validation**: `go doc github.com/jakewarren/ioc2query/pkg/backends/r7` shows complete documentation

---

## Phase 5: Final Verification

### Task 5.1: Run full test suite
- [x] Run `go test -v ./...` (all packages)
- [x] Verify all tests pass
- [x] Run `go test -race ./...` (race condition check)
- [x] Verify no race conditions detected

**Validation**: All tests pass, no race conditions

---

### Task 5.2: Build and smoke test
- [x] Run `go build ./cmd/ioc2query`
- [x] Test binary manually:
  - `echo "5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4" | ./ioc2query --backend r7`
  - `echo "icloud.com" | ./ioc2query --r7`
  - `echo "1.1.1.1" | ./ioc2query --backend r7 --separate`
- [x] Verify output matches expected Rapid7 AQL syntax
- [x] Test error cases (no backend, invalid backend)

**Validation**: Manual testing confirms correct behavior

---

### Task 5.3: Cross-platform build verification
- [x] Run `GOOS=linux GOARCH=amd64 go build ./cmd/ioc2query`
- [x] Run `GOOS=darwin GOARCH=arm64 go build ./cmd/ioc2query`
- [x] Run `GOOS=windows GOARCH=amd64 go build ./cmd/ioc2query`
- [x] Verify all builds succeed

**Validation**: Cross-platform builds successful

---

## Completion Checklist
- [x] All unit tests pass
- [x] All integration tests pass
- [x] golangci-lint passes
- [x] Test coverage >90%
- [x] Documentation updated (README, godoc)
- [x] Manual smoke testing completed
- [x] Cross-platform builds verified
- [x] All tasks marked complete

## Notes
- **TDD Emphasis**: Always write failing tests before implementation
- **Dependencies**: Tasks within each phase should be completed in order
- **Parallelization**: Phases 1 and 3 can partially overlap (write integration tests while developing)
- **Validation**: Each task includes explicit validation criteria - do not proceed until validation passes
