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
