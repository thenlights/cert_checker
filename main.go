package main

import (
	"log"

	"sslchecker/datautils"
	"sslchecker/sslutils"
	"sslchecker/termutils"
)

func main() {

	hosts := datautils.InfoFrom(datautils.Csv, "knowledge.csv", ",")
	known := datautils.KnownHostsFromSource(datautils.Csv, "ip_to_host.csv", ",")

	domain := termutils.Ask("Insert domain: ")

	cert, err := sslutils.Check(domain)
	if err != nil {
		log.Fatal(err)
	}

	if value, ok := known[cert.IP]; ok {
		termutils.PrintDebug("Recognized Host: " + value)
	}

	// Print
	if dictionary, ok := hosts[domain]; ok {
		termutils.PrintLine()
		termutils.PrintDebug("Retrieved Data form Knowledge base:")
		termutils.PrintDomain(dictionary["domain"])
		termutils.PrintServerApp(dictionary["server_app"])
		termutils.PrintDestination(dictionary["to"])
		termutils.PrintLoadBalancer(dictionary["balancer"])
		termutils.PrintRedirectOnly(dictionary["redirect"])
		termutils.PrintLine()
	}

	termutils.PrintIssuer(cert.Issuer)
	termutils.PrintCommonName(cert.CommonName)
	termutils.PrintIssuedTo(cert.SubjectOrg)
	termutils.PrintExpiration(cert.Expiration, cert.ExpiresIn)
	termutils.PrintSans(cert.Sans)
	termutils.PrintVersion(cert.Version)
	termutils.PrintDoubleLine()

}
