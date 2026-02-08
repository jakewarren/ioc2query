# Proposal: Add Rapid7 Backend Support

## Overview
Add Rapid7 InsightIDR backend to enable query generation for Rapid7 SIEM platform. This extends `ioc2query` to support a second vendor query language alongside the existing SentinelOne S1QL backend.

## Problem Statement
Security analysts using Rapid7 InsightIDR currently cannot use `ioc2query` to transform IOCs into vendor-specific detection queries. The project's architecture was designed to support multiple backends, but only SentinelOne S1QL is currently implemented.

## Proposed Solution
Implement a Rapid7 backend (`r7`) that:
- Implements the `backends.Backend` interface
- Generates Rapid7 InsightIDR AQL (Advanced Query Language) queries
- Supports all existing IOC types: MD5, SHA1, SHA256 hashes, domains, and IPv4 addresses
- Uses Rapid7-specific query syntax with `where()` clauses and appropriate operators

## Scope

### In Scope
- New `pkg/backends/r7/` package implementing Rapid7 query generation
- Backend identifier: `r7` (used with `--backend r7` flag)
- Query generation for:
  - **Hashes**: All types (MD5/SHA1/SHA256) using wildcard pattern `*.hashes.*`
  - **Domains**: Using `ICONTAINS-ANY` operator on `query,url` fields
  - **IPv4 addresses**: Using `IN` operator on `source_address,destination_address` fields
- Combined query mode (single query with all IOCs)
- Separate query mode (individual query per IOC)
- Unit tests following TDD approach
- Integration tests with realistic IOC examples
- CLI integration via existing `--backend` flag

### Out of Scope
- IPv6 addresses (not in v1 scope)
- URL extraction as separate IOC type (domains handled via ICONTAINS-ANY on url field)
- Custom field mappings or configuration
- Query optimization or performance tuning (initial implementation prioritizes correctness)
- Rapid7 API integration (query generation only)

## Affected Capabilities
- **cli-interface** (MODIFIED): Add "r7" as valid backend option
- **r7-backend** (ADDED): New capability for Rapid7 query generation

## Dependencies
- Existing `pkg/backends.Backend` interface (no changes required)
- Existing `pkg/extraction.IOCSet` structure (no changes required)
- CLI flag parsing (minimal change to add "r7" option)

## Risks & Considerations
- **Query syntax accuracy**: Examples provided must match production Rapid7 InsightIDR syntax
  - *Mitigation*: Test with real Rapid7 platform if possible, validate query format
- **Field name stability**: Rapid7 may change field names across platform versions
  - *Mitigation*: Document which InsightIDR version syntax is based on
- **Operator behavior**: ICONTAINS-ANY behavior with domains needs verification
  - *Mitigation*: Add comprehensive test cases covering edge cases

## Success Criteria
- `ioc2query --backend r7` generates valid Rapid7 InsightIDR queries
- All IOC types (hashes, domains, IPs) produce correct query syntax
- Unit test coverage matches S1QL backend quality (table-driven tests)
- Integration tests validate end-to-end functionality
- golangci-lint passes with no new warnings
- Documentation includes Rapid7 query examples

## Timeline Estimate
- Implementation: 2-3 hours (following S1QL backend patterns)
- Testing: 1-2 hours (unit + integration tests)
- Documentation: 30 minutes
- **Total**: ~4-6 hours

## Alternatives Considered
1. **Generic query builder**: Too complex for current needs, over-engineering
2. **Template-based approach**: Less type-safe than direct Go code generation
3. **Different backend identifier**: Considered "rapid7" and "insightidr", selected "r7" for brevity

## References
- Rapid7 InsightIDR Advanced Query Language (AQL) documentation
- Existing S1QL backend implementation: `pkg/backends/s1ql/`
- Backend interface definition: `pkg/backends/interface.go`
