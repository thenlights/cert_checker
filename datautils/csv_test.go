package datautils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testdataPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("..", "testdata", name)
}

func TestReadKnowledgeCsv(t *testing.T) {
	path := testdataPath(t, "knowledge.csv")
	data, err := ReadKnowledgeCsv(path, ",")
	if err != nil {
		t.Fatalf("ReadKnowledgeCsv() error = %v", err)
	}

	row, ok := data["example.com"]
	if !ok {
		t.Fatal("missing example.com entry")
	}
	if row["server_app"] != "app-server" {
		t.Errorf("server_app = %q, want app-server", row["server_app"])
	}
	if row["balancer"] != "no" {
		t.Errorf("balancer = %q, want no", row["balancer"])
	}

	short, ok := data["short.row"]
	if !ok {
		t.Fatal("missing short.row entry")
	}
	if _, ok := short["balancer"]; ok {
		t.Error("short row should not have balancer column filled from missing cells")
	}
}

func TestReadKnowledgeCsv_semicolonSeparator(t *testing.T) {
	path := testdataPath(t, "knowledge_semicolon.csv")
	data, err := ReadKnowledgeCsv(path, ";")
	if err != nil {
		t.Fatalf("ReadKnowledgeCsv() error = %v", err)
	}

	row, ok := data["semi.example"]
	if !ok {
		t.Fatal("missing semi.example entry")
	}
	if row["server_app"] != "semi-app" {
		t.Errorf("server_app = %q, want semi-app", row["server_app"])
	}
}

func TestReadKnowledgeCsv_missingFile(t *testing.T) {
	_, err := ReadKnowledgeCsv(filepath.Join(t.TempDir(), "missing.csv"), ",")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadHostsCsv(t *testing.T) {
	path := testdataPath(t, "ip_to_host.csv")
	data, err := ReadHostsCsv(path, ",")
	if err != nil {
		t.Fatalf("ReadHostsCsv() error = %v", err)
	}

	if data["10.0.0.1"] != "app-server" {
		t.Errorf("10.0.0.1 = %q, want app-server", data["10.0.0.1"])
	}
	if data["192.168.1.5"] != "db.internal" {
		t.Errorf("192.168.1.5 = %q, want db.internal", data["192.168.1.5"])
	}
}

func TestReadHostsCsv_invalidRow(t *testing.T) {
	path := testdataPath(t, "invalid_hosts.csv")
	_, err := ReadHostsCsv(path, ",")
	if err == nil {
		t.Fatal("expected error for invalid row")
	}
	if !strings.Contains(err.Error(), "row 2") {
		t.Errorf("error = %v, want mention of row 2", err)
	}
}

func TestReadKnowledgeCsv_emptySeparatorUsesComma(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "k.csv")
	content := "domain,value\nfoo,bar\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := ReadKnowledgeCsv(path, "")
	if err != nil {
		t.Fatal(err)
	}
	if data["foo"]["value"] != "bar" {
		t.Errorf("got %#v", data["foo"])
	}
}
