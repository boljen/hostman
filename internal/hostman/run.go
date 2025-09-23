package hostman

import (
	"errors"
	"fmt"
)

func Run(cfg Config) error {
	if cfg.Watchmode {
		return RunAndWatch(cfg)
	} else {
		return RunOnce(cfg)
	}
}

func RunOnce(cfg Config) error {
	cfgFile, err := ResolveConfigFilePath(cfg.Filename)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not resolve config file '%s'", cfg.Filename))
	}

	project, err := ParseProjectFile(cfgFile)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not parse config file '%s', got error: %s", cfgFile, err))
	}

	sources, err := project.GetMapping()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get mapping sources, got error: %s", err))
	}

	hostsFile, err := OpenHostsFile(cfg.Hostsfile)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open hosts file '%s', got error: %s", cfg.Hostsfile, err))
	}

	return hostsFile.Update(project.Project, project.filepath, sources)
}

func RunAndWatch(cfg Config) error {

	return nil
}
