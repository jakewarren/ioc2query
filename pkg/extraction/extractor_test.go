package extraction

import (
	"testing"
)

func TestExtractor_Extract(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *IOCSet
		wantErr bool
	}{
		{
			name: "extract mixed IOCs",
			input: `
			MD5: 5d41402abc4b2a76b9719d911017c592
			SHA1: 356a192b7913b04c54574d18c28d46e6395428ab
			SHA256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
			Domain: malicious.com
			IP: 192.168.1.100
			`,
			want: &IOCSet{
				MD5Hashes:     []string{"5d41402abc4b2a76b9719d911017c592"},
				SHA1Hashes:    []string{"356a192b7913b04c54574d18c28d46e6395428ab"},
				SHA256Hashes:  []string{"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
				Domains:       []string{"malicious.com"},
				IPv4Addresses: []string{"192.168.1.100"},
			},
			wantErr: false,
		},
		{
			name: "deduplicate IOCs",
			input: `
			MD5: 5d41402abc4b2a76b9719d911017c592
			MD5: 5D41402ABC4B2A76B9719D911017C592
			Domain: example.com
			Domain: example.com
			IP: 192.168.1.1
			IP: 192.168.1.1
			`,
			want: &IOCSet{
				MD5Hashes:     []string{"5d41402abc4b2a76b9719d911017c592"},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{"example.com"},
				IPv4Addresses: []string{"192.168.1.1"},
			},
			wantErr: false,
		},
		{
			name:    "error on empty input",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "error on whitespace-only input",
			input:   "   \n\t  ",
			want:    nil,
			wantErr: true,
		},
		{
			name: "handle no IOCs found",
			input: `
			This is just some text with no IOCs.
			Nothing to see here.
			`,
			want: &IOCSet{
				MD5Hashes:     []string{},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{},
				IPv4Addresses: []string{},
			},
			wantErr: false,
		},
		{
			name: "extract only hashes",
			input: `
			MD5: 098f6bcd4621d373cade4e832627b4f6
			SHA1: a94a8fe5ccb19ba61c4c0873d391e987982fbbd3
			SHA256: 9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08
			`,
			want: &IOCSet{
				MD5Hashes:     []string{"098f6bcd4621d373cade4e832627b4f6"},
				SHA1Hashes:    []string{"a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"},
				SHA256Hashes:  []string{"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"},
				Domains:       []string{},
				IPv4Addresses: []string{},
			},
			wantErr: false,
		},
		{
			name: "extract only network indicators",
			input: `
			Domain: evil.com
			Domain: badsite.net
			IP: 10.0.0.1
			IP: 172.16.0.1
			`,
			want: &IOCSet{
				MD5Hashes:     []string{},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{"evil.com", "badsite.net"},
				IPv4Addresses: []string{"10.0.0.1", "172.16.0.1"},
			},
			wantErr: false,
		},
		{
			name: "extract defanged IOCs",
			input: `
			Domain: evil[.]com
			Domain: badsite(dot)net
			IP: 192[.]168[.]1[.]1
			URL: hxxp://malicious.example[.]com/path
			`,
			want: &IOCSet{
				MD5Hashes:     []string{},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{"evil.com", "badsite.net", "malicious.example.com"},
				IPv4Addresses: []string{"192.168.1.1"},
			},
			wantErr: false,
		},
		{
			name: "ignore invalid domains",
			input: `
			Invalid domain: example.c
			Not a domain: test.123
			`,
			want: &IOCSet{
				MD5Hashes:     []string{},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{},
				IPv4Addresses: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := New()
			got, err := e.Extract(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extractor.Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Compare results
			if !stringSliceEqual(got.MD5Hashes, tt.want.MD5Hashes) {
				t.Errorf("MD5Hashes = %v, want %v", got.MD5Hashes, tt.want.MD5Hashes)
			}
			if !stringSliceEqual(got.SHA1Hashes, tt.want.SHA1Hashes) {
				t.Errorf("SHA1Hashes = %v, want %v", got.SHA1Hashes, tt.want.SHA1Hashes)
			}
			if !stringSliceEqual(got.SHA256Hashes, tt.want.SHA256Hashes) {
				t.Errorf("SHA256Hashes = %v, want %v", got.SHA256Hashes, tt.want.SHA256Hashes)
			}
			if !stringSliceEqual(got.Domains, tt.want.Domains) {
				t.Errorf("Domains = %v, want %v", got.Domains, tt.want.Domains)
			}
			if !stringSliceEqual(got.IPv4Addresses, tt.want.IPv4Addresses) {
				t.Errorf("IPv4Addresses = %v, want %v", got.IPv4Addresses, tt.want.IPv4Addresses)
			}
		})
	}
}

func TestIOCSet_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		iocs *IOCSet
		want bool
	}{
		{
			name: "empty IOCSet",
			iocs: &IOCSet{
				MD5Hashes:     []string{},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{},
				IPv4Addresses: []string{},
			},
			want: true,
		},
		{
			name: "non-empty IOCSet",
			iocs: &IOCSet{
				MD5Hashes:     []string{"hash"},
				SHA1Hashes:    []string{},
				SHA256Hashes:  []string{},
				Domains:       []string{},
				IPv4Addresses: []string{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.iocs.IsEmpty(); got != tt.want {
				t.Errorf("IOCSet.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIOCSet_Count(t *testing.T) {
	tests := []struct {
		name string
		iocs *IOCSet
		want int
	}{
		{
			name: "empty IOCSet",
			iocs: &IOCSet{},
			want: 0,
		},
		{
			name: "single IOC",
			iocs: &IOCSet{
				MD5Hashes: []string{"hash"},
			},
			want: 1,
		},
		{
			name: "multiple IOCs",
			iocs: &IOCSet{
				MD5Hashes:     []string{"hash1", "hash2"},
				Domains:       []string{"example.com"},
				IPv4Addresses: []string{"1.2.3.4", "5.6.7.8"},
			},
			want: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.iocs.Count(); got != tt.want {
				t.Errorf("IOCSet.Count() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices (order-independent)
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	
	// Create maps for order-independent comparison
	aMap := make(map[string]bool)
	bMap := make(map[string]bool)
	
	for _, v := range a {
		aMap[v] = true
	}
	for _, v := range b {
		bMap[v] = true
	}
	
	// Check if all elements match
	for k := range aMap {
		if !bMap[k] {
			return false
		}
	}
	for k := range bMap {
		if !aMap[k] {
			return false
		}
	}
	
	return true
}
