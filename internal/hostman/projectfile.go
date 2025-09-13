package hostman

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
	"regexp"
)

type ProjectFile struct {
	Project string          `hcl:"project"`
	Sources []string        `hcl:"sources"`
	Static  []*StaticSource `hcl:"static,block"`
	Http    []*HTTPSource   `hcl:"http,block"`

	filepath string
}

func ParseProjectFile(filename string) (*ProjectFile, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseProjectFileData(filename, data)
}

func ParseProjectFileData(filepath string, data []byte) (*ProjectFile, error) {
	file := ProjectFile{
		filepath: filepath,
	}
	if err := hclsimple.Decode(filepath, data, nil, &file); err != nil {
		return nil, err
	}
	if err := file.Validate(); err != nil {
		return nil, err
	}
	return &file, nil
}

func (file *ProjectFile) Validate() error {

	/*

		mappedSources, err := validateAndMapSources(file)
		if err != nil {
			return nil, err
		}
		file.mappedSources = mappedSources

		// Build sortedSources from the Sources list
		file.sortedSources = make([]Source, 0, len(file.Sources))
		seen := make(map[string]struct{})
		for _, name := range file.Sources {
			if name == "" {
				return nil, errors.New("sources list contains empty name")
			}
			if _, dup := seen[name]; dup {
				return nil, fmt.Errorf("duplicate source in sources list: %s", name)
			}
			seen[name] = struct{}{}
			src, ok := file.mappedSources[name]
			if !ok {
				return nil, fmt.Errorf("source referenced but not defined: %s", name)
			}
			file.sortedSources = append(file.sortedSources, src)
		}*/

	return nil
}

func (file *ProjectFile) getSourcesList() []Source {
	sources := make([]Source, 0, len(file.Static)+len(file.Http))
	for i := 0; i < len(file.Static); i++ {
		sources = append(sources, file.Static[i])
	}
	for i := 0; i < len(file.Http); i++ {
		sources = append(sources, file.Http[i])
	}
	return sources
}

func validateAndMapSources(file ProjectFile) (map[string]Source, error) {
	sourcesMap := make(map[string]Source, len(file.Static)+len(file.Http))
	sources := make([]Source, 0, len(file.Static)+len(file.Http))
	for _, s := range file.Static {
		sources = append(sources, s)
	}
	for _, s := range file.Http {
		sources = append(sources, s)
	}

	for _, s := range sources {
		name := s.GetName()
		if err := validateName(name); err != nil {
			return nil, err
		}
		if _, exists := sourcesMap[name]; exists {
			return nil, fmt.Errorf("duplicate source name: %s", name)
		}
		if err := s.Validate(); err != nil {
			return nil, fmt.Errorf("invalid source '%s': %w", name, err)
		}
		sourcesMap[name] = s
	}
	return sourcesMap, nil
}

var validateNameRegex = regexp.MustCompile(`\s`)

func validateName(name string) error {
	if name == "" {
		return errors.New("source name cannot be empty")
	}
	if !validateNameRegex.MatchString(name) {
		return fmt.Errorf("source name '%s' cannot contain whitespace", name)
	}
	return nil
}
