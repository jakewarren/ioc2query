// Package r7 provides Rapid7 InsightIDR LEQL query generation from IOCs.
// It implements the backends.Backend interface to generate queries
// compatible with Rapid7's Log Entry Query Language.
package r7

import (
	"fmt"
	"os"
	"strings"

	"github.com/jakewarren/ioc2query/pkg/backends"
	"github.com/jakewarren/ioc2query/pkg/extraction"
)

// Ensure R7Backend implements the Backend interface
var _ backends.Backend = (*R7Backend)(nil)

// R7Backend generates Rapid7 InsightIDR LEQL queries
type R7Backend struct {
	config *Config
}

// Config holds R7 backend configuration
type Config struct {
	// Reserved for future extensibility
}

// New creates a new R7 backend with the given configuration
func New(config *Config) *R7Backend {
	if config == nil {
		config = &Config{}
	}
	return &R7Backend{
		config: config,
	}
}

// Name returns the backend identifier
func (b *R7Backend) Name() string {
	return "r7"
}

// GenerateQuery creates queries grouped by IOC type
func (b *R7Backend) GenerateQuery(iocs *extraction.IOCSet) ([]string, error) {
	if iocs == nil || iocs.IsEmpty() {
		return nil, fmt.Errorf("IOC set is empty")
	}

	var queries []string

	// Combine all hash types into a single query
	var allHashes []string
	allHashes = append(allHashes, iocs.MD5Hashes...)
	allHashes = append(allHashes, iocs.SHA1Hashes...)
	allHashes = append(allHashes, iocs.SHA256Hashes...)
	if len(allHashes) > 0 {
		queries = append(queries, b.generateHashQuery(allHashes))
	}

	// Generate domain query
	if len(iocs.Domains) > 0 {
		queries = append(queries, b.generateDomainQuery(iocs.Domains))
	}

	// Generate IP query
	if len(iocs.IPv4Addresses) > 0 {
		queries = append(queries, b.generateIPQuery(iocs.IPv4Addresses))
	}

	// Warn if query is very large
	if iocs.Count() > 1000 {
		fmt.Fprintf(os.Stderr, "Warning: Query contains %d IOCs, which may be very large\n", iocs.Count())
	}

	return queries, nil
}

// GenerateQueries creates one query per IOC type
func (b *R7Backend) GenerateQueries(iocs *extraction.IOCSet) ([]string, error) {
	if iocs == nil || iocs.IsEmpty() {
		return nil, fmt.Errorf("IOC set is empty")
	}

	var queries []string

	var allHashes []string
	allHashes = append(allHashes, iocs.MD5Hashes...)
	allHashes = append(allHashes, iocs.SHA1Hashes...)
	allHashes = append(allHashes, iocs.SHA256Hashes...)
	if len(allHashes) > 0 {
		queries = append(queries, b.generateHashQuery(allHashes))
	}

	if len(iocs.Domains) > 0 {
		queries = append(queries, b.generateDomainQuery(iocs.Domains))
	}

	if len(iocs.IPv4Addresses) > 0 {
		queries = append(queries, b.generateIPQuery(iocs.IPv4Addresses))
	}

	return queries, nil
}

// generateHashQuery creates a query for file hashes
// Uses wildcard pattern to match all hash types
func (b *R7Backend) generateHashQuery(hashes []string) string {
	hashList := b.formatStringList(hashes)
	return fmt.Sprintf(`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN [%s])`, hashList)
}

// generateDomainQuery creates a query for domains
// Uses ICONTAINS-ANY for case-insensitive partial matching
func (b *R7Backend) generateDomainQuery(domains []string) string {
	domainList := b.formatStringList(domains)
	return fmt.Sprintf(`where("query","url" ICONTAINS-ANY [%s])`, domainList)
}

// generateIPQuery creates a query for IP addresses
// Uses IN operator for exact matching on source and destination
func (b *R7Backend) generateIPQuery(ips []string) string {
	ipList := b.formatStringList(ips)
	return fmt.Sprintf(`where("source_address","destination_address" IN [%s])`, ipList)
}

// formatStringList formats a list of strings for R7 LEQL array syntax
func (b *R7Backend) formatStringList(items []string) string {
	quoted := make([]string, len(items))
	for i, item := range items {
		quoted[i] = fmt.Sprintf(`'%s'`, b.escapeString(item))
	}
	return strings.Join(quoted, ", ")
}

// escapeString escapes special characters in strings for R7 LEQL
func (b *R7Backend) escapeString(s string) string {
	// Escape backslashes and single quotes
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `'`, `\'`)
	return s
}
