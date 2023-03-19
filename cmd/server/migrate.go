package main

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var dsn string

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate",
		Short:   "Runs SQL Migrations.",
		Long:    "migrate runs SQL migrations found in the db/migrations directory to the provided DSN. For SQLite, the DSN should look like 'sqlite3://path/to/file.db'",
		Example: "url-shortener-api migrate -d sqlite://db/url.db",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := migrate.New("file://db/migrations", dsn)
			if err != nil {
				return err
			}

			err = m.Up()
			if err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("error running migrations: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&dsn, "dsn", "d", "", "DSN for database.")
	cmd.MarkFlagRequired("dsn")

	return cmd
}
