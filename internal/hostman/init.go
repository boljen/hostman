package hostman

import (
	"os"
	"path/filepath"
)

const ProjectFileTemplate = `project = "demo"

sources = [
	"single"
	"multiple"
	"remote"
]

static "single" {
	ip = "127.0.0.1"
	host = "example.com"
}

static "multiple" {
	ip = "127.0.0.1"
	hosts = [
		"a.example.com", 
		"b.example.com"
	]
}

http "remote" {
	endpoint = "https://raw.githubusercontent.com/boljen/hostman/refs/heads/master/examples/http/response.json"
}
`

func InitProjectFile(dir string) error {
	filename := filepath.Join(dir, "hostman.hcl")
	err := os.WriteFile(filename, []byte(ProjectFileTemplate), 0o644)
	return err
}
