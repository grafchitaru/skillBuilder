package main

import (
	"database/sql"
	"fmt"
	"github.com/grafchitaru/skillBuilder/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx
	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

const MigrationsDir = "./migrations"

func migrateCommand() *cli.Command {
	var nativeDB *sql.DB
	cfg := *config.NewConfig()
	return &cli.Command{
		Name:        "migration",
		Description: "Goose migration cli (https://github.com/pressly/goose)",
		Before: func(c *cli.Context) error {
			var err error
			nativeDB, err = sql.Open("pgx", cfg.PostgresDatabaseDsn)
			if err != nil {
				fmt.Println("failed to connect to postgres: %w", err)
				return err
			}
			return nil
		},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			command := ""
			dir := MigrationsDir
			if len(args) > 0 {
				command = args[0]
				args = args[1:]
			}

			if command == "create" {
				if len(args) > 0 {
					args = append(args, "sql")
				}
			}

			return goose.RunWithOptionsContext(ctx.Context, command, nativeDB, dir, args, goose.WithAllowMissing())
		},
	}
}
