package hostman

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

type RunArgs struct {
	Watchmode bool
	Filename  string
	Hostsfile string
}

func Run(cfg RunArgs) error {
	if cfg.Watchmode {
		return RunAndWatch(cfg)
	} else {
		return RunOnce(cfg)
	}
}

func RunOnce(cfg RunArgs) error {
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

func RunAndWatch(cfg RunArgs) error {

	fmt.Println("Starting")
	if err := RunOnce(cfg); err != nil {
		return err
	}

	ticker := time.Tick(5 * time.Second)

	for {
		select {
		case <-ticker:
			if err := RunOnce(cfg); err != nil {
				log.Printf("Error while refreshing: %s", err)
			}
		}
	}

}

type ListArgs struct {
	Hostsfile string
	Project   string
}

func List(args ListArgs) error {
	data, err := os.ReadFile(args.Hostsfile)
	if err != nil {
		return err
	}

	log.Printf("Hosts file content:\n%s", string(data))

	return nil
}
