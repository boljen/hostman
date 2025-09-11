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

func ResolveConfigFile(filename string) (*ConfigFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfig(filename, data)
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
