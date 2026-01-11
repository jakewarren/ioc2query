# ioc-extraction Specification

## Purpose
TBD - created by archiving change add-s1-query-support. Update Purpose after archive.
## Requirements
### Requirement: System SHALL extract MD5 hashes
The system MUST extract all valid MD5 hash values (32 hexadecimal characters) from input text.

**Acceptance Criteria:**
- Recognizes both uppercase and lowercase MD5 hashes
- Normalizes all extracted hashes to lowercase
- Handles hashes embedded in surrounding text
- Ignores invalid MD5 patterns (incorrect length, non-hex characters)

#### Scenario: Extract MD5 from threat intelligence report
**Given** input text contains "File hash: 5D41402ABC4B2A76B9719D911017C592 was detected"  
**When** the extraction function processes the input  
**Then** the system returns ["5d41402abc4b2a76b9719d911017c592"] in the MD5 hashes list

#### Scenario: Handle multiple MD5 hashes
**Given** input text contains "Hashes: 5d41402abc4b2a76b9719d911017c592 and 098F6BCD4621D373CADE4E832627B4F6"  
**When** the extraction function processes the input  
**Then** the system returns both hashes normalized to lowercase in the MD5 hashes list

---

### Requirement: System SHALL extract SHA1 hashes
The system MUST extract all valid SHA1 hash values (40 hexadecimal characters) from input text.

**Acceptance Criteria:**
- Recognizes both uppercase and lowercase SHA1 hashes
- Normalizes all extracted hashes to lowercase
- Distinguishes SHA1 from MD5 and SHA256 based on length
- Handles hashes in various text contexts

#### Scenario: Extract SHA1 from malware analysis
**Given** input text contains "SHA1: 2AAE6C35C94FCFB415DBE95F408B9CE91EE846ED detected in memory"  
**When** the extraction function processes the input  
**Then** the system returns ["2aae6c35c94fcfb415dbe95f408b9ce91ee846ed"] in the SHA1 hashes list

---

### Requirement: System SHALL extract SHA256 hashes
The system MUST extract all valid SHA256 hash values (64 hexadecimal characters) from input text.

**Acceptance Criteria:**
- Recognizes both uppercase and lowercase SHA256 hashes
- Normalizes all extracted hashes to lowercase
- Correctly identifies 64-character hex strings as SHA256
- Handles multi-line and formatted hash representations

#### Scenario: Extract SHA256 from incident report
**Given** input text contains "File SHA256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"  
**When** the extraction function processes the input  
**Then** the system returns the hash in the SHA256 hashes list

---

### Requirement: System SHALL extract domain names
The system MUST extract all valid domain names from input text, including both fanged and defanged formats.

**Acceptance Criteria:**
- Extracts standard domain formats (example.com, subdomain.example.com)
- Handles defanged domains (example[.]com, example[dot]com)
- Normalizes defanged domains to standard format
- Validates TLD format (at least 2 characters)
- Excludes invalid domain patterns

#### Scenario: Extract fanged domain
**Given** input text contains "C2 server at malicious.example.com was contacted"  
**When** the extraction function processes the input  
**Then** the system returns ["malicious.example.com"] in the domains list

#### Scenario: Extract defanged domain
**Given** input text contains "Domain: evil[.]example[.]com was observed"  
**When** the extraction function processes the input  
**Then** the system returns ["evil.example.com"] (fanged) in the domains list

#### Scenario: Extract subdomain
**Given** input text contains "api.malware-c2.example.org responded with payload"  
**When** the extraction function processes the input  
**Then** the system returns ["api.malware-c2.example.org"] in the domains list

---

### Requirement: System SHALL extract IPv4 addresses
The system MUST extract all valid IPv4 addresses from input text, including both fanged and defanged formats.

**Acceptance Criteria:**
- Extracts standard IPv4 format (dotted decimal: 192.168.1.1)
- Handles defanged IPs (192[.]168[.]1[.]1)
- Normalizes defanged IPs to standard format
- Validates IP octet ranges (0-255)
- Excludes invalid IP patterns (e.g., 999.999.999.999)

#### Scenario: Extract standard IPv4
**Given** input text contains "Connection to 192.168.1.100 detected"  
**When** the extraction function processes the input  
**Then** the system returns ["192.168.1.100"] in the IPv4 addresses list

#### Scenario: Extract defanged IPv4
**Given** input text contains "Attacker IP: 10[.]0[.]0[.]5 observed"  
**When** the extraction function processes the input  
**Then** the system returns ["10.0.0.5"] (fanged) in the IPv4 addresses list

#### Scenario: Extract multiple IPs from log
**Given** input text contains "Traffic from 192.168.1.50 to 8.8.8.8 blocked"  
**When** the extraction function processes the input  
**Then** the system returns ["192.168.1.50", "8.8.8.8"] in the IPv4 addresses list

---

### Requirement: System SHALL deduplicate extracted IOCs
The system MUST deduplicate all extracted IOCs within each type category.

**Acceptance Criteria:**
- Removes exact duplicate hashes (case-insensitive comparison)
- Removes duplicate domains (case-insensitive comparison)
- Removes duplicate IP addresses
- Preserves order of first occurrence
- Applies deduplication per IOC type independently

#### Scenario: Deduplicate repeated MD5 hashes
**Given** input contains "Hash 5d41402abc4b2a76b9719d911017c592 and again 5D41402ABC4B2A76B9719D911017C592"  
**When** the extraction function processes the input  
**Then** the system returns only one instance of the hash in the MD5 list

#### Scenario: Deduplicate domains with different casing
**Given** input contains "malicious.com and MALICIOUS.COM and Malicious.Com"  
**When** the extraction function processes the input  
**Then** the system returns ["malicious.com"] (single deduplicated entry)

---

### Requirement: System SHALL handle empty or invalid input
The system MUST gracefully handle empty or invalid input without crashing.

**Acceptance Criteria:**
- Returns empty IOCSet for empty string input
- Returns empty IOCSet when no valid IOCs found
- Logs warnings for malformed but non-fatal patterns
- Does not crash on unexpected input formats

#### Scenario: Process empty input
**Given** input text is an empty string ""  
**When** the extraction function processes the input  
**Then** the system returns an empty IOCSet with no errors

#### Scenario: Process input with no IOCs
**Given** input text contains "This is a normal sentence with no indicators"  
**When** the extraction function processes the input  
**Then** the system returns an empty IOCSet with all lists empty

---

### Requirement: System SHALL use go-ioc library
The system MUST use the `github.com/vertoforce/go-ioc` library for IOC pattern recognition and extraction.

**Acceptance Criteria:**
- Integrates go-ioc parser for regex-based extraction
- Leverages built-in defanging/fanging capabilities
- Wraps library in internal interface for testability
- Handles library errors appropriately

#### Scenario: Initialize go-ioc parser
**Given** the extraction module is initialized  
**When** a new Extractor instance is created  
**Then** the system successfully initializes the go-ioc parser without errors

#### Scenario: Extract using go-ioc
**Given** input contains mixed IOC types  
**When** the go-ioc parser processes the input  
**Then** the system correctly extracts all supported IOC types using the library's regex patterns

