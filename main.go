package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/boljen/hostman/internal/hostman"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {

	if err := (&cli.Command{
		Name: "hostman",

		Commands: []*cli.Command{
			{
				Name:  "apply",
				Usage: "Apply a project file with domain mapping to the operating systems hosts file",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "file",
						Usage:   "Name of the hostman project file. If only a filename is given it will look for the file in the current directory or any of the parent directories.",
						Aliases: []string{"f"},
						Value:   "hostman.hcl",
					},
					&cli.StringFlag{
						Name:    "hostsfile",
						Usage:   "location of the os hosts file",
						Aliases: []string{"h"},
						Value:   HOST_FILE_NATIVE_PATH,
					},
					&cli.BoolFlag{
						Name:    "watch",
						Usage:   "Enable watch mode",
						Aliases: []string{"w"},
					},
				},

				Action: func(ctx context.Context, cmd *cli.Command) error {
					watchMode := cmd.Bool("watch")
					filename := cmd.String("file")
					hostsFile := cmd.String("hostsfile")

					return hostman.Run(hostman.RunArgs{
						Watchmode: watchMode,
						Filename:  filename,
						Hostsfile: hostsFile,
					})
				},
			},
			{
				Name:  "init",
				Usage: "Initializes a new hostman project file (hostman.hcl) in the current directory",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cwd, err := os.Getwd()
					if err != nil {
						return err
					}
					return hostman.InitProjectFile(cwd)
				},
			},
			{
				Name:  "cat",
				Usage: "Outputs the content of the hosts file",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					data, err := os.ReadFile(cmd.String("hostsfile"))
					if err != nil {
						return errors.New(fmt.Sprintf("Could not read hosts file '%s' because error '%s'", cmd.String("hostsfile"), err))
					}
					os.Stdout.Write(data)
					return nil
				},
			},
			{
				Name:  "list",
				Usage: "Lists current projects in the hosts file together with their host mapping and original file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project",
						Usage: "name of the specific project to list hosts configuration",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return hostman.List(hostman.ListArgs{
						Hostsfile: cmd.String("hostsfile"),
						Project:   cmd.String("project"),
					})
				},
			},
			{
				Name:  "clean",
				Usage: "removes hostman managed blocks in the hosts file. If no project is specified, all managed blocks are removed",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project",
						Usage: "name of the specific project to cleanup",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return errors.New("Not yet implemented")
				},
			},
		},
	}).Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
