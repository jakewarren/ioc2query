// Package extraction provides IOC (Indicators of Compromise) extraction
// from raw text input. It uses the go-ioc library to identify and extract
// various types of IOCs including file hashes, domains, and IP addresses.
package extraction

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vertoforce/go-ioc/ioc"
)

// IOCSet holds deduplicated IOCs by type
type IOCSet struct {
	MD5Hashes     []string
	SHA1Hashes    []string
	SHA256Hashes  []string
	Domains       []string
	IPv4Addresses []string
}

// Extractor wraps go-ioc functionality
type Extractor struct{}

// New creates a new IOC extractor
func New() *Extractor {
	return &Extractor{}
}

// Extract returns deduplicated IOCs from input text
func (e *Extractor) Extract(input string) (*IOCSet, error) {
	if strings.TrimSpace(input) == "" {
		return nil, fmt.Errorf("input is empty")
	}

	// Pre-process input to refang defanged indicators
	input = refang(input)

	iocSet := &IOCSet{
		MD5Hashes:     make([]string, 0),
		SHA1Hashes:    make([]string, 0),
		SHA256Hashes:  make([]string, 0),
		Domains:       make([]string, 0),
		IPv4Addresses: make([]string, 0),
	}

	// Extract IOCs using go-ioc
	parsedIOCs := ioc.GetIOCs(input, true)

	// Deduplicate and normalize MD5 hashes
	md5Map := make(map[string]bool)
	sha1Map := make(map[string]bool)
	sha256Map := make(map[string]bool)
	domainMap := make(map[string]bool)
	ipMap := make(map[string]bool)

	for _, parsedIOC := range parsedIOCs {
		switch parsedIOC.Type {
		case ioc.MD5:
			normalized := strings.ToLower(parsedIOC.IOC)
			if !md5Map[normalized] {
				md5Map[normalized] = true
				iocSet.MD5Hashes = append(iocSet.MD5Hashes, normalized)
			}
		case ioc.SHA1:
			normalized := strings.ToLower(parsedIOC.IOC)
			if !sha1Map[normalized] {
				sha1Map[normalized] = true
				iocSet.SHA1Hashes = append(iocSet.SHA1Hashes, normalized)
			}
		case ioc.SHA256:
			normalized := strings.ToLower(parsedIOC.IOC)
			if !sha256Map[normalized] {
				sha256Map[normalized] = true
				iocSet.SHA256Hashes = append(iocSet.SHA256Hashes, normalized)
			}
		case ioc.Domain:
			domain := strings.ToLower(parsedIOC.IOC)
			if !domainMap[domain] {
				domainMap[domain] = true
				iocSet.Domains = append(iocSet.Domains, domain)
			}
		case ioc.IPv4:
			if !ipMap[parsedIOC.IOC] {
				ipMap[parsedIOC.IOC] = true
				iocSet.IPv4Addresses = append(iocSet.IPv4Addresses, parsedIOC.IOC)
			}
		}
	}

	return iocSet, nil
}

// IsEmpty returns true if the IOCSet contains no IOCs
func (iocs *IOCSet) IsEmpty() bool {
	return len(iocs.MD5Hashes) == 0 &&
		len(iocs.SHA1Hashes) == 0 &&
		len(iocs.SHA256Hashes) == 0 &&
		len(iocs.Domains) == 0 &&
		len(iocs.IPv4Addresses) == 0
}

// Count returns the total number of IOCs in the set
func (iocs *IOCSet) Count() int {
	return len(iocs.MD5Hashes) +
		len(iocs.SHA1Hashes) +
		len(iocs.SHA256Hashes) +
		len(iocs.Domains) +
		len(iocs.IPv4Addresses)
}

// refang converts defanged IOCs to their standard format
func refang(input string) string {
	// Handle [.] and (.) and [dot] and (dot)
	input = strings.ReplaceAll(input, "[.]", ".")
	input = strings.ReplaceAll(input, "(.)", ".")
	input = strings.ReplaceAll(input, "[dot]", ".")
	input = strings.ReplaceAll(input, "(dot)", ".")

	// Handle hxxp
	reHxxp := regexp.MustCompile(`(?i)hxxp`)
	input = reHxxp.ReplaceAllString(input, "http")

	return input
}
