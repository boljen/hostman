package hostman

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveFile_FilenameSearchesParents(t *testing.T) {
	// Arrange
	base := t.TempDir()
	child := filepath.Join(base, "a", "b")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	cfgName := "hostman.hcl"
	cfgPath := filepath.Join(base, cfgName)
	if err := os.WriteFile(cfgPath, []byte("project=\"x\""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	// Act
	resolved, err := resolveFileFrom(cfgName, child)
	if err != nil {
		t.Fatalf("expected to resolve, got err: %v", err)
	}

	// Assert
	if resolved != cfgPath {
		abs, _ := filepath.Abs(cfgPath)
		if resolved != abs {
			t.Fatalf("resolved path mismatch: got %q want %q", resolved, cfgPath)
		}
	}
}

func TestResolveFile_PathOnlyNoSearch(t *testing.T) {
	// Arrange
	base := t.TempDir()
	parent := filepath.Join(base, "p")
	child := filepath.Join(parent, "c")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	name := "hostman.hcl"
	parentFile := filepath.Join(parent, name)
	if err := os.WriteFile(parentFile, []byte("project=\"x\""), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	nonExisting := filepath.Join(child, name)

	// Act
	_, err := resolveFileFrom(nonExisting, child)

	// Assert
	if err == nil {
		t.Fatalf("expected error for non-existing explicit path")
	}
}

func TestIsFilename(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"hostman.hcl", true},
		{"./hostman.hcl", false},
		{"dir/hostman.hcl", false},
		{"dir\\hostman.hcl", false},
		{"C:\\tmp\\hostman.hcl", false},
	}

	for _, c := range cases {
		if got := isFilename(c.in); got != c.want {
			t.Errorf("isFilename(%q)=%v want %v", c.in, got, c.want)
		}
	}
}
