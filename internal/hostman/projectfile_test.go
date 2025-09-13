package hostman

import (
	"path/filepath"
	"testing"
)

func TestParseProjectFile_ExamplesSimple(t *testing.T) {
	cfgPath := filepath.Join("..", "..", "examples", "simple", "hostman.hcl")

	pf, err := ParseProjectFile(cfgPath)
	if err != nil {
		t.Fatalf("ParseProjectFile failed: %v", err)
	}
	if pf == nil {
		t.Fatal("ParseProjectFile returned nil ProjectFile without error")
	}

	// Basic assertions based on examples/simple/hostman.hcl
	if pf.Project != "test" {
		t.Fatalf("unexpected project name: got %q, want %q", pf.Project, "test")
	}

	if len(pf.Sources) != 1 {
		t.Fatalf("unexpected number of sources: got %d, want %d", len(pf.Sources), 1)
	}
	if pf.Sources[0] != "domaintest" {
		t.Fatalf("unexpected first source name: got %q, want %q", pf.Sources[0], "domaintest")
	}

	if len(pf.Static) != 1 {
		t.Fatalf("unexpected number of static blocks: got %d, want %d", len(pf.Static), 1)
	}
	st := pf.Static[0]
	if st.Name != "domaintest" {
		t.Fatalf("unexpected static name: got %q, want %q", st.Name, "domaintest")
	}
	if st.Hosts == nil || len(*st.Hosts) != 1 || (*st.Hosts)[0] != "domaintest" {
		t.Fatalf("unexpected static hosts: got %#v, want [\"domaintest\"]", st.Hosts)
	}
	if st.Ip != "127.0.0.1" {
		t.Fatalf("unexpected static ip: got %q, want %q", st.Ip, "127.0.0.1")
	}
}

func TestParseProjectFile_NonexistentFile(t *testing.T) {
	cfgPath := filepath.Join("..", "..", "examples", "simple", "hostman.hcl-doesnotexist")

	if _, err := ParseProjectFile(cfgPath); err == nil {
		t.Fatal("ParseProjectFile didnt return an error for a non-existent file")
	}
}

func TestParseProjectFileData_InvalidHCL(t *testing.T) {
	invalid := []byte(`project = "demo"
	static "name" {
	  ip = "127.0.0.1"
	  hosts = ["a", "b"
	// no closing brackets/braces`)

	if _, err := ParseProjectFileData("inline-invalid.hcl", invalid); err == nil {
		t.Fatalf("ParseProjectFileData did not return an error for invalid HCL input")
	}
}
