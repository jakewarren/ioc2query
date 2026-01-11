# Specification: CLI Interface

## ADDED Requirements

### Requirement: System SHALL accept input from stdin
The system MUST accept raw text input from standard input (stdin) by default.

**Acceptance Criteria:**
- Reads from stdin when no input file specified
- Processes entire stdin stream until EOF
- Handles large input streams efficiently
- Works with piped input from other commands

#### Scenario: Pipe input to CLI
**Given** a file "iocs.txt" contains IOC data  
**When** user runs `cat iocs.txt | ioc2query --backend s1`  
**Then** the system reads the input from stdin and generates queries

#### Scenario: Interactive stdin input
**Given** user runs `ioc2query --backend s1` without piping  
**When** user types IOC data and presses Ctrl+D (EOF)  
**Then** the system processes the typed input and generates queries

---

### Requirement: System SHALL accept input from file
The system MUST accept input from a file specified via command-line flag.

**Acceptance Criteria:**
- Provides `--input` or `-i` flag for file path
- Reads entire file contents
- Returns error if file doesn't exist or can't be read
- Supports relative and absolute file paths

#### Scenario: Read from input file
**Given** a file "indicators.txt" exists with IOC data  
**When** user runs `ioc2query --backend s1 --input indicators.txt`  
**Then** the system reads the file contents and generates queries

#### Scenario: Handle missing input file
**Given** file "missing.txt" does not exist  
**When** user runs `ioc2query --backend s1 --input missing.txt`  
**Then** the system returns error "file not found" with exit code 1

---

### Requirement: System SHALL output to stdout
The system MUST output generated queries to standard output (stdout) by default.

**Acceptance Criteria:**
- Writes query output to stdout when no output file specified
- Only outputs queries, not status messages (those go to stderr)
- Supports redirection to files via shell operators
- Clean output suitable for piping to other commands

#### Scenario: Output to stdout
**Given** IOC input is provided  
**When** user runs `ioc2query --backend s1` without output flag  
**Then** the system writes generated query to stdout

#### Scenario: Redirect output to file
**Given** IOC input is provided  
**When** user runs `ioc2query --backend s1 > output.txt`  
**Then** the system writes query to output.txt via stdout redirection

---

### Requirement: System SHALL output to file
The system MUST support writing output to a file specified via command-line flag.

**Acceptance Criteria:**
- Provides `--output` or `-o` flag for file path
- Creates output file if it doesn't exist
- Overwrites existing output file (with warning)
- Returns error if file can't be written

#### Scenario: Write to output file
**Given** IOC input is provided  
**When** user runs `ioc2query --backend s1 --output query.txt`  
**Then** the system creates query.txt with the generated query

#### Scenario: Handle write permission error
**Given** output path is not writable  
**When** user runs `ioc2query --backend s1 --output /root/query.txt`  
**Then** the system returns error "permission denied" with exit code 1

---

### Requirement: System SHALL require backend selection
The system MUST require backend selection via `--backend` or `-b` flag, or use the `--s1` shortcut.

**Acceptance Criteria:**
- Flag accepts backend identifier "s1"
- Provides `--s1` flag as shortcut for `--backend s1`
- Returns error if flag is missing
- Returns error if backend identifier is invalid
- Case-insensitive backend matching

#### Scenario: Specify S1 backend
**Given** user provides IOC input  
**When** user runs `ioc2query --backend s1`  
**Then** the system uses the S1QL backend for query generation

#### Scenario: Use S1 shortcut
**Given** user provides IOC input  
**When** user runs `ioc2query --s1`  
**Then** the system uses the S1QL backend for query generation

#### Scenario: Missing backend flag
**Given** user provides IOC input  
**When** user runs `ioc2query` without --backend or --s1 flag  
**Then** the system returns error "backend required" with exit code 2

#### Scenario: Invalid backend
**Given** user provides IOC input  
**When** user runs `ioc2query --backend invalid`  
**Then** the system returns error "unknown backend: invalid" with exit code 2

---

### Requirement: System SHALL support separate queries mode
The system MUST support generating separate queries per IOC via `--separate` or `-s` flag.

**Acceptance Criteria:**
- Flag triggers GenerateQueries instead of GenerateQuery
- Each IOC gets its own query in output
- Queries are separated by newlines or delimiters
- Mode works with both stdin and file input

#### Scenario: Generate separate queries
**Given** IOCSet contains 3 MD5 hashes  
**When** user runs `ioc2query --backend s1 --separate`  
**Then** the system outputs 3 separate query strings, one per hash

---

### Requirement: System SHALL support verbose logging
The system MUST support verbose output via `--verbose` or `-v` flag.

**Acceptance Criteria:**
- Enables detailed logging to stderr
- Shows extraction statistics (IOC counts by type)
- Shows processing steps and timing
- Does not interfere with query output to stdout

#### Scenario: Verbose extraction stats
**Given** input contains 5 MD5 hashes and 3 domains  
**When** user runs `ioc2query --backend s1 --verbose`  
**Then** stderr shows "Extracted: 5 MD5, 3 domains" before outputting query

---

### Requirement: System SHALL provide help information
The system MUST provide help and usage information via `--help` or `-h` flag.

**Acceptance Criteria:**
- Shows tool description and purpose
- Lists all available flags with descriptions
- Provides usage examples
- Shows supported backends and IOC types
- Displays version information

#### Scenario: Display help
**Given** user is unsure how to use the tool  
**When** user runs `ioc2query --help`  
**Then** the system displays comprehensive usage information and exits with code 0

---

### Requirement: System SHALL use appropriate exit codes
The system MUST use appropriate exit codes to indicate success or failure.

**Acceptance Criteria:**
- Exit code 0 for successful execution
- Exit code 1 for input/output errors
- Exit code 2 for invalid arguments
- Exit code 3 for extraction errors
- Exit code 4 for query generation errors

#### Scenario: Success exit code
**Given** valid IOC input and correct flags  
**When** command executes successfully  
**Then** the system exits with code 0

#### Scenario: Invalid argument exit code
**Given** user provides invalid flag  
**When** command parses arguments  
**Then** the system exits with code 2 and error message

---

### Requirement: System SHALL write errors to stderr
The system MUST write error messages to standard error (stderr), not stdout.

**Acceptance Criteria:**
- All errors go to stderr
- Stderr messages don't pollute stdout query output
- Error messages are clear and actionable
- Errors include context (file path, line number when relevant)

#### Scenario: File not found error
**Given** user specifies non-existent input file  
**When** command attempts to read file  
**Then** stderr contains "Error: file not found: <path>" and stdout is empty

---

### Requirement: System SHALL run cross-platform
The system MUST run on Linux, macOS, and Windows without platform-specific modifications.

**Acceptance Criteria:**
- Compiles for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64
- File path handling works on all platforms
- stdin/stdout/stderr work consistently
- No platform-specific code in critical paths

#### Scenario: Compile for multiple platforms
**Given** Go cross-compilation is configured  
**When** build scripts run  
**Then** binaries are created for Linux, macOS (Intel/ARM), and Windows

---

### Requirement: System SHALL handle empty input
The system MUST handle empty input gracefully without crashing.

**Acceptance Criteria:**
- Detects empty stdin or empty file
- Returns informative error message
- Exits with appropriate exit code
- Does not generate empty queries

#### Scenario: Empty stdin
**Given** stdin is empty (immediate EOF)  
**When** command reads input  
**Then** stderr shows "Error: no input provided" and exits with code 3

---

### Requirement: Performance for large inputs
The system SHALL process inputs with up to 1000 IOCs in under 1 second.

**Acceptance Criteria:**
- Handles large text inputs efficiently
- Memory usage remains reasonable (<100MB for 1000 IOCs)
- Streaming processing where possible
- No significant performance degradation

#### Scenario: Process 1000 IOCs
**Given** input file contains 1000 mixed IOCs  
**When** user runs `ioc2query --backend s1 --input large.txt`  
**Then** the system generates query in less than 1 second
