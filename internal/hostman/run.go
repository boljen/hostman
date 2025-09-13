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

	fmt.Println(project)

	//fmt.Println(cmd.String("hostsfile"))
	//fmt.Println(cmd.Bool("watch"))

	// If watch => watch filename for changes
	// + watch http endpoints for changes
	//

	return nil
}

func RunAndWatch(cfg Config) error {

	return nil
}
