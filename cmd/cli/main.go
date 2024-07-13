package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

func main() {
	cliApp := &cli.App{
		Commands: []*cli.Command{
			migrateCommand(),
		},
		Before: func(c *cli.Context) error {
			goose.SetBaseFS(os.DirFS(MigrationsDir))

			c.Context, _ = signal.NotifyContext(c.Context, os.Interrupt)

			return nil
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		fmt.Println("failed to run cli: %w", err)
	}
}
