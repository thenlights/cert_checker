package sslutils

import (
	"crypto/tls"
	"math"
	"net"
	"time"
)

type DomainData struct {
	Expiration string
	ExpiresIn  int
	Sans       []string
	Issuer     string
	Org        string
	CommonName string
	Version    int
	Domain     string
	HostName   string
	Ip         string
}

func Check(domain string) DomainData {

	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}

	err = conn.VerifyHostname(domain)
	if err != nil {
		panic("Hostname doesn't match with certificate: " + err.Error())
	}

	resolved, err := net.LookupIP(domain)
	if err != nil {
		panic("Unable to resolve host IP: " + err.Error())
	}

	name, err := net.LookupAddr(resolved[0].String())
	if err != nil {
		panic("Unable to resolve host Name: " + err.Error())
	}
	cert := conn.ConnectionState().PeerCertificates[0]
	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter

	expiresIn := int(math.Floor(expiry.Sub(time.Now()).Hours() / 24))

	var issuer, ip string
	switch {
	case len(cert.Subject.Organization) > 0:
		issuer = cert.Subject.Organization[0]
	case len(cert.Issuer.Organization) > 0:
		issuer = cert.Issuer.Organization[0]
	default:
		issuer = "N/A"
	}

	if len(resolved[0]) > 0 {
		ip = resolved[0].String()
	} else {
		ip = "N/A"
	}

	data := DomainData{
		Expiration: expiry.Format(time.RFC850),
		ExpiresIn:  expiresIn,
		Sans:       cert.DNSNames,
		Issuer:     issuer,
		Org:        cert.Issuer.Organization[0],
		CommonName: cert.Subject.CommonName,
		Version:    cert.Version,
		Domain:     domain,
		HostName:   name[0],
		Ip:         ip,
	}
	
	return data
}
