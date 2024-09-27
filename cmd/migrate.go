package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"gin.go.dev/pkg/storage/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/spf13/cobra"
	"strconv"
)

var cmdMigrate = &cobra.Command{
	Use:   "migrate",
	Short: "Migrates the database",
}

var cmdMigrateUp = &cobra.Command{
	Use:   "up",
	Short: "Apply all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrateDatabase(func(m *migrate.Migrate) error {
			return m.Up()
		}); err != nil {
			fmt.Println(err)
		}
	},
}

var cmdMigrateDown = &cobra.Command{
	Use:   "down",
	Short: "Apply all down migrations",
	Run: func(cmd *cobra.Command, args []string) {
		if err := migrateDatabase(func(m *migrate.Migrate) error {
			return m.Down()
		}); err != nil {
			fmt.Println(err)
		}
	},
}

var cmdMigrateStep = &cobra.Command{
	Use:   "step [n]",
	Short: "Apply n up or down migrations",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return fmt.Errorf("invalid step count: %v", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		n, _ := strconv.Atoi(args[0])
		if err := migrateDatabase(func(m *migrate.Migrate) error {
			return m.Steps(n)
		}); err != nil {
			fmt.Println(err)
		}
	},
}

func migrateDatabase(migrateFunc func(*migrate.Migrate) error) error {
	if cfg == nil {
		return errors.New("config not initialized")
	}

	conn, err := sql.Open("postgres", cfg.Database.URL().String())
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	fs, err := iofs.New(db.Migrations, "migrations")
	migrator, err := migrate.NewWithInstance("iofs", fs, "postgres", driver)
	if err != nil {
		return err
	}

	if err = migrateFunc(migrator); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
		} else {
			return err
		}
	} else {
		fmt.Println("Migrations applied successfully")
	}
	return nil
}

func init() {
	cmdMigrate.AddCommand(cmdMigrateUp)
	cmdMigrate.AddCommand(cmdMigrateDown)
	cmdMigrate.AddCommand(cmdMigrateStep)
}
