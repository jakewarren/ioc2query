# ioc2query

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/jakewarren/ioc2query/blob/master/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=shields)](http://makeapullrequest.com)
[![Hosted on GitHub Pages](https://img.shields.io/badge/hosted%20on-GitHub%20Pages-222222?logo=githubpages&logoColor=white)](https://jakewarren.github.io/ioc2query/)


Transform Indicators of Compromise (IOCs) into vendor-specific detection queries for threat hunting.

## Overview

`ioc2query` is a command-line tool that extracts IOCs from raw text (reports, logs, threat feeds) and automatically generates detection queries for security platforms. Currently supports SentinelOne S1QL query language.

## Features

- **Automatic IOC Extraction**: Parses MD5, SHA1, SHA256 hashes, domains, and IPv4 addresses from any text
- **Deduplication**: Automatically removes duplicate IOCs
- **Web Interface**: Browser-based version powered by WebAssembly
- **S1QL Query Generation**: Creates SentinelOne-compatible S1QLv2 queries
- **Flexible Input**: Accepts stdin or file input
- **Multiple Output Modes**: Combined queries or separate queries per IOC type

## Web Interface

`ioc2query` is available as a client-side web application using WebAssembly. This allows you to generate queries directly in your browser without installing any tools. No data leaves your machine.

**[Try the Web Interface](https://jakewarren.github.io/ioc2query/)**

### Local Development of Web Interface

To run the web interface locally, you need a simple static web server.

1. Build the WASM binary:
   ```bash
   make wasm
   ```
2. Serve the `web` directory:
   ```bash
   # Using Python
   python3 -m http.server -d web 8080
   ```
3. Open `http://localhost:8080` in your browser.


## Installation

### From Source

```bash
go install github.com/jakewarren/ioc2query/cmd/ioc2query@latest
```

### Build Locally

```bash
git clone https://github.com/jakewarren/ioc2query.git
cd ioc2query
go build -o ioc2query ./cmd/ioc2query
```

## Usage

### Basic Usage

```bash
# From stdin
cat iocs.txt | ioc2query --backend s1

# From file
ioc2query --backend s1 --input iocs.txt

# Save to file
ioc2query -b s1 -i iocs.txt -o query.txt
```

### Options

```
  -b, --backend string   Backend to use (required, currently supports: s1)
  -i, --input string     Input file (default: stdin)
  -o, --output string    Output file (default: stdout)
  -s, --separate         Generate separate queries per IOC type
  -v, --verbose          Enable verbose logging
  -h, --help             Show usage information
```

### Examples

**Extract IOCs and generate S1QL query:**

```bash
$ echo "Observed malicious activity from 192.168.1.100 connecting to evil.com.
File hash: 5d41402abc4b2a76b9719d911017c592" | ioc2query -b s1

src.process.image.md5 IN ("5d41402abc4b2a76b9719d911017c592") || tgt.file.md5 IN ("5d41402abc4b2a76b9719d911017c592") OR
event.dns.request IN ("evil.com") OR
(src.ip.address = "192.168.1.100" || dst.ip.address = "192.168.1.100")
```

**Generate separate queries:**

```bash
$ ioc2query -b s1 -i threat_report.txt --separate

src.process.image.md5 IN ("hash1", "hash2", "hash3") || tgt.file.md5 IN ("hash1", "hash2", "hash3")

src.process.image.sha256 IN ("hash4", "hash5") || tgt.file.sha256 IN ("hash4", "hash5")

event.dns.request IN ("malicious.com", "evil.net")

(src.ip.address IN ("1.2.3.4", "5.6.7.8") || dst.ip.address IN ("1.2.3.4", "5.6.7.8"))
```

**Verbose output:**

```bash
$ ioc2query -b s1 -i iocs.txt -v
Extracting IOCs from input...
Extracted 42 IOCs:
  MD5: 12
  SHA1: 8
  SHA256: 15
  Domains: 5
  IPs: 2
Generating queries using s1 backend...
[query output]
Query generation complete.
```

## Supported IOC Types

| Type   | Description      | Example                                                            |
| ------ | ---------------- | ------------------------------------------------------------------ |
| MD5    | MD5 file hash    | `5d41402abc4b2a76b9719d911017c592`                                 |
| SHA1   | SHA1 file hash   | `356a192b7913b04c54574d18c28d46e6395428ab`                         |
| SHA256 | SHA256 file hash | `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` |
| Domain | Domain names     | `malicious.com`, `evil.net`                                        |
| IPv4   | IPv4 addresses   | `192.168.1.100`, `10.0.0.1`                                        |

## S1QL Query Patterns

### File Hashes

```sql
-- Single MD5 hash
src.process.image.md5 IN ("5d41402abc4b2a76b9719d911017c592") || tgt.file.md5 IN ("5d41402abc4b2a76b9719d911017c592")

-- Multiple hashes
src.process.image.md5 IN ("hash1", "hash2", "hash3") || tgt.file.md5 IN ("hash1", "hash2", "hash3")
```

### Domains

```sql
-- Single domain
event.dns.request IN ("malicious.com")

-- Multiple domains
event.dns.request IN ("evil.com", "bad.net", "malware.org")
```

### IP Addresses

```sql
-- Single IP (checks both source and destination)
(src.ip.address = "192.168.1.1" || dst.ip.address = "192.168.1.1")

-- Multiple IPs
(src.ip.address IN ("10.0.0.1", "172.16.0.1") || dst.ip.address IN ("10.0.0.1", "172.16.0.1"))
```


## License

MIT License - see LICENSE file for details

## Acknowledgments

- [go-ioc](https://github.com/vertoforce/go-ioc) for IOC extraction
