// Commandline entrypoint for migration task.
package cobracli

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/b4fun/migrated"
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// NewMigrationCommand creates a cobra command.
func NewMigrationCommand(
	mainCmd string,
	config *migrated.MigratorConfig,
) *cobra.Command {
	logger := config.Logger

	abortWithError := func(err error) {
		logger.Errorf("migrate: %+v", err)
		os.Exit(1)
	}

	c := &cobra.Command{
		Use:   "migration (create|up|down)",
		Short: fmt.Sprintf("%s migration", mainCmd),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	c.AddCommand(&cobra.Command{
		Use:   "create NAME",
		Short: "create a migration",
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			migrationBaseDir := config.MustGetMigrationBaseDir()
			name := args[0]

			base, err := migrated.CreateMigration(migrationBaseDir, time.Now().Unix(), name)
			if err != nil {
				abortWithError(errors.Wrap(err, "migration.CreateMigration"))
			}
			logger.Infof("created %s\n", base)
		},
	})

	upCommand := &cobra.Command{
		Use:   "up",
		Short: "apply migrations",
		Run: func(c *cobra.Command, args []string) {
			limit, err := c.LocalFlags().GetInt("limit")
			if err != nil {
				abortWithError(errors.WithStack(err))
			}

			migrator := config.MustNewMigrator()

			_, sourceDriver := config.MustCreateSource()
			currentVersion, _, err := migrator.Version()
			if err != nil && err != migrate.ErrNilVersion {
				abortWithError(errors.WithStack(err))
			}
			migrated.LogUpMigrationPlan(config.Logger, sourceDriver, currentVersion, limit)

			err = migrated.UpMigration(migrator, limit)
			switch err {
			case nil:
			case migrate.ErrNoChange:
				// it's fine
				logger.Warnf("no changes\n")
			case os.ErrNotExist:
				// it's unexpected, but should exit with error
				logger.Warnf("no more migrations can be applied\n")
			default:
				abortWithError(errors.Wrap(err, "up"))
			}
			logger.Infof("migration(s) applied\n")
		},
	}
	upCommand.Flags().IntP("limit", "n", -1, "apply limit, defaults to apply all")
	c.AddCommand(upCommand)

	c.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "rollback one migration",
		Run: func(c *cobra.Command, args []string) {
			migrator := config.MustNewMigrator()
			currentVersion, _, err := migrator.Version()
			switch err {
			case nil:
				logger.Infof("rollback from %d\n", currentVersion)
			case migrate.ErrNilVersion:
				// not things to do
			default:
				abortWithError(errors.Wrap(err, "down"))
			}

			err = migrated.DownMigration(migrator)
			switch err {
			case nil:
			case os.ErrNotExist:
				// it's unexpected, but should exit with error
				logger.Warnf("no more migrations can be rollbacked\n")
			default:
				abortWithError(errors.Wrap(err, "down"))
			}
			logger.Infof("migration rollbacked\n")
		},
	})

	c.AddCommand(&cobra.Command{
		Use:   "force",
		Short: "force to migration version",
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			version, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				abortWithError(errors.WithStack(err))
			}

			err = migrated.ForceMigration(config.MustNewMigrator(), int(version))
			if err != nil {
				abortWithError(errors.Wrap(err, "migration.ForceMigration"))
			}
			logger.Warnf("forced to version %d\n", version)
		},
	})

	return c
}
