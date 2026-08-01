// Package backends defines the interface for query generation backends.
// Each backend implements vendor-specific query generation from extracted IOCs.
package backends

import "github.com/jakewarren/ioc2query/pkg/extraction"

// Backend generates vendor-specific queries from IOCs
type Backend interface {
	// Name returns the backend identifier (e.g., "s1", "rapid7")
	Name() string

	// GenerateQuery creates queries grouped by IOC type
	GenerateQuery(iocs *extraction.IOCSet) ([]string, error)

	// GenerateQueries creates one query per IOC type
	GenerateQueries(iocs *extraction.IOCSet) ([]string, error)
}
