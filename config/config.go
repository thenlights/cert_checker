package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	CsvSeparator     string `json:"csv_separator"`
	KnowledgeCsvFile string `json:"knowledge_csv_file"`
	IpToHostFile     string `json:"ip_to_host_file"`
}

// Tries loading a configuration from file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		CsvSeparator:     ",",
		KnowledgeCsvFile: "knowledge.csv",
		IpToHostFile:     "ip_to_host.csv",
	}

}
