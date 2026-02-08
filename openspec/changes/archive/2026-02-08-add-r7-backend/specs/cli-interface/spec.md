# cli-interface Specification

## Purpose
Command-line interface for ioc2query with support for multiple backend query languages.

## MODIFIED Requirements

### Requirement: System SHALL require backend selection
The system MUST require backend selection via `--backend` or `-b` flag, or use backend shortcuts.

**Acceptance Criteria:**
- Flag accepts backend identifiers "s1" or "r7"
- Provides `--s1` flag as shortcut for `--backend s1`
- Provides `--r7` flag as shortcut for `--backend r7`
- Returns error if flag is missing
- Returns error if backend identifier is invalid
- Case-insensitive backend matching

#### Scenario: Specify S1 backend
**Given** user provides IOC input  
**When** user runs `ioc2query --backend s1`  
**Then** the system uses the S1QL backend for query generation

#### Scenario: Specify R7 backend
**Given** user provides IOC input  
**When** user runs `ioc2query --backend r7`  
**Then** the system uses the Rapid7 backend for query generation

#### Scenario: Use S1 shortcut
**Given** user provides IOC input  
**When** user runs `ioc2query --s1`  
**Then** the system uses the S1QL backend for query generation

#### Scenario: Use R7 shortcut
**Given** user provides IOC input  
**When** user runs `ioc2query --r7`  
**Then** the system uses the Rapid7 backend for query generation

#### Scenario: Missing backend flag
**Given** user provides IOC input  
**When** user runs `ioc2query` without --backend, --s1, or --r7 flag  
**Then** the system returns error "backend required" with exit code 2

#### Scenario: Invalid backend
**Given** user provides IOC input  
**When** user runs `ioc2query --backend invalid`  
**Then** the system returns error "unknown backend: invalid" with exit code 2

#### Scenario: Case-insensitive backend matching
**Given** user provides IOC input  
**When** user runs `ioc2query --backend R7` or `ioc2query --backend S1`  
**Then** the system accepts the backend identifier and uses the appropriate backend
