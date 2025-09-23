package hostman

import (
	"path/filepath"
	"strings"
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

func TestGetMappingSources_OrderAndSelection(t *testing.T) {
	pf := &ProjectFile{
		Sources: []string{"two", "one"},
		Static:  []*StaticMapping{},
		Http:    []*HTTPMapping{},
	}

	host := "example"
	pf.Static = []*StaticMapping{{Name: "one", Host: &host, Ip: "127.0.0.1"}}
	pf.Http = []*HTTPMapping{{Name: "two", Endpoint: "http://127.invalid/"}}

	got, err := pf.GetMappingSources()
	if err != nil {
		t.Fatalf("GetMappingSources returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("unexpected number of mappings: got %d, want %d", len(got), 2)
	}
	if got[0].GetName() != "two" || got[1].GetName() != "one" {
		t.Fatalf("unexpected order or names: got [%q, %q]", got[0].GetName(), got[1].GetName())
	}
}

func TestGetMappingSources_SourceNotFound(t *testing.T) {
	host := "h"
	pf := &ProjectFile{
		Sources: []string{"missing"},
		Static:  []*StaticMapping{{Name: "present", Host: &host, Ip: "127.0.0.1"}},
	}

	_, err := pf.GetMappingSources()
	if err == nil {
		t.Fatalf("expected error when source name is missing")
	}
	if !contains(err.Error(), "source 'missing' not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetMappingSources_DuplicateNamesError(t *testing.T) {
	host1 := "a"
	host2 := "b"
	pf := &ProjectFile{
		Sources: []string{"dup"},
		Static:  []*StaticMapping{{Name: "dup", Host: &host1, Ip: "1.1.1.1"}},
		Http:    []*HTTPMapping{{Name: "dup", Endpoint: "http://127.invalid/"}},
	}

	_, err := pf.GetMappingSources()
	if err == nil {
		t.Fatalf("expected error due to duplicate source names, got nil")
	}
	// joined error should mention duplicate source name
	if !contains(err.Error(), "duplicate source name: dup") {
		t.Fatalf("error does not mention duplicate: %v", err)
	}
	_ = host2 // keep variable used in case of future extension
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }
