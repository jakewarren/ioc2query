# r7-backend Specification

## Purpose
Backend implementation for generating Rapid7 InsightIDR Advanced Query Language (AQL) queries from extracted IOCs. Enables security analysts to transform indicators of compromise into Rapid7-specific detection queries.

## ADDED Requirements

### Requirement: System SHALL generate query for all hash types
The system MUST generate valid Rapid7 AQL queries to search for file hashes (MD5, SHA1, SHA256) using wildcard pattern matching.

**Acceptance Criteria:**
- Uses `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN [...])` syntax
- Combines all hash types (MD5, SHA1, SHA256) into single hash list
- Properly formats hash values in single-quoted array
- Handles single and multiple hashes
- Generates syntactically valid Rapid7 AQL

#### Scenario: Single hash query
**Given** IOCSet contains one SHA256 hash "5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4"  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element: `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4'])`

#### Scenario: Multiple hashes of different types
**Given** IOCSet contains MD5 "5d41402abc4b2a76b9719d911017c592", SHA1 "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d", and SHA256 "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element containing all three hashes in single list: `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5d41402abc4b2a76b9719d911017c592', 'aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d', 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855'])`

---

### Requirement: System SHALL generate query for domain names
The system MUST generate valid Rapid7 AQL queries to search for domain names in DNS and URL events using case-insensitive partial matching.

**Acceptance Criteria:**
- Uses `where("query","url" ICONTAINS-ANY [...])` syntax
- Properly formats domain values in single-quoted array
- Handles multiple domains
- Generates syntactically correct Rapid7 AQL
- Case-insensitive matching enabled via ICONTAINS-ANY

#### Scenario: Single domain query
**Given** IOCSet contains domain "icloud.com"  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element: `where("query","url" ICONTAINS-ANY ['icloud.com'])`

#### Scenario: Multiple domains query
**Given** IOCSet contains domains ["evil.com", "malware.net", "bad.org"]  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element: `where("query","url" ICONTAINS-ANY ['evil.com', 'malware.net', 'bad.org'])`

---

### Requirement: System SHALL generate query for IPv4 addresses
The system MUST generate valid Rapid7 AQL queries to search for IPv4 addresses in network events.

**Acceptance Criteria:**
- Uses `where("source_address","destination_address" IN [...])` syntax
- Searches both source and destination address fields
- Handles single and multiple IP addresses
- Properly formats IP addresses in single-quoted array
- Generates syntactically valid Rapid7 AQL

#### Scenario: Single IPv4 query
**Given** IOCSet contains IP address "1.1.1.1"  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element: `where("source_address","destination_address" IN ['1.1.1.1'])`

#### Scenario: Multiple IPv4 addresses
**Given** IOCSet contains IPs ["10.0.0.1", "192.168.1.100", "8.8.8.8"]  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element: `where("source_address","destination_address" IN ['10.0.0.1', '192.168.1.100', '8.8.8.8'])`

---

### Requirement: System SHALL generate queries grouped by IOC type
The system MUST generate Rapid7 AQL queries grouped by IOC type (hashes, domains, IPs), returning them as separate query strings.

**Acceptance Criteria:**
- GenerateQuery returns []string containing separate queries per IOC type
- Hashes (all types combined) generate one query element if present
- Domains generate one query element if present  
- IPs generate one query element if present
- Each query element is independently valid and executable
- Presentation layer (CLI/web) handles joining with appropriate formatting
- Each IOC type uses appropriate operator (IN for hashes/IPs, ICONTAINS-ANY for domains)

#### Scenario: Mixed IOC types
**Given** IOCSet contains hashes, domains, and IPv4 addresses  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with three elements: hash query, domain query, IP query

#### Scenario: All IOC types present
**Given** IOCSet contains MD5, SHA1, SHA256, domains, and IPs  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 3 elements:
- Element 0: `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['md5hash', 'sha1hash', 'sha256hash'])`
- Element 1: `where("query","url" ICONTAINS-ANY ['evil.com'])`
- Element 2: `where("source_address","destination_address" IN ['1.2.3.4'])`

#### Scenario: Single IOC type
**Given** IOCSet contains only hashes  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with 1 element containing the hash query

#### Scenario: CLI formatting
**Given** GenerateQuery returns []string with multiple elements  
**When** CLI outputs the queries  
**Then** the system joins elements with "\n\n" (double newline) for readability

---

### Requirement: System SHALL generate separate queries
The system MUST support generating individual queries for each IOC when requested.

**Acceptance Criteria:**
- GenerateQueries returns slice of individual query strings
- One query per IOC indicator
- Each query is independently valid
- Maintains consistent formatting across queries

#### Scenario: Separate queries for multiple IOCs
**Given** IOCSet contains 3 hashes and 2 domains  
**When** GenerateQueries is called for R7 backend  
**Then** the system returns 5 separate query strings, one for each IOC

#### Scenario: Separate queries maintain correct syntax
**Given** IOCSet contains one MD5 hash "abc123" and one domain "test.com"  
**When** GenerateQueries is called for R7 backend  
**Then** the system returns 2 queries: `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['abc123'])` and `where("query","url" ICONTAINS-ANY ['test.com'])`

---

### Requirement: System SHALL handle empty IOCSet
The system MUST handle empty IOCSets gracefully without generating invalid queries.

**Acceptance Criteria:**
- Returns error when IOCSet is empty or all lists are empty
- Error message clearly indicates no IOCs provided
- Does not generate empty or invalid query syntax
- Error is descriptive for debugging

#### Scenario: Empty IOCSet error
**Given** IOCSet has all empty lists (no hashes, domains, or IPs)  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns an error: "IOC set is empty"

#### Scenario: Nil IOCSet error
**Given** IOCSet is nil  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns an error indicating invalid input

---

### Requirement: System SHALL escape special characters
The system MUST properly escape special characters in IOC values to prevent query syntax errors.

**Acceptance Criteria:**
- Escapes single quotes in domain names and other string values
- Handles backslashes correctly
- Prevents query syntax breaking
- Maintains semantic meaning of IOCs after escaping

#### Scenario: Domain with single quote
**Given** IOCSet contains domain "test's-site.com"  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with properly escaped single quote: `where("query","url" ICONTAINS-ANY ['test\'s-site.com'])`

#### Scenario: Domain with backslash
**Given** IOCSet contains domain "test\\example.com"  
**When** GenerateQuery is called for R7 backend  
**Then** the system returns []string with properly escaped backslash: `where("query","url" ICONTAINS-ANY ['test\\\\example.com'])`

---

### Requirement: System SHALL warn on large queries
The system SHALL warn users when generating queries with very large numbers of IOCs.

**Acceptance Criteria:**
- Outputs warning to stderr when IOC count exceeds 1000
- Warning message indicates query size
- Does not prevent query generation
- Warning is visible but does not interfere with query output

#### Scenario: Large IOC set warning
**Given** IOCSet contains 1500 total IOCs across all types  
**When** GenerateQuery is called for R7 backend  
**Then** the system writes warning to stderr: "Warning: Query contains 1500 IOCs, which may be very large"  
**And** the system still generates and returns the complete query

---

### Requirement: System SHALL implement backend interface
The R7 backend MUST implement the Backend interface defined in pkg/backends.

**Acceptance Criteria:**
- Implements Name() method returning "r7"
- Implements GenerateQuery(iocs) method
- Implements GenerateQueries(iocs) method
- Returns appropriate errors for invalid inputs
- Is compatible with existing CLI backend selection

#### Scenario: Backend name identification
**Given** R7 backend instance is created  
**When** Name() method is called  
**Then** the system returns "r7" as identifier

#### Scenario: Backend interface compatibility
**Given** R7 backend instance is created  
**When** assigned to backends.Backend interface variable  
**Then** the system compiles without error and backend methods are callable

---

### Requirement: Format queries consistently
The system SHALL format generated queries with consistent spacing and quotation for readability.

**Acceptance Criteria:**
- Uses single quotes for string values in arrays
- Consistent spacing around operators and commas
- Field names in double quotes
- Maintains Rapid7 AQL conventions

#### Scenario: Consistent array formatting
**Given** IOCSet contains multiple hashes  
**When** GenerateQuery is called for R7 backend  
**Then** the system formats array as `['hash1', 'hash2', 'hash3']` with single quotes and comma-space separation

#### Scenario: Consistent field name formatting
**Given** any IOCSet with data  
**When** GenerateQuery is called for R7 backend  
**Then** field names are wrapped in double quotes: `where("field1","field2" OPERATOR [...])`
