# Design: Add Rapid7 Backend Support

## Architecture Overview
The Rapid7 backend follows the same architectural pattern established by the S1QL backend, implementing the `backends.Backend` interface with vendor-specific query generation logic.

```
┌─────────────────┐
│   CLI Layer     │
│  (cmd/ioc2query)│
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  IOC Extraction │
│(pkg/extraction) │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────┐
│   Backend Interface         │
│  (pkg/backends.Backend)     │
└──────┬──────────────┬───────┘
       │              │
       ▼              ▼
┌────────────┐  ┌────────────┐
│ S1QL       │  │ R7         │
│ Backend    │  │ Backend    │
└────────────┘  └────────────┘
```

## Design Decisions

### 1. Backend Identifier: "r7"
**Decision**: Use "r7" as the backend identifier instead of "rapid7" or "insightidr"

**Rationale**:
- Brevity: Shorter CLI usage (`--backend r7`)
- Consistency: Matches pattern of "s1" for SentinelOne
- Clarity: Still clearly identifies Rapid7
- User preference: Confirmed via user feedback

**Alternatives Rejected**:
- `rapid7`: Longer to type, less consistent with "s1" pattern
- `insightidr`: Too product-specific, may confuse if moving to other Rapid7 products

### 2. Hash Query Format: Wildcard Pattern
**Decision**: Use `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['hash'])` for all hash types

**Rationale**:
- Simplicity: Single query pattern for all hash types (MD5/SHA1/SHA256)
- User specification: Exact format requested
- Rapid7 capability: InsightIDR supports wildcard field matching
- Deduplication: IOC extraction already deduplicates hashes by type, so no risk of duplicate matching

**Alternatives Rejected**:
- Type-specific fields: Would require detecting hash length and generating different queries
  - More complex implementation
  - No added value given wildcard support

**Query Structure**:
```
where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['hash1', 'hash2', ...])
```

### 3. Domain/URL Handling: ICONTAINS-ANY
**Decision**: Use `ICONTAINS-ANY` operator for domains with `query,url` fields

**Rationale**:
- Case-insensitive matching: Domains are case-insensitive by nature
- Partial matching: Allows matching subdomains (e.g., "evil.com" matches "sub.evil.com")
- User specification: Confirmed format from examples
- Fields: Both DNS query logs and URL access logs searched

**Query Structure**:
```
where("query","url" ICONTAINS-ANY ['domain1.com', 'domain2.org', ...])
```

**Trade-off**: ICONTAINS-ANY may produce more false positives than exact matching, but this is acceptable for IOC hunting where recall > precision.

### 4. IP Address Handling: IN Operator
**Decision**: Use `IN` operator for IPs with `source_address,destination_address` fields

**Rationale**:
- Exact matching: IP addresses should match exactly (no partial matching needed)
- Bidirectional: Search both source and destination to catch all traffic
- Efficiency: IN operator is more efficient than CONTAINS for exact matches

**Query Structure**:
```
where("source_address","destination_address" IN ['1.1.1.1', '2.2.2.2', ...])
```

### 5. Query Grouping Strategy
**Decision**: Return queries grouped by IOC type as []string, with presentation layer handling formatting

**Grouped Mode** (default `GenerateQuery`):
- Returns []string with one element per IOC type
- Hashes (all types combined) as one query element
- Domains as one query element
- IPs as one query element
- CLI/web layer joins with "\n\n" (double newline) for readability
- Each query can be copied individually into Rapid7 UI

**Separate Mode** (`--separate` flag, uses `GenerateQueries`):
- One query per individual IOC indicator
- Returns []string with one element per IOC
- Allows individual copy/paste into Rapid7 UI
- Useful for testing specific indicators

**Rationale**:
- Separation of concerns: Backend generates queries, CLI handles presentation
- Flexibility: Different UIs can format differently (CLI uses newlines, web could use different styling)
- Usability: Each IOC type query can be run separately in Rapid7 if needed
- Backend interface consistency: Both backends return []string

**Implementation**:
```go
func (b *R7Backend) GenerateQuery(iocs *extraction.IOCSet) ([]string, error) {
    var queries []string
    
    // Combine all hash types into single query
    var allHashes []string
    allHashes = append(allHashes, iocs.MD5Hashes...)
    allHashes = append(allHashes, iocs.SHA1Hashes...)
    allHashes = append(allHashes, iocs.SHA256Hashes...)
    if len(allHashes) > 0 {
        queries = append(queries, b.generateHashQuery(allHashes))
    }
    
    if len(iocs.Domains) > 0 {
        queries = append(queries, b.generateDomainQuery(iocs.Domains))
    }
    
    if len(iocs.IPv4Addresses) > 0 {
        queries = append(queries, b.generateIPQuery(iocs.IPv4Addresses))
    }
    
    return queries, nil
}
```

**CLI Formatting**:
```go
queries, err := backend.GenerateQuery(iocs)
output := strings.Join(queries, "\n\n")  // Double newline separator
```

## Implementation Patterns

### Package Structure
```
pkg/backends/r7/
├── r7.go           # Main backend implementation
└── r7_test.go      # Table-driven unit tests
```

### Type Definition
```go
type R7Backend struct {
    config *Config
}

type Config struct {
    // Future: Add configuration options if needed
}
```

### Query Generation Methods
Following S1QL pattern:
- `generateHashQuery(allHashes []string) string` - Combined hash query
- `generateDomainQuery(domains []string) string` - Domain ICONTAINS-ANY query
- `generateIPQuery(ips []string) string` - IP address IN query
- `formatStringList(items []string) string` - Helper for list formatting
- `escapeString(s string) string` - Escape special characters

## Testing Strategy

### Unit Tests (Table-Driven)
Pattern from S1QL backend:
```go
func TestR7Backend_GenerateQuery(t *testing.T) {
    backend := New(nil)
    tests := []struct {
        name     string
        iocs     *extraction.IOCSet
        want     []string
        wantErr  bool
    }{
        {
            name: "single MD5 hash",
            iocs: &extraction.IOCSet{
                MD5Hashes: []string{"5d41402abc4b2a76b9719d911017c592"},
            },
            want: []string{`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5d41402abc4b2a76b9719d911017c592'])`},
            wantErr: false,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := backend.GenerateQuery(tt.iocs)
            if (err != nil) != tt.wantErr {
                t.Errorf("GenerateQuery() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if len(got) != len(tt.want) {
                t.Errorf("GenerateQuery() returned %d queries, want %d", len(got), len(tt.want))
                return
            }
            for i := range got {
                if got[i] != tt.want[i] {
                    t.Errorf("GenerateQuery()[%d] = %v, want %v", i, got[i], tt.want[i])
                }
            }
        })
    }
}
```

### Integration Tests
- End-to-end CLI tests with `--backend r7`
- Sample IOC files with expected output
- Compare generated queries against golden files

### Test Coverage Requirements
- All IOC types (MD5, SHA1, SHA256, domains, IPs)
- Single and multiple IOCs per type
- Combined IOC types in one query
- Empty IOCSet error handling
- Special character escaping
- Separate query mode

## Error Handling
Match S1QL backend patterns:
- Empty IOCSet → descriptive error
- Invalid input → clear error message
- Large query warning (stderr) if IOC count > 1000

## String Escaping
Rapid7 AQL special characters to escape:
- Single quotes: `'` → `\'` (used for string delimiters)
- Backslashes: `\` → `\\`

Pattern:
```go
func (b *R7Backend) escapeString(s string) string {
    s = strings.ReplaceAll(s, `\`, `\\`)
    s = strings.ReplaceAll(s, `'`, `\'`)
    return s
}
```

## CLI Integration
No changes required to CLI structure:
- `--backend r7` automatically works via backend registry pattern
- Backend selection in main.go:
  ```go
  case "r7":
      backend = r7.New(nil)
  ```

## Documentation Requirements
- Update README.md with Rapid7 backend examples
- Add query syntax examples for each IOC type
- Document ICONTAINS-ANY vs IN operator differences
- Note InsightIDR version compatibility

## Future Extensions (Out of Scope)
- Custom field mappings via config file
- Time range filters (`where(...) AND timestamp > X`)
- Advanced operators (regex, CIDR matching)
- Query optimization (field order, index hints)
- Multi-platform support (InsightIDR vs InsightConnect)

## Open Questions
None - all clarifications received from user.
