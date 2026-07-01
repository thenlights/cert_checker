package sslutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"
)

func TestDistinguishedName(t *testing.T) {
	tests := []struct {
		name string
		in   pkix.Name
		want string
	}{
		{
			name: "common name preferred",
			in:   pkix.Name{CommonName: "Test CA", Organization: []string{"CA Org"}},
			want: "Test CA",
		},
		{
			name: "organization fallback",
			in:   pkix.Name{Organization: []string{"CA Org"}},
			want: "CA Org",
		},
		{
			name: "empty",
			in:   pkix.Name{},
			want: "N/A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := distinguishedName(tt.in); got != tt.want {
				t.Errorf("distinguishedName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSubjectOrganization(t *testing.T) {
	tests := []struct {
		name string
		in   pkix.Name
		want string
	}{
		{
			name: "organization present",
			in:   pkix.Name{CommonName: "www.example.com", Organization: []string{"Example Inc"}},
			want: "Example Inc",
		},
		{
			name: "no organization",
			in:   pkix.Name{CommonName: "www.example.com"},
			want: "N/A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := subjectOrganization(tt.in); got != tt.want {
				t.Errorf("subjectOrganization() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoteIP(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "192.0.2.1:443")
	if err != nil {
		t.Fatal(err)
	}

	ip, err := remoteIP(addr)
	if err != nil {
		t.Fatalf("remoteIP() error = %v", err)
	}
	if ip != "192.0.2.1" {
		t.Errorf("remoteIP() = %q, want 192.0.2.1", ip)
	}
}

func TestRemoteIP_invalidAddress(t *testing.T) {
	_, err := remoteIP(&net.IPAddr{IP: net.ParseIP("192.0.2.1")})
	if err == nil {
		t.Fatal("expected error for address without port")
	}
}

func TestCheckAddr(t *testing.T) {
	const serverName = "test.example.com"
	notAfter := time.Now().Add(45 * 24 * time.Hour)

	addr, cleanup := startTestTLSServer(t, x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   serverName,
			Organization: []string{"Example Inc"},
		},
		Issuer: pkix.Name{
			CommonName:   "Test CA",
			Organization: []string{"Test CA Org"},
		},
		DNSNames:  []string{serverName, "alt.example.com"},
		NotBefore: time.Now().Add(-time.Hour),
		NotAfter:  notAfter,
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
	})
	defer cleanup()

	data, err := checkTLSAddr(addr, serverName, true)
	if err != nil {
		t.Fatalf("CheckAddr() error = %v", err)
	}

	if data.Domain != serverName {
		t.Errorf("Domain = %q, want %q", data.Domain, serverName)
	}
	if data.CommonName != serverName {
		t.Errorf("CommonName = %q, want %q", data.CommonName, serverName)
	}
	if data.SubjectOrg != "Example Inc" {
		t.Errorf("SubjectOrg = %q, want Example Inc", data.SubjectOrg)
	}
	if data.Issuer != "Test CA" {
		t.Errorf("Issuer = %q, want Test CA", data.Issuer)
	}
	if len(data.Sans) != 2 {
		t.Fatalf("Sans = %v, want 2 entries", data.Sans)
	}
	if data.IP == "" || data.IP == "N/A" {
		t.Errorf("IP = %q, want non-empty connected address", data.IP)
	}
	if data.ExpiresIn < 40 || data.ExpiresIn > 45 {
		t.Errorf("ExpiresIn = %d, want roughly 45", data.ExpiresIn)
	}
}

func TestCheckAddr_hostnameMismatch(t *testing.T) {
	addr, cleanup := startTestTLSServer(t, x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{CommonName: "other.example.com"},
		Issuer:                pkix.Name{CommonName: "Test CA"},
		DNSNames:              []string{"other.example.com"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	})
	defer cleanup()

	_, err := checkTLSAddr(addr, "test.example.com", true)
	if err == nil {
		t.Fatal("expected hostname verification error")
	}
	if !strings.Contains(err.Error(), "hostname doesn't match") {
		t.Errorf("error = %v", err)
	}
}

func startTestTLSServer(t *testing.T, leafTemplate x509.Certificate) (addr string, cleanup func()) {
	t.Helper()

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate CA key: %v", err)
	}

	caTemplate := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Test CA", Organization: []string{"Test CA Org"}},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create CA certificate: %v", err)
	}

	leafKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate leaf key: %v", err)
	}

	caCert, err := x509.ParseCertificate(caDER)
	if err != nil {
		t.Fatalf("parse CA certificate: %v", err)
	}

	leafDER, err := x509.CreateCertificate(rand.Reader, &leafTemplate, caCert, &leafKey.PublicKey, caKey)
	if err != nil {
		t.Fatalf("create leaf certificate: %v", err)
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{leafDER, caDER},
		PrivateKey:  leafKey,
	}

	listener, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS12,
	})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				serverConn := tls.Server(c, &tls.Config{
					Certificates: []tls.Certificate{tlsCert},
					MinVersion:   tls.VersionTLS12,
				})
				_ = serverConn.Handshake()
			}(conn)
		}
	}()

	return listener.Addr().String(), func() {
		_ = listener.Close()
		<-done
	}
}
