package hostman

import (
	"errors"
	"fmt"
	"log"
	"strings"
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

	if err := hostsFile.Update(project.Project, project.filepath, sources); err != nil {
		return errors.New(fmt.Sprintf("Could not update hosts file content, got error: %s", err))
	}

	if err := hostsFile.Save(); err != nil {
		return errors.New(fmt.Sprintf("Could not save hosts file content, got error: %s", err))
	}

	return nil

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
}

func List(args ListArgs) error {
	hostsFile, err := OpenHostsFile(args.Hostsfile)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open hosts file '%s', got error: %s", args.Hostsfile, err))
	}

	projects := hostsFile.GetProjects()

	for _, p := range projects {
		fmt.Printf("Project: %s\n", p)

		hostsData, err := hostsFile.GetProjectData(p)
		if err != nil {
			fmt.Printf("ERROR: Could not get hosts for project '%s', got error: %s\n", p, err)
		}

		hostsData = strings.TrimSpace(hostsData)

		lines := strings.Lines(hostsData)

		for s := range lines {
			fmt.Print("\t" + s + "\n")
		}
	}

	return nil
}

type CleanArgs struct {
	Project   string
	Hostsfile string
}

func Clean(args CleanArgs) error {
	hostsFile, err := OpenHostsFile(args.Hostsfile)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open hosts file '%s', got error: %s", args.Hostsfile, err))
	}

	if args.Project == "" {
		if err := hostsFile.RemoveAllProjects(); err != nil {
			return err
		}
	} else {
		if err := hostsFile.RemoveProject(args.Project); err != nil {
			return err
		}
	}

	return hostsFile.Save()
}
