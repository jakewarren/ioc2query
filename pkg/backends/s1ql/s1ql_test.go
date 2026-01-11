package s1ql

import (
	"strings"
	"testing"

	"github.com/jakewarren/ioc2query/pkg/extraction"
)

func TestS1QLBackend_Name(t *testing.T) {
	backend := New(nil)
	if got := backend.Name(); got != "s1" {
		t.Errorf("Name() = %v, want %v", got, "s1")
	}
}

func TestS1QLBackend_GenerateQuery(t *testing.T) {
	tests := []struct {
		name    string
		iocs    *extraction.IOCSet
		want    string
		wantErr bool
	}{
		{
			name: "single MD5 hash",
			iocs: &extraction.IOCSet{
				MD5Hashes: []string{"5d41402abc4b2a76b9719d911017c592"},
			},
			want:    `src.process.image.md5 IN ("5d41402abc4b2a76b9719d911017c592") || tgt.file.md5 IN ("5d41402abc4b2a76b9719d911017c592")`,
			wantErr: false,
		},
		{
			name: "multiple MD5 hashes",
			iocs: &extraction.IOCSet{
				MD5Hashes: []string{"hash1", "hash2"},
			},
			want:    `src.process.image.md5 IN ("hash1", "hash2") || tgt.file.md5 IN ("hash1", "hash2")`,
			wantErr: false,
		},
		{
			name: "single SHA1 hash",
			iocs: &extraction.IOCSet{
				SHA1Hashes: []string{"356a192b7913b04c54574d18c28d46e6395428ab"},
			},
			want:    `src.process.image.sha1 IN ("356a192b7913b04c54574d18c28d46e6395428ab") || tgt.file.sha1 IN ("356a192b7913b04c54574d18c28d46e6395428ab")`,
			wantErr: false,
		},
		{
			name: "multiple SHA256 hashes",
			iocs: &extraction.IOCSet{
				SHA256Hashes: []string{"hash1", "hash2", "hash3"},
			},
			want:    `src.process.image.sha256 IN ("hash1", "hash2", "hash3") || tgt.file.sha256 IN ("hash1", "hash2", "hash3")`,
			wantErr: false,
		},
		{
			name: "single domain",
			iocs: &extraction.IOCSet{
				Domains: []string{"malicious.com"},
			},
			want:    `event.dns.request IN ("malicious.com")`,
			wantErr: false,
		},
		{
			name: "multiple domains",
			iocs: &extraction.IOCSet{
				Domains: []string{"evil.com", "bad.net"},
			},
			want:    `event.dns.request IN ("evil.com", "bad.net")`,
			wantErr: false,
		},
		{
			name: "single IP",
			iocs: &extraction.IOCSet{
				IPv4Addresses: []string{"192.168.1.1"},
			},
			want:    `(src.ip.address = "192.168.1.1" || dst.ip.address = "192.168.1.1")`,
			wantErr: false,
		},
		{
			name: "multiple IPs",
			iocs: &extraction.IOCSet{
				IPv4Addresses: []string{"10.0.0.1", "172.16.0.1"},
			},
			want:    `(src.ip.address IN ("10.0.0.1", "172.16.0.1") || dst.ip.address IN ("10.0.0.1", "172.16.0.1"))`,
			wantErr: false,
		},
		{
			name: "combined query with all IOC types",
			iocs: &extraction.IOCSet{
				MD5Hashes:     []string{"md5hash"},
				SHA1Hashes:    []string{"sha1hash"},
				SHA256Hashes:  []string{"sha256hash"},
				Domains:       []string{"evil.com"},
				IPv4Addresses: []string{"1.2.3.4"},
			},
			want: `(src.process.image.md5 IN ("md5hash") || tgt.file.md5 IN ("md5hash") || src.process.image.sha1 IN ("sha1hash") || tgt.file.sha1 IN ("sha1hash") || src.process.image.sha256 IN ("sha256hash") || tgt.file.sha256 IN ("sha256hash")) ||
(event.dns.request IN ("evil.com") || (src.ip.address = "1.2.3.4" || dst.ip.address = "1.2.3.4"))`,
			wantErr: false,
		},
		{
			name:    "empty IOC set",
			iocs:    &extraction.IOCSet{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "nil IOC set",
			iocs:    nil,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := New(nil)
			got, err := backend.GenerateQuery(tt.iocs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS1QLBackend_GenerateQueries(t *testing.T) {
	tests := []struct {
		name    string
		iocs    *extraction.IOCSet
		want    []string
		wantErr bool
	}{
		{
			name: "separate queries for each indicator",
			iocs: &extraction.IOCSet{
				MD5Hashes:     []string{"md5_1", "md5_2"},
				Domains:       []string{"evil.com"},
				IPv4Addresses: []string{"1.2.3.4"},
			},
			want: []string{
				`src.process.image.md5 IN ("md5_1") || tgt.file.md5 IN ("md5_1")`,
				`src.process.image.md5 IN ("md5_2") || tgt.file.md5 IN ("md5_2")`,
				`event.dns.request IN ("evil.com")`,
				`(src.ip.address = "1.2.3.4" || dst.ip.address = "1.2.3.4")`,
			},
			wantErr: false,
		},
		{
			name: "only hash queries",
			iocs: &extraction.IOCSet{
				MD5Hashes:    []string{"hash1", "hash2"},
				SHA256Hashes: []string{"hash3"},
			},
			want: []string{
				`src.process.image.md5 IN ("hash1") || tgt.file.md5 IN ("hash1")`,
				`src.process.image.md5 IN ("hash2") || tgt.file.md5 IN ("hash2")`,
				`src.process.image.sha256 IN ("hash3") || tgt.file.sha256 IN ("hash3")`,
			},
			wantErr: false,
		},
		{
			name:    "empty IOC set",
			iocs:    &extraction.IOCSet{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := New(nil)
			got, err := backend.GenerateQueries(tt.iocs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQueries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("GenerateQueries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestS1QLBackend_escapeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no special characters",
			input: "example.com",
			want:  "example.com",
		},
		{
			name:  "with quotes",
			input: `domain"with"quotes`,
			want:  `domain\"with\"quotes`,
		},
		{
			name:  "with backslashes",
			input: `domain\with\backslash`,
			want:  `domain\\with\\backslash`,
		},
		{
			name:  "with both quotes and backslashes",
			input: `domain\"test`,
			want:  `domain\\\"test`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			backend := New(nil)
			got := backend.escapeString(tt.input)
			if got != tt.want {
				t.Errorf("escapeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		// Normalize whitespace for comparison
		aNorm := strings.TrimSpace(a[i])
		bNorm := strings.TrimSpace(b[i])
		if aNorm != bNorm {
			return false
		}
	}
	return true
}
