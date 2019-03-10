package migrated

import (
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
)

// LoggerT defines a common logger interface used by migrator.
type LoggerT interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
}

// Migrator configuration object.
type MigratorConfig struct {
	Logger LoggerT

	// base dir of the migration directories
	MigrationBaseDir string
	// migration database driver constructor
	CreateDatabase func() (string, database.Driver, error)
	// migration source constructor
	CreateSource func() (string, source.Driver, error)
}

// MustNewMigrator creates a migrator instance from dsn string.
// It panics on error.
func (c *MigratorConfig) MustNewMigrator() *migrate.Migrate {
	sourceName, source, err := c.CreateSource()
	c.abortIfError(err)
	driverName, driver, err := c.CreateDatabase()
	c.abortIfError(err)
	migrator, err := migrate.NewWithInstance(
		sourceName, source,
		driverName, driver,
	)
	c.abortIfError(err)

	return migrator
}

// MustGetMigrationBaseDir returns migration files base dir.
// It panics on error.
func (c *MigratorConfig) MustGetMigrationBaseDir() string {
	migrationBaseDir, err := filepath.Abs(c.MigrationBaseDir)
	c.abortIfError(err)
	return migrationBaseDir
}

func (c *MigratorConfig) abortIfError(err error) {
	if err == nil {
		return
	}

	c.Logger.Errorf("%+v\n", err)
	os.Exit(1)
}
