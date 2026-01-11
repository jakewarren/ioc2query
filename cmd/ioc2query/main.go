package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jakewarren/ioc2query/pkg/backends"
	"github.com/jakewarren/ioc2query/pkg/backends/s1ql"
	"github.com/jakewarren/ioc2query/pkg/extraction"
)

const (
	exitSuccess         = 0
	exitIOError         = 1
	exitInvalidArgs     = 2
	exitExtractionError = 3
	exitQueryGenError   = 4
)

var (
	backendFlag  = flag.String("backend", "", "Backend to use (s1)")
	s1Flag       = flag.Bool("s1", false, "Shortcut for --backend s1")
	inputFlag    = flag.String("input", "", "Input file (default: stdin)")
	outputFlag   = flag.String("output", "", "Output file (default: stdout)")
	separateFlag = flag.Bool("separate", false, "Generate separate queries per IOC type")
	verboseFlag  = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
	flag.StringVar(backendFlag, "b", "", "Backend to use (s1)")
	flag.StringVar(inputFlag, "i", "", "Input file (default: stdin)")
	flag.StringVar(outputFlag, "o", "", "Output file (default: stdout)")
	flag.BoolVar(separateFlag, "s", false, "Generate separate queries per IOC type")
	flag.BoolVar(verboseFlag, "v", false, "Enable verbose logging")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(determineExitCode(err))
	}
}

func run() error {
	// Handle shortcuts
	if *s1Flag {
		*backendFlag = "s1"
	}

	// Validate backend flag is provided
	if *backendFlag == "" {
		return &InvalidArgsError{fmt.Errorf("backend required (use --s1 or --backend s1)")}
	}

	// Read input
	input, err := readInput()
	if err != nil {
		return &IOError{err}
	}

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "Extracting IOCs from input...\n")
	}

	// Extract IOCs
	extractor := extraction.New()
	iocs, err := extractor.Extract(input)
	if err != nil {
		return &ExtractionError{err}
	}

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "Extracted %d IOCs:\n", iocs.Count())
		fmt.Fprintf(os.Stderr, "  MD5: %d\n", len(iocs.MD5Hashes))
		fmt.Fprintf(os.Stderr, "  SHA1: %d\n", len(iocs.SHA1Hashes))
		fmt.Fprintf(os.Stderr, "  SHA256: %d\n", len(iocs.SHA256Hashes))
		fmt.Fprintf(os.Stderr, "  Domains: %d\n", len(iocs.Domains))
		fmt.Fprintf(os.Stderr, "  IPs: %d\n", len(iocs.IPv4Addresses))
	}

	if iocs.IsEmpty() {
		return &ExtractionError{fmt.Errorf("no IOCs found in input")}
	}

	// Select backend
	backend, err := selectBackend(*backendFlag)
	if err != nil {
		return &InvalidArgsError{err}
	}

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "Generating queries using %s backend...\n", backend.Name())
	}

	// Generate queries
	var output string
	if *separateFlag {
		queries, err := backend.GenerateQueries(iocs)
		if err != nil {
			return &QueryGenError{err}
		}
		for i, query := range queries {
			if i > 0 {
				output += "\n\n"
			}
			output += query
		}
	} else {
		query, err := backend.GenerateQuery(iocs)
		if err != nil {
			return &QueryGenError{err}
		}
		output = query
	}

	// Write output
	if err := writeOutput(output); err != nil {
		return &IOError{err}
	}

	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "Query generation complete.\n")
	}

	return nil
}

func readInput() (string, error) {
	var reader io.Reader
	if *inputFlag != "" {
		file, err := os.Open(*inputFlag)
		if err != nil {
			return "", fmt.Errorf("failed to open input file: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close input file: %v\n", err)
			}
		}()
		reader = file
	} else {
		reader = os.Stdin
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	return string(data), nil
}

func writeOutput(output string) error {
	var writer io.Writer
	if *outputFlag != "" {
		file, err := os.Create(*outputFlag)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to close output file: %v\n", err)
			}
		}()
		writer = file
	} else {
		writer = os.Stdout
	}

	_, err := fmt.Fprint(writer, output)
	if err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	// Add newline if writing to stdout
	if *outputFlag == "" {
		if _, err := fmt.Fprintln(writer); err != nil {
			return fmt.Errorf("failed to write newline: %w", err)
		}
	}

	return nil
}

func selectBackend(name string) (backends.Backend, error) {
	name = strings.ToLower(name)
	switch name {
	case "s1":
		return s1ql.New(nil), nil
	default:
		return nil, fmt.Errorf("unknown backend: %s (supported: s1)", name)
	}
}

func determineExitCode(err error) int {
	switch err.(type) {
	case *IOError:
		return exitIOError
	case *InvalidArgsError:
		return exitInvalidArgs
	case *ExtractionError:
		return exitExtractionError
	case *QueryGenError:
		return exitQueryGenError
	default:
		return exitIOError
	}
}

// Error types for different failure modes
type IOError struct{ err error }

func (e *IOError) Error() string { return e.err.Error() }

type InvalidArgsError struct{ err error }

func (e *InvalidArgsError) Error() string { return e.err.Error() }

type ExtractionError struct{ err error }

func (e *ExtractionError) Error() string { return e.err.Error() }

type QueryGenError struct{ err error }

func (e *QueryGenError) Error() string { return e.err.Error() }
