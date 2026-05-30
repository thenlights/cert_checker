package datautils

const Csv = "CSV"

type SourceFormat string

type IpToHostName map[string]string

type KnowledgeBase map[string]map[string]string

func KnownHosts() IpToHostName {
	return IpToHostName{
		"127.0.0.1":       "localhost",
		"142.250.181.174": "google.com",
		"52.142.124.215":  "duck.com",
		"150.171.28.10":   "bing.com",
	}
}

func InfoFrom(source SourceFormat, filename string, separator string) KnowledgeBase {
	switch source {
	case Csv:
		return ReadCsv(filename, separator)
	default:
		return KnowledgeBase{}
	}
}
