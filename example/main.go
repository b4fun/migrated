package main

import (
	"database/sql"
	"log"

	"github.com/b4fun/migrated"
	migratedcli "github.com/b4fun/migrated/cobracli"
	"github.com/b4fun/migrated/example/migration"
	"github.com/b4fun/migrated/gobindata"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

type logger struct{}

func (l *logger) Infof(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func (l *logger) Warnf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func (l *logger) Errorf(format string, a ...interface{}) {
	log.Printf(format, a...)
}

func main() {
	cliName := "testctl"
	cli := &cobra.Command{
		Use: cliName,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	migrator := &migrated.MigratorConfig{
		Logger:           &logger{},
		MigrationBaseDir: "./migration",
		CreateDatabase: func() (string, database.Driver, error) {
			db, err := sql.Open("sqlite3", "./test.db")
			if err != nil {
				return "", nil, err
			}
			driver, err := sqlite3.WithInstance(db, &sqlite3.Config{
				MigrationsTable: "migrations",
			})
			if err != nil {
				return "", nil, err
			}
			return "sqlite3", driver, nil
		},
		CreateSource: gobindata.NewSourceFrom(
			migration.AssetNames,
			migration.Asset,
		),
	}
	cli.AddCommand(migratedcli.NewMigrationCommand(cliName, migrator))

	cli.Execute()
}
