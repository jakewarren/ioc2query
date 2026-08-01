// Package s1ql provides SentinelOne S1QLv2 query generation from IOCs.
// It implements the backends.Backend interface to generate queries
// compatible with SentinelOne's query language.
package s1ql

import (
	"fmt"
	"os"
	"strings"

	"github.com/jakewarren/ioc2query/pkg/extraction"
)

// S1QLBackend generates SentinelOne S1QLv2 queries
type S1QLBackend struct {
	config *Config
}

// Config holds S1QL backend configuration
type Config struct {
	CombineWithOR   bool // Combine IOCs in single query vs separate
	IncludeComments bool // Add comments explaining query parts
}

// New creates a new S1QL backend with the given configuration
func New(config *Config) *S1QLBackend {
	if config == nil {
		config = &Config{
			CombineWithOR:   true,
			IncludeComments: false,
		}
	}
	return &S1QLBackend{
		config: config,
	}
}

// Name returns the backend identifier
func (b *S1QLBackend) Name() string {
	return "s1"
}

// GenerateQuery creates queries grouped by IOC category
func (b *S1QLBackend) GenerateQuery(iocs *extraction.IOCSet) ([]string, error) {
	if iocs == nil || iocs.IsEmpty() {
		return nil, fmt.Errorf("IOC set is empty")
	}

	var fileParts []string
	if len(iocs.MD5Hashes) > 0 {
		fileParts = append(fileParts, b.generateMD5Query(iocs.MD5Hashes))
	}
	if len(iocs.SHA1Hashes) > 0 {
		fileParts = append(fileParts, b.generateSHA1Query(iocs.SHA1Hashes))
	}
	if len(iocs.SHA256Hashes) > 0 {
		fileParts = append(fileParts, b.generateSHA256Query(iocs.SHA256Hashes))
	}

	var networkParts []string
	if len(iocs.Domains) > 0 {
		networkParts = append(networkParts, b.generateDomainQuery(iocs.Domains))
	}
	if len(iocs.IPv4Addresses) > 0 {
		networkParts = append(networkParts, b.generateIPQuery(iocs.IPv4Addresses))
	}

	var groups []string
	if len(fileParts) > 0 {
		groups = append(groups, strings.Join(fileParts, " || "))
	}
	if len(networkParts) > 0 {
		groups = append(groups, strings.Join(networkParts, " || "))
	}

	// Warn if query is very large
	if iocs.Count() > 1000 {
		fmt.Fprintf(os.Stderr, "Warning: Query contains %d IOCs, which may be very large\n", iocs.Count())
	}

	// SentinelOne supports consolidated queries across all IOC types
	return []string{strings.Join(groups, " ||\n")}, nil
}

// GenerateQueries creates one query per IOC type
func (b *S1QLBackend) GenerateQueries(iocs *extraction.IOCSet) ([]string, error) {
	if iocs == nil || iocs.IsEmpty() {
		return nil, fmt.Errorf("IOC set is empty")
	}

	var queries []string

	if len(iocs.MD5Hashes) > 0 {
		queries = append(queries, b.generateMD5Query(iocs.MD5Hashes))
	}
	if len(iocs.SHA1Hashes) > 0 {
		queries = append(queries, b.generateSHA1Query(iocs.SHA1Hashes))
	}
	if len(iocs.SHA256Hashes) > 0 {
		queries = append(queries, b.generateSHA256Query(iocs.SHA256Hashes))
	}
	if len(iocs.Domains) > 0 {
		queries = append(queries, b.generateDomainQuery(iocs.Domains))
	}
	if len(iocs.IPv4Addresses) > 0 {
		queries = append(queries, b.generateIPQuery(iocs.IPv4Addresses))
	}

	return queries, nil
}

// generateAnyFieldQuery builds an S1QL power-query clause that searches
// multiple fields for the same set of values, e.g. any(f1, f2) in (v1, v2).
func (b *S1QLBackend) generateAnyFieldQuery(fields []string, values []string) string {
	return fmt.Sprintf(`any(%s) in (%s)`, strings.Join(fields, ", "), b.formatStringList(values))
}

// generateMD5Query creates a query for MD5 hashes
func (b *S1QLBackend) generateMD5Query(hashes []string) string {
	return b.generateAnyFieldQuery([]string{"src.process.image.md5", "tgt.file.md5"}, hashes)
}

// generateSHA1Query creates a query for SHA1 hashes
func (b *S1QLBackend) generateSHA1Query(hashes []string) string {
	return b.generateAnyFieldQuery([]string{"src.process.image.sha1", "tgt.file.sha1"}, hashes)
}

// generateSHA256Query creates a query for SHA256 hashes
func (b *S1QLBackend) generateSHA256Query(hashes []string) string {
	return b.generateAnyFieldQuery([]string{"src.process.image.sha256", "tgt.file.sha256"}, hashes)
}

// generateDomainQuery creates a query for domains
func (b *S1QLBackend) generateDomainQuery(domains []string) string {
	domainList := b.formatStringList(domains)
	return fmt.Sprintf(`event.dns.request in (%s)`, domainList)
}

// generateIPQuery creates a query for IP addresses
func (b *S1QLBackend) generateIPQuery(ips []string) string {
	return b.generateAnyFieldQuery([]string{"src.ip.address", "dst.ip.address"}, ips)
}

// formatStringList formats a list of strings for S1QL IN clauses
func (b *S1QLBackend) formatStringList(items []string) string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf(`"%s"`, b.escapeString(item))
	}
	return strings.Join(quoted, ", ")
}

// escapeString escapes special characters in strings for S1QL
func (b *S1QLBackend) escapeString(s string) string {
	// Escape backslashes and quotes
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
