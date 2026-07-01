package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.CsvSeparator != "," {
		t.Errorf("CsvSeparator = %q, want %q", cfg.CsvSeparator, ",")
	}
	if cfg.KnowledgeCsvFile != "knowledge.csv" {
		t.Errorf("KnowledgeCsvFile = %q, want knowledge.csv", cfg.KnowledgeCsvFile)
	}
	if cfg.IpToHostFile != "ip_to_host.csv" {
		t.Errorf("IpToHostFile = %q, want ip_to_host.csv", cfg.IpToHostFile)
	}
}

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "csv_separator": ";",
  "knowledge_csv_file": "custom_knowledge.csv",
  "ip_to_host_file": "custom_hosts.csv"
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.CsvSeparator != ";" {
		t.Errorf("CsvSeparator = %q, want ;", cfg.CsvSeparator)
	}
	if cfg.KnowledgeCsvFile != "custom_knowledge.csv" {
		t.Errorf("KnowledgeCsvFile = %q", cfg.KnowledgeCsvFile)
	}
	if cfg.IpToHostFile != "custom_hosts.csv" {
		t.Errorf("IpToHostFile = %q", cfg.IpToHostFile)
	}
}

func TestLoadConfig_missingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_invalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
