package datautils

import (
	"testing"
)

func TestInfoFrom_csv(t *testing.T) {
	data, err := InfoFrom(Csv, testdataPath(t, "knowledge.csv"), ",")
	if err != nil {
		t.Fatalf("InfoFrom() error = %v", err)
	}
	if len(data) < 1 {
		t.Fatal("expected at least one knowledge row")
	}
}

func TestInfoFrom_unknownSource(t *testing.T) {
	data, err := InfoFrom(SourceFormat("JSON"), "ignored", ",")
	if err != nil {
		t.Fatalf("InfoFrom() error = %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(data))
	}
}

func TestKnownHostsFromSource_csv(t *testing.T) {
	data, err := KnownHostsFromSource(Csv, testdataPath(t, "ip_to_host.csv"), ",")
	if err != nil {
		t.Fatalf("KnownHostsFromSource() error = %v", err)
	}
	if data["10.0.0.1"] != "app-server" {
		t.Errorf("lookup = %q, want app-server", data["10.0.0.1"])
	}
}

func TestKnownHostsFromSource_unknownSource(t *testing.T) {
	data, err := KnownHostsFromSource(SourceFormat("YAML"), "ignored", ",")
	if err != nil {
		t.Fatalf("KnownHostsFromSource() error = %v", err)
	}
	if len(data) != 0 {
		t.Fatal("expected empty map for unknown source")
	}
}
