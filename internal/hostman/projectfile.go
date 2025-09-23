package hostman

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"os"
	"regexp"
	"strconv"
)

type ProjectFile struct {
	Project string           `hcl:"project"`
	Sources []string         `hcl:"sources"`
	Static  []*StaticMapping `hcl:"static,block"`
	Http    []*HTTPMapping   `hcl:"http,block"`

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

func (file *ProjectFile) GetMappingSources() ([]Mapping, error) {
	mapped, err := file.validateAndMapSources()
	if err != nil {
		return nil, err
	}

	list := make([]Mapping, 0, len(file.Sources))

	for _, src := range file.Sources {
		if _, has := mapped[src]; !has {
			return nil, fmt.Errorf("source '%s' not found", src)
		}
		list = append(list, mapped[src])
	}

	return list, nil
}

func (file *ProjectFile) Validate() error {
	if err := file.validateSourceConfigs(); err != nil {
		return err
	}
	sourcesMap, err := file.validateAndMapSources()
	if err != nil {
		return err
	}
	if file.Sources == nil {
		return errors.New("sources cannot be empty")
	}
	for i, src := range file.Sources {
		if _, has := sourcesMap[src]; !has {
			return fmt.Errorf("sources["+strconv.Itoa(i)+"] with name '%s' not found", src)
		}
	}
	return nil
}

func (file *ProjectFile) validateSourceConfigs() error {
	sources := file.getAvailableSources()
	errorList := make([]error, 0)
	for _, source := range sources {
		if err := source.Validate(); err != nil {
			errorList = append(errorList, err)
		}
		if err := validateName(source.GetName()); err != nil {
			errorList = append(errorList, err)
		}
	}
	if len(errorList) > 0 {
		return errors.Join(errorList...)
	}
	return nil
}

func (file *ProjectFile) getAvailableSources() []Mapping {
	sources := make([]Mapping, 0, len(file.Static)+len(file.Http))
	for i := 0; i < len(file.Static); i++ {
		sources = append(sources, file.Static[i])
	}
	for i := 0; i < len(file.Http); i++ {
		sources = append(sources, file.Http[i])
	}
	return sources
}

func (file *ProjectFile) validateAndMapSources() (map[string]Mapping, error) {
	sourcesMap := make(map[string]Mapping, len(file.Static)+len(file.Http))
	sources := make([]Mapping, 0, len(file.Static)+len(file.Http))
	for _, s := range file.Static {
		sources = append(sources, s)
	}
	for _, s := range file.Http {
		sources = append(sources, s)
	}

	errorList := make([]error, 0)

	for _, s := range sources {
		name := s.GetName()
		if _, exists := sourcesMap[name]; exists {
			errorList = append(errorList, fmt.Errorf("duplicate source name: %s", name))
			continue
		}
		sourcesMap[name] = s
	}

	if len(errorList) > 0 {
		return nil, errors.Join(errorList...)
	}

	return sourcesMap, nil
}

var validateNameRegex = regexp.MustCompile(`\S`)

func validateName(name string) error {
	if name == "" {
		return errors.New("source name cannot be empty")
	}
	if !validateNameRegex.MatchString(name) {
		return fmt.Errorf("source name '%s' cannot contain whitespace", name)
	}
	return nil
}
