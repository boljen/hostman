package hostman

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenHostsFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "hosts")
	content := "127.0.0.1 localhost\n::1 localhost"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp hosts file: %v", err)
	}

	h, err := OpenHostsFile(path)
	if err != nil {
		t.Fatalf("OpenHostsFile unexpected error: %v", err)
	}
	if h == nil {
		t.Fatalf("OpenHostsFile returned nil without error")
	}
	if h.location != path {
		t.Fatalf("unexpected location: got %q want %q", h.location, path)
	}
	if h.content != content {
		t.Fatalf("unexpected content: got %q want %q", h.content, content)
	}
}

func TestOpenHostsFile_Nonexistent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no-such-file")
	f, err := OpenHostsFile(path)
	if err == nil {
		t.Fatalf("expected error for non-existent file, got none and file=%#v", f)
	}
	if f != nil {
		t.Fatalf("expected returned file to be nil on error, got: %#v", f)
	}
}

func TestGetProjects_one(t *testing.T) {
	hostsFile := HostsFile{
		content: `### hostman-project-start demo ###`,
	}

	projects := hostsFile.GetProjects()

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}
	if projects[0] != "demo" {
		t.Fatalf("expected project name to be 'demo', got %q", projects[0])
	}

}

func TestGetProjects_multiple(t *testing.T) {
	hostsFile := HostsFile{
		content: `### hostman-project-start demo ###

### hostman-project-start demo2 ###
`,
	}

	projects := hostsFile.GetProjects()

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
	if projects[0] != "demo" {
		t.Fatalf("expected project name to be 'demo', got %q", projects[0])
	}
	if projects[1] != "demo2" {
		t.Fatalf("expected project name to be 'demo2', got %q", projects[0])
	}

}

func TestGetProjects_none(t *testing.T) {
	hostsFile := HostsFile{
		content: `### hostman-project-start demo #a##`,
	}

	projects := hostsFile.GetProjects()

	if len(projects) != 0 {
		t.Fatalf("expected 0 project, got %d", len(projects))
	}

}

func TestRemoveProject_one(t *testing.T) {
	hostsFile := HostsFile{
		content: `
127.0.0.1 localhost
### hostman-project-start demo ###


### hostman-project-end demo ###
127.0.0.1 localhost2
`,
	}

	if err := hostsFile.RemoveProject("demo"); err != nil {
		t.Fatalf("RemoveProject unexpected error: %v", err)
	}

	if hostsFile.content != `
127.0.0.1 localhost
127.0.0.1 localhost2
` {
		t.Fatalf("expected content to be empty after removal")
	}

}

func TestRemoveProject_all(t *testing.T) {
	hostsFile := HostsFile{
		content: `
127.0.0.1 localhost
### hostman-project-start demo ###


### hostman-project-end demo ###
### hostman-project-start demo2 ###


### hostman-project-end demo2 ###
127.0.0.1 localhost2
`,
	}

	if err := hostsFile.RemoveAllProjects(); err != nil {
		t.Fatalf("RemoveProject unexpected error: %v", err)
	}

	if hostsFile.content != `
127.0.0.1 localhost
127.0.0.1 localhost2
` {
		t.Fatalf("expected content to be empty after removal")
	}

}
