package datautils

const Csv = "CSV"

type SourceFormat string

type IpToHostName map[string]string

type KnowledgeBase map[string]map[string]string

func KnownHostsFromSource(source SourceFormat, filename string, separator string) IpToHostName {
	switch source {
	case Csv:
		return ReadHostsCsv(filename, separator)
	default:
		return IpToHostName{}
	}
}

func InfoFrom(source SourceFormat, filename string, separator string) KnowledgeBase {
	switch source {
	case Csv:
		return ReadKnowledgeCsv(filename, separator)
	default:
		return KnowledgeBase{}
	}
}
