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
		Name:  "hostman",
		Usage: "Apply a project file with domain mapping to the operating systems hosts file",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Name of the hostman project config file. If only a filename is given it will look for the file in the current directory or any of the parent directories.",
				DefaultText: "hostman.hcl",
			},
			&cli.StringFlag{
				Name:        "hostsfile",
				Usage:       "location of the os hosts file",
				DefaultText: HOST_FILE_NATIVE_PATH,
			},
			&cli.BoolFlag{
				Name:  "watch",
				Usage: "Enable watch mode",
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			cfgFile, err := hostman.ResolveConfigFilePath(cmd.String("config"))
			if err != nil {
				return errors.New(fmt.Sprintf("Could not resolve config file '%s'", cmd.String("config")))
			}

			fmt.Println(cfgFile)

			//fmt.Println(cmd.String("hostsfile"))
			//fmt.Println(cmd.Bool("watch"))

			// If watch => watch filename for changes
			// + watch http endpoints for changes
			//

			return nil
		},

		Commands: []*cli.Command{
			{
				Name:  "projects",
				Usage: "Lists current projects in the hosts file together with their host mapping and original file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "project",
						Usage: "name of the specific project to list hosts configuration",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return errors.New("Not yet implemented")
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
