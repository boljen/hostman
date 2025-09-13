package hostman

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitProjectFile_CreatesFileWithTemplate(t *testing.T) {
	// Arrange
	tmp := t.TempDir()

	// Act
	if err := InitProjectFile(tmp); err != nil {
		t.Fatalf("InitProjectFile error: %v", err)
	}

	// Assert
	path := filepath.Join(tmp, "hostman.hcl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}

	if string(data) != ProjectFileTemplate {
		t.Fatalf("template mismatch in created file\n--- got ---\n%q\n--- want ---\n%q", string(data), ProjectFileTemplate)
	}
}

func TestInitProjectFile_OverwritesExisting(t *testing.T) {
	// Arrange
	tmp := t.TempDir()
	path := filepath.Join(tmp, "hostman.hcl")

	// Write some other content first
	if err := os.WriteFile(path, []byte("something else"), 0o644); err != nil {
		t.Fatalf("pre-write failed: %v", err)
	}

	// Act
	if err := InitProjectFile(tmp); err != nil {
		t.Fatalf("InitProjectFile error: %v", err)
	}

	// Assert
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file after init: %v", err)
	}

	if string(data) != ProjectFileTemplate {
		t.Fatalf("expected file to be overwritten with template")
	}
}
