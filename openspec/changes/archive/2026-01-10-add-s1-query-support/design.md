# Design: Add S1 Query Support

## Architecture Overview

```
┌─────────────────┐
│   CLI Layer     │  (cobra/flags, input/output handling)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  IOC Extractor  │  (uses go-ioc library, deduplication)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Backend Router  │  (selects query generator)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  S1QL Backend   │  (generates S1QLv2 queries)
└─────────────────┘
```

## Module Design

### 1. IOC Extraction Module (`pkg/extraction`)

**Purpose:** Extract and normalize IOCs from raw text input.

**Key Components:**
```go
// IOCSet holds deduplicated IOCs by type
type IOCSet struct {
    MD5Hashes    []string
    SHA1Hashes   []string
    SHA256Hashes []string
    Domains      []string
    IPv4Addresses []string
}

// Extractor wraps go-ioc functionality
type Extractor struct {
    parser *ioc.IOCParser
}

// Extract returns deduplicated IOCs from input text
func (e *Extractor) Extract(input string) (*IOCSet, error)
```

**Implementation Details:**
- Use `github.com/vertoforce/go-ioc` for extraction
- Normalize all hashes to lowercase
- Deduplicate using Go maps (order preservation via slice)
- Validate IOC format (regex validation for domains/IPs)
- Handle defanged IOCs (e.g., "example[.]com" → "example.com")

**Error Handling:**
- Return error if input is empty
- Log warnings for malformed IOCs (don't fail extraction)
- Return empty IOCSet if no valid IOCs found

### 2. Backend Interface (`pkg/backends`)

**Purpose:** Define common interface for all query backends.

**Interface Design:**
```go
// Backend generates vendor-specific queries from IOCs
type Backend interface {
    // Name returns the backend identifier (e.g., "s1", "rapid7")
    Name() string
    
    // GenerateQuery creates a single combined query
    GenerateQuery(iocs *extraction.IOCSet) (string, error)
    
    // GenerateQueries creates individual queries per IOC
    GenerateQueries(iocs *extraction.IOCSet) ([]string, error)
}
```

**Design Rationale:**
- Interface enables easy addition of future backends (Rapid7, etc.)
- Two generation modes support different analyst workflows
- Simple string output for CLI integration

### 3. S1QL Backend (`pkg/backends/s1ql`)

**Purpose:** Generate SentinelOne S1QLv2 queries.

**Key Components:**
```go
type S1QLBackend struct {
    config *Config
}

type Config struct {
    CombineWithOR bool // Combine IOCs in single query vs separate
    IncludeComments bool // Add comments explaining query parts
}

func New(config *Config) *S1QLBackend
func (b *S1QLBackend) Name() string
func (b *S1QLBackend) GenerateQuery(iocs *extraction.IOCSet) (string, error)
func (b *S1QLBackend) GenerateQueries(iocs *extraction.IOCSet) ([]string, error)
```

**S1QL Query Patterns:**

1. **File Hash Queries:**
```sql
-- MD5
EventType = "File" AND Md5 IN ("hash1", "hash2", ...)

-- SHA1
EventType = "File" AND Sha1 IN ("hash1", "hash2", ...)

-- SHA256
EventType = "File" AND Sha256 IN ("hash1", "hash2", ...)

-- Combined hashes
(EventType = "File" AND Md5 IN (...)) OR 
(EventType = "File" AND Sha1 IN (...)) OR
(EventType = "File" AND Sha256 IN (...))
```

2. **Domain Queries:**
```sql
-- Single domain
EventType = "Network" AND DnsRequest CONTAINS "malicious.com"

-- Multiple domains
EventType = "Network" AND (
    DnsRequest CONTAINS "domain1.com" OR
    DnsRequest CONTAINS "domain2.com" OR
    ...
)
```

3. **IPv4 Queries:**
```sql
-- Single IP
EventType = "Network" AND (
    SrcIp = "192.168.1.1" OR
    DstIp = "192.168.1.1"
)

-- Multiple IPs
EventType = "Network" AND (
    SrcIp IN ("ip1", "ip2", ...) OR
    DstIp IN ("ip1", "ip2", ...)
)
```

4. **Combined Query (All IOC Types):**
```sql
(EventType = "File" AND Md5 IN (...)) OR
(EventType = "File" AND Sha1 IN (...)) OR
(EventType = "File" AND Sha256 IN (...)) OR
(EventType = "Network" AND (DnsRequest CONTAINS "..." OR ...)) OR
(EventType = "Network" AND (SrcIp IN (...) OR DstIp IN (...)))
```

**Implementation Details:**
- Escape special characters in domain names
- Quote all string literals properly
- Format multi-line queries for readability
- Limit query size (warn if >1000 IOCs)

**Error Handling:**
- Return error if IOCSet is empty
- Warn if query becomes very large
- Validate S1QL syntax basics

### 4. CLI Interface (`cmd/ioc2query`)

**Purpose:** Provide command-line interface for the tool.

**Command Structure:**
```bash
# Basic usage (stdin)
cat iocs.txt | ioc2query --backend s1

# File input
ioc2query --backend s1 --input iocs.txt

# Output to file
ioc2query --backend s1 --input iocs.txt --output query.txt

# Separate queries per IOC
ioc2query --backend s1 --input iocs.txt --separate

# Verbose output
ioc2query --backend s1 --input iocs.txt --verbose
```

**Flags:**
- `--backend, -b` (required): Backend to use (s1, rapid7)
- `--input, -i` (optional): Input file (default: stdin)
- `--output, -o` (optional): Output file (default: stdout)
- `--separate, -s` (optional): Generate separate queries per IOC
- `--verbose, -v` (optional): Enable verbose logging
- `--help, -h`: Show usage information

**Exit Codes:**
- `0`: Success
- `1`: Input/output error
- `2`: Invalid arguments
- `3`: Extraction error
- `4`: Query generation error

**Output Format:**
- Default: Clean query output (stdout)
- Verbose: Logs extraction stats to stderr
- Errors: Clear error messages to stderr

## Data Flow

1. **Input:** Raw text from stdin or file
2. **Extraction:** Parse text → IOCSet (deduplicated)
3. **Backend Selection:** Choose query generator
4. **Query Generation:** IOCSet → S1QL queries
5. **Output:** Formatted queries to stdout or file

## Testing Strategy

### Unit Tests
- **Extraction:** Test each IOC type extraction, deduplication, edge cases
- **S1QL Backend:** Test query generation for each IOC type combination
- **CLI:** Test flag parsing, error handling, input/output

### Integration Tests
- End-to-end: Input file → CLI → Query output validation
- Sample inputs with known expected outputs
- Cross-platform testing (Linux, macOS, Windows)

### Test Coverage Goals
- Target: >85% code coverage
- Critical paths: 100% coverage (extraction, query generation)

## Performance Considerations
- IOC extraction: O(n) where n = input size
- Deduplication: O(m) where m = number of IOCs
- Query generation: O(m) - linear with IOC count
- Target: <1 second for inputs with <1000 IOCs

## Security Considerations
- No external network calls (offline tool)
- Input sanitization to prevent injection (query escaping)
- No credential handling (query generation only)
- Safe handling of potentially malicious indicator strings

## Extensibility
- Backend interface allows easy addition of new query languages
- Config struct enables per-backend customization
- Modular design supports future enhancements (query optimization, etc.)

## Dependencies
- **Runtime:** `github.com/vertoforce/go-ioc` for IOC extraction
- **CLI Framework:** cobra or standard library (TBD in implementation)
- **Build:** Go 1.21+ for language features
- **CI/CD:** GitHub Actions for automated testing

## Open Design Questions
1. Should we implement query result pagination for very large IOC sets?
2. Do we need query validation against actual S1QL parser?
3. Should we support query templates for custom S1QL patterns?

## Future Enhancements (Out of Scope)
- Query execution via SentinelOne API
- Query result parsing and formatting
- Interactive query builder mode
- Support for S1QL advanced features (aggregations, joins)
