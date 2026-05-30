package sslutils

import (
	"crypto/tls"
	"crypto/x509/pkix"
	"fmt"
	"math"
	"net"
	"strings"
	"time"
)

type DomainData struct {
	Expiration string
	ExpiresIn  int
	Sans       []string
	Issuer     string
	SubjectOrg string
	CommonName string
	Version    int
	Domain     string
	HostName   string
	IP         string
}

func Check(domain string) (DomainData, error) {
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return DomainData{}, fmt.Errorf("server doesn't support SSL certificate: %w", err)
	}
	defer conn.Close()

	if err := conn.VerifyHostname(domain); err != nil {
		return DomainData{}, fmt.Errorf("hostname doesn't match certificate: %w", err)
	}

	state := conn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return DomainData{}, fmt.Errorf("no peer certificate")
	}
	cert := state.PeerCertificates[0]

	ip, err := remoteIP(conn.RemoteAddr())
	if err != nil {
		return DomainData{}, fmt.Errorf("remote address: %w", err)
	}

	hostName := reverseHostname(ip)
	expiry := cert.NotAfter
	expiresIn := int(math.Floor(expiry.Sub(time.Now()).Hours() / 24))

	return DomainData{
		Expiration: expiry.Format(time.RFC850),
		ExpiresIn:  expiresIn,
		Sans:       cert.DNSNames,
		Issuer:     distinguishedName(cert.Issuer),
		SubjectOrg: subjectOrganization(cert.Subject),
		CommonName: cert.Subject.CommonName,
		Version:    cert.Version,
		Domain:     domain,
		HostName:   hostName,
		IP:         ip,
	}, nil
}

func remoteIP(addr net.Addr) (string, error) {
	host, _, err := net.SplitHostPort(addr.String())
	if err != nil {
		return "", err
	}
	return host, nil
}

func reverseHostname(ip string) string {
	names, err := net.LookupAddr(ip)
	if err != nil || len(names) == 0 {
		return "N/A"
	}
	return strings.TrimSuffix(names[0], ".")
}

func distinguishedName(name pkix.Name) string {
	if name.CommonName != "" {
		return name.CommonName
	}
	if len(name.Organization) > 0 {
		return name.Organization[0]
	}
	return "N/A"
}

func subjectOrganization(name pkix.Name) string {
	if len(name.Organization) > 0 {
		return name.Organization[0]
	}
	return "N/A"
}
