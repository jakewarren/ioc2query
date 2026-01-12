# s1ql-backend Specification

## Purpose
TBD - created by archiving change add-s1-query-support. Update Purpose after archive.
## Requirements
### Requirement: System SHALL generate query for MD5 hashes
The system MUST generate valid S1QLv2 queries to search for MD5 file hashes.

**Acceptance Criteria:**
- Uses `src.process.image.md5 in (...) || tgt.file.md5 in (...)` syntax
- Properly quotes and formats hash values in query
- Handles single and multiple MD5 hashes
- Generates syntactically valid S1QL

#### Scenario: Single MD5 hash query
**Given** IOCSet contains one MD5 hash "5d41402abc4b2a76b9719d911017c592"  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query: `src.process.image.md5 in ("5d41402abc4b2a76b9719d911017c592") || tgt.file.md5 in ("5d41402abc4b2a76b9719d911017c592")`

#### Scenario: Multiple MD5 hashes query
**Given** IOCSet contains MD5 hashes ["hash1", "hash2", "hash3"]  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query with `src.process.image.md5 in ("hash1", "hash2", "hash3") || tgt.file.md5 in ("hash1", "hash2", "hash3")`

---

### Requirement: System SHALL generate query for SHA1 hashes
The system MUST generate valid S1QLv2 queries to search for SHA1 file hashes.

**Acceptance Criteria:**
- Uses `src.process.image.sha1 in (...) || tgt.file.sha1 in (...)` syntax
- Properly formats SHA1 hash values
- Handles multiple SHA1 hashes in single query
- Maintains S1QL syntax correctness

#### Scenario: Single SHA1 hash query
**Given** IOCSet contains one SHA1 hash "2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query: `src.process.image.sha1 in ("2aae6c35c94fcfb415dbe95f408b9ce91ee846ed") || tgt.file.sha1 in ("2aae6c35c94fcfb415dbe95f408b9ce91ee846ed")`

---

### Requirement: System SHALL generate query for SHA256 hashes
The system MUST generate valid S1QLv2 queries to search for SHA256 file hashes.

**Acceptance Criteria:**
- Uses `src.process.image.sha256 in (...) || tgt.file.sha256 in (...)` syntax
- Properly formats SHA256 hash values (64 hex chars)
- Supports multiple SHA256 hashes
- Generates valid S1QL syntax

#### Scenario: Single SHA256 hash query
**Given** IOCSet contains one SHA256 hash "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query: `src.process.image.sha256 in ("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855") || tgt.file.sha256 in ("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")`

---

### Requirement: System SHALL generate query for domain names
The system MUST generate valid S1QLv2 queries to search for domain names in network events.

**Acceptance Criteria:**
- Uses `event.dns.request in (...)` syntax
- Properly escapes special characters in domain names
- Handles multiple domains with in operator
- Generates syntactically correct S1QL

#### Scenario: Single domain query
**Given** IOCSet contains domain "malicious.example.com"  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query: `event.dns.request in ("malicious.example.com")`

#### Scenario: Multiple domains query
**Given** IOCSet contains domains ["evil.com", "bad.org", "malware.net"]  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query: `event.dns.request in ("evil.com", "bad.org", "malware.net")`

---

### Requirement: System SHALL generate query for IPv4 addresses
The system MUST generate valid S1QLv2 queries to search for IPv4 addresses in network events.

**Acceptance Criteria:**
- Searches both src.ip.address and dst.ip.address fields
- Uses `(src.ip.address in (...) || dst.ip.address in (...))` syntax
- Handles single and multiple IP addresses
- Properly formats IP addresses in query

#### Scenario: Single IPv4 query
**Given** IOCSet contains IP address "192.168.1.100"  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query: `(src.ip.address = "192.168.1.100" || dst.ip.address = "192.168.1.100")`

#### Scenario: Multiple IPv4 query
**Given** IOCSet contains IPs ["10.0.0.1", "10.0.0.2", "10.0.0.3"]  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns query with `src.ip.address in (...) || dst.ip.address in (...)` for all IPs

---

### Requirement: System SHALL generate combined query
The system MUST generate a single combined S1QL query when IOCSet contains multiple IOC types.

**Acceptance Criteria:**
- Combines all IOC types with OR logic
- Groups file hash queries together
- Groups network queries together
- Maintains proper parentheses and precedence
- Query is syntactically valid S1QL

#### Scenario: Mixed IOC types
**Given** IOCSet contains MD5 hashes, domains, and IPv4 addresses  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns a single query combining all types with OR operators

#### Scenario: All IOC types present
**Given** IOCSet contains MD5, SHA1, SHA256, domains, and IPs  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns properly structured query: `(file conditions) || (network conditions)`

---

### Requirement: System SHALL generate separate queries
The system MUST support generating individual queries for each IOC when requested.

**Acceptance Criteria:**
- GenerateQueries returns slice of individual query strings
- One query per IOC indicator
- Each query is independently valid
- Maintains consistent formatting across queries

#### Scenario: Separate queries for multiple IOCs
**Given** IOCSet contains 3 MD5 hashes and 2 domains  
**When** GenerateQueries is called for S1QL backend  
**Then** the system returns 5 separate query strings, one for each IOC

---

### Requirement: System SHALL handle empty IOCSet
The system MUST handle empty IOCSets gracefully without generating invalid queries.

**Acceptance Criteria:**
- Returns error when IOCSet is empty or all lists are empty
- Error message clearly indicates no IOCs provided
- Does not generate empty or invalid query syntax

#### Scenario: Empty IOCSet
**Given** IOCSet has all empty lists (no hashes, domains, or IPs)  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns an error indicating no IOCs to process

---

### Requirement: Format queries for readability
The system SHALL format generated queries with proper indentation and line breaks for readability.

**Acceptance Criteria:**
- Multi-line queries use consistent indentation
- OR operators appear on appropriate lines
- Hash lists are formatted cleanly
- Queries remain valid when formatted

#### Scenario: Format complex query
**Given** IOCSet contains multiple IOC types  
**When** GenerateQuery is called for S1QL backend  
**Then** the system returns a formatted multi-line query that is human-readable

---

### Requirement: System SHALL escape special characters
The system MUST properly escape special characters in IOC values to prevent query injection or syntax errors.

**Acceptance Criteria:**
- Escapes quotes in domain names
- Handles backslashes and other special SQL characters
- Prevents query syntax breaking
- Maintains semantic meaning of IOCs

#### Scenario: Domain with special characters
**Given** IOCSet contains domain "test-site.example.com"  
**When** GenerateQuery is called for S1QL backend  
**Then** the system properly escapes or handles the hyphen without breaking query syntax

---

### Requirement: System SHALL implement backend interface
The S1QL backend MUST implement the Backend interface defined in pkg/backends.

**Acceptance Criteria:**
- Implements Name() method returning "s1"
- Implements GenerateQuery(iocs) method
- Implements GenerateQueries(iocs) method
- Returns appropriate errors for invalid inputs

#### Scenario: Backend name identification
**Given** S1QL backend instance is created  
**When** Name() method is called  
**Then** the system returns "s1" as identifier

---

### Requirement: Query size warning
The system SHALL warn when generated queries become very large (>1000 IOCs).

**Acceptance Criteria:**
- Logs warning when IOC count exceeds threshold
- Warning does not prevent query generation
- Suggests using separate queries for performance
- Threshold is configurable

#### Scenario: Large IOC set
**Given** IOCSet contains 1500 MD5 hashes  
**When** GenerateQuery is called for S1QL backend  
**Then** the system logs a warning about query size but still generates the query

