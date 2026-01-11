# Change: Add S1 Query Support

## Why
Security analysts need to rapidly transform IOCs into SentinelOne-compatible detection queries for threat hunting across endpoints. Manual query writing is time-consuming and error-prone.

## What Changes
- Add IOC extraction capability using `github.com/vertoforce/go-ioc` library
- Implement S1QLv2 query generation backend for file hashes, domains, and IP addresses
- Create CLI interface with stdin/file input and flexible output options
- Provide comprehensive testing suite with realistic IOC samples

## Impact
- **Affected specs:** Creates 3 new capabilities: `ioc-extraction`, `s1ql-backend`, `cli-interface`
- **Affected code:** New Go modules in `pkg/extraction/`, `pkg/backends/s1ql/`, and `cmd/ioc2query/`
- **Dependencies:** Adds `github.com/vertoforce/go-ioc` (MIT license)
- **Users:** Enables incident responders to generate S1QL queries in seconds vs minutes
