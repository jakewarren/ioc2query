package r7

import (
	"strings"
	"testing"

	"github.com/jakewarren/ioc2query/pkg/extraction"
)

func TestR7Backend_Name(t *testing.T) {
	backend := New(nil)
	if got := backend.Name(); got != "r7" {
		t.Errorf("Name() = %v, want %v", got, "r7")
	}
}

func TestR7Backend_generateHashQuery(t *testing.T) {
	backend := New(nil)
	tests := []struct {
		name   string
		hashes []string
		want   string
	}{
		{
			name:   "single MD5 hash",
			hashes: []string{"5d41402abc4b2a76b9719d911017c592"},
			want:   `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5d41402abc4b2a76b9719d911017c592'])`,
		},
		{
			name:   "single SHA1 hash",
			hashes: []string{"356a192b7913b04c54574d18c28d46e6395428ab"},
			want:   `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['356a192b7913b04c54574d18c28d46e6395428ab'])`,
		},
		{
			name:   "single SHA256 hash",
			hashes: []string{"5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4"},
			want:   `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4'])`,
		},
		{
			name:   "multiple hashes of same type",
			hashes: []string{"hash1", "hash2", "hash3"},
			want:   `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['hash1', 'hash2', 'hash3'])`,
		},
		{
			name:   "multiple hashes of mixed types",
			hashes: []string{"5d41402abc4b2a76b9719d911017c592", "356a192b7913b04c54574d18c28d46e6395428ab", "5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4"},
			want:   `where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5d41402abc4b2a76b9719d911017c592', '356a192b7913b04c54574d18c28d46e6395428ab', '5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4'])`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backend.generateHashQuery(tt.hashes)
			if got != tt.want {
				t.Errorf("generateHashQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestR7Backend_generateDomainQuery(t *testing.T) {
	backend := New(nil)
	tests := []struct {
		name    string
		domains []string
		want    string
	}{
		{
			name:    "single domain",
			domains: []string{"malicious.com"},
			want:    `where("query","url" ICONTAINS-ANY ['malicious.com'])`,
		},
		{
			name:    "multiple domains",
			domains: []string{"evil.com", "bad.net"},
			want:    `where("query","url" ICONTAINS-ANY ['evil.com', 'bad.net'])`,
		},
		{
			name:    "domain with hyphen",
			domains: []string{"evil-domain.com"},
			want:    `where("query","url" ICONTAINS-ANY ['evil-domain.com'])`,
		},
		{
			name:    "domain with special characters requiring escaping",
			domains: []string{"domain'with'quotes.com"},
			want:    `where("query","url" ICONTAINS-ANY ['domain\'with\'quotes.com'])`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backend.generateDomainQuery(tt.domains)
			if got != tt.want {
				t.Errorf("generateDomainQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestR7Backend_generateIPQuery(t *testing.T) {
	backend := New(nil)
	tests := []struct {
		name string
		ips  []string
		want string
	}{
		{
			name: "single IPv4 address",
			ips:  []string{"192.168.1.1"},
			want: `where("source_address","destination_address" IN ['192.168.1.1'])`,
		},
		{
			name: "multiple IPv4 addresses",
			ips:  []string{"10.0.0.1", "172.16.0.1"},
			want: `where("source_address","destination_address" IN ['10.0.0.1', '172.16.0.1'])`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backend.generateIPQuery(tt.ips)
			if got != tt.want {
				t.Errorf("generateIPQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestR7Backend_GenerateQuery(t *testing.T) {
	backend := New(nil)
	tests := []struct {
		name    string
		iocs    *extraction.IOCSet
		want    []string
		wantErr bool
	}{
		{
			name: "only MD5 hashes",
			iocs: &extraction.IOCSet{
				MD5Hashes: []string{"5d41402abc4b2a76b9719d911017c592"},
			},
			want:    []string{`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5d41402abc4b2a76b9719d911017c592'])`},
			wantErr: false,
		},
		{
			name: "only SHA1 hashes",
			iocs: &extraction.IOCSet{
				SHA1Hashes: []string{"356a192b7913b04c54574d18c28d46e6395428ab"},
			},
			want:    []string{`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['356a192b7913b04c54574d18c28d46e6395428ab'])`},
			wantErr: false,
		},
		{
			name: "only SHA256 hashes",
			iocs: &extraction.IOCSet{
				SHA256Hashes: []string{"5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4"},
			},
			want:    []string{`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['5394bb17630ed1c849ebc50d6d11a0c5d99037c2073b261f32bd66a618dd4df4'])`},
			wantErr: false,
		},
		{
			name: "mixed hash types",
			iocs: &extraction.IOCSet{
				MD5Hashes:    []string{"md5hash"},
				SHA1Hashes:   []string{"sha1hash"},
				SHA256Hashes: []string{"sha256hash"},
			},
			want:    []string{`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['md5hash', 'sha1hash', 'sha256hash'])`},
			wantErr: false,
		},
		{
			name: "only domains",
			iocs: &extraction.IOCSet{
				Domains: []string{"malicious.com"},
			},
			want:    []string{`where("query","url" ICONTAINS-ANY ['malicious.com'])`},
			wantErr: false,
		},
		{
			name: "only IPs",
			iocs: &extraction.IOCSet{
				IPv4Addresses: []string{"192.168.1.1"},
			},
			want:    []string{`where("source_address","destination_address" IN ['192.168.1.1'])`},
			wantErr: false,
		},
		{
			name: "hashes and domains",
			iocs: &extraction.IOCSet{
				MD5Hashes: []string{"md5hash"},
				Domains:   []string{"evil.com"},
			},
			want: []string{
				`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['md5hash'])`,
				`where("query","url" ICONTAINS-ANY ['evil.com'])`,
			},
			wantErr: false,
		},
		{
			name: "hashes and IPs",
			iocs: &extraction.IOCSet{
				SHA256Hashes:  []string{"sha256hash"},
				IPv4Addresses: []string{"1.2.3.4"},
			},
			want: []string{
				`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['sha256hash'])`,
				`where("source_address","destination_address" IN ['1.2.3.4'])`,
			},
			wantErr: false,
		},
		{
			name: "domains and IPs",
			iocs: &extraction.IOCSet{
				Domains:       []string{"evil.com"},
				IPv4Addresses: []string{"1.2.3.4"},
			},
			want: []string{
				`where("query","url" ICONTAINS-ANY ['evil.com'])`,
				`where("source_address","destination_address" IN ['1.2.3.4'])`,
			},
			wantErr: false,
		},
		{
			name: "all IOC types combined",
			iocs: &extraction.IOCSet{
				MD5Hashes:     []string{"md5hash"},
				SHA1Hashes:    []string{"sha1hash"},
				SHA256Hashes:  []string{"sha256hash"},
				Domains:       []string{"evil.com"},
				IPv4Addresses: []string{"1.2.3.4"},
			},
			want: []string{
				`where("parent_process.exe_file.hashes.*","process.exe_file.hashes.*" IN ['md5hash', 'sha1hash', 'sha256hash'])`,
				`where("query","url" ICONTAINS-ANY ['evil.com'])`,
				`where("source_address","destination_address" IN ['1.2.3.4'])`,
			},
			wantErr: false,
		},
		{
			name:    "empty IOCSet",
			iocs:    &extraction.IOCSet{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "nil IOCSet",
			iocs:    nil,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := backend.GenerateQuery(tt.iocs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("GenerateQuery() returned %d queries, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GenerateQuery()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestR7Backend_GenerateQueries(t *testing.T) {
	backend := New(nil)
	tests := []struct {
		name      string
		iocs      *extraction.IOCSet
		wantCount int
		wantErr   bool
		validate  func(t *testing.T, queries []string)
	}{
		{
			name: "multiple hashes across hash sub-types grouped into one hash query",
			iocs: &extraction.IOCSet{
				MD5Hashes:  []string{"hash1"},
				SHA1Hashes: []string{"hash2"},
			},
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, queries []string) {
				if !strings.Contains(queries[0], "hash1") {
					t.Errorf("Query should contain hash1")
				}
				if !strings.Contains(queries[0], "hash2") {
					t.Errorf("Query should contain hash2")
				}
			},
		},
		{
			name: "multiple domains grouped into one domain query",
			iocs: &extraction.IOCSet{
				Domains: []string{"evil.com", "bad.net"},
			},
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, queries []string) {
				if !strings.Contains(queries[0], "evil.com") {
					t.Errorf("Query should contain evil.com")
				}
				if !strings.Contains(queries[0], "bad.net") {
					t.Errorf("Query should contain bad.net")
				}
			},
		},
		{
			name: "multiple IPs grouped into one IP query",
			iocs: &extraction.IOCSet{
				IPv4Addresses: []string{"1.1.1.1", "2.2.2.2"},
			},
			wantCount: 1,
			wantErr:   false,
			validate: func(t *testing.T, queries []string) {
				if !strings.Contains(queries[0], "1.1.1.1") {
					t.Errorf("Query should contain 1.1.1.1")
				}
				if !strings.Contains(queries[0], "2.2.2.2") {
					t.Errorf("Query should contain 2.2.2.2")
				}
			},
		},
		{
			name: "mixed types",
			iocs: &extraction.IOCSet{
				MD5Hashes:     []string{"hash1"},
				Domains:       []string{"evil.com"},
				IPv4Addresses: []string{"1.2.3.4"},
			},
			wantCount: 3,
			wantErr:   false,
			validate: func(t *testing.T, queries []string) {
				foundHash := false
				foundDomain := false
				foundIP := false
				for _, q := range queries {
					if strings.Contains(q, "hash1") {
						foundHash = true
					}
					if strings.Contains(q, "evil.com") {
						foundDomain = true
					}
					if strings.Contains(q, "1.2.3.4") {
						foundIP = true
					}
				}
				if !foundHash {
					t.Errorf("Should have a query with hash1")
				}
				if !foundDomain {
					t.Errorf("Should have a query with evil.com")
				}
				if !foundIP {
					t.Errorf("Should have a query with 1.2.3.4")
				}
			},
		},
		{
			name:      "empty IOCSet",
			iocs:      &extraction.IOCSet{},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := backend.GenerateQueries(tt.iocs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateQueries() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("GenerateQueries() returned %d queries, want %d", len(got), tt.wantCount)
				return
			}
			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}

func TestR7Backend_escapeString(t *testing.T) {
	backend := New(nil)
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no special characters",
			input: "normal-domain.com",
			want:  "normal-domain.com",
		},
		{
			name:  "single quote",
			input: "domain'with'quotes",
			want:  "domain\\'with\\'quotes",
		},
		{
			name:  "backslash",
			input: `domain\with\backslash`,
			want:  `domain\\with\\backslash`,
		},
		{
			name:  "both backslash and quote",
			input: `domain\'test`,
			want:  `domain\\\'test`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backend.escapeString(tt.input)
			if got != tt.want {
				t.Errorf("escapeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
