package datautils

const Csv = "CSV"

type SourceFormat string

type IpToHostName map[string]string

type KnowledgeBase map[string]map[string]string

func KnownHostsFromSource(source SourceFormat, filename string, separator string) (IpToHostName, error) {
	switch source {
	case Csv:
		return ReadHostsCsv(filename, separator)
	default:
		return IpToHostName{}, nil
	}
}

func InfoFrom(source SourceFormat, filename string, separator string) (KnowledgeBase, error) {
	switch source {
	case Csv:
		return ReadKnowledgeCsv(filename, separator)
	default:
		return KnowledgeBase{}, nil
	}
}
