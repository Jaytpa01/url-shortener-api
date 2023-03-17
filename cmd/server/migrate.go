package main

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var dsn string

func migrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Runs SQL Migrations.",
		Long:  "migrate runs SQL migrations found in the db/migrations directory to the provided DSN.",
		RunE: func(cmd *cobra.Command, args []string) error {

			db, err := sql.Open("sqlite3", dsn)
			if err != nil {
				return fmt.Errorf("couldn't open db connection: %w", err)
			}

			driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
			if err != nil {
				return fmt.Errorf("couldn't create sqlite3 driver: %w", err)
			}

			m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
			if err != nil {
				return fmt.Errorf("error creating migration instance: %w", err)
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
