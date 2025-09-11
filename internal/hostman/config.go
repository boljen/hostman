package hostman

import (
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
)

type ConfigFile struct {
	Project string         `hcl:"project"`
	Sources []string       `hcl:"sources"`
	Static  []StaticSource `hcl:"static,block"`
	Http    []HTTPSource   `hcl:"http,block"`
}

// ResolveConfigFile resolves a configuration file path based on user input.
// If input is just a filename (no path separators), it searches the current
// directory and then ascends parent directories until the root to find the file.
// If input contains a path, it resolves that specific path only (relative to cwd if relative).
// Returns the absolute path to the existing file or an error if not found.
func ResolveConfigFile(filename string) (*ConfigFile, error) {
	resolved, err := resolveConfigFilePath(filename)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return nil, err
	}

	return ParseConfig(resolved, data)
}

func ParseConfig(path string, data []byte) (*ConfigFile, error) {
	file := ConfigFile{}
	if err := hclsimple.Decode(path, data, nil, &file); err != nil {
		return nil, err
	}

	//

	// Sources => make sure they all exist

	return &file, nil
}
