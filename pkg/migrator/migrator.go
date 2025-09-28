package migrator

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Config struct {
	Dsn               string
	DBName            string
	Schema            string
	MigrationsTable   string
	MigrationsDirPath string
}

func Up(c Config) error {
	db, err := sql.Open("pgx", c.Dsn)
	if err != nil {
		return err
	}
	driver, err := pgx.WithInstance(db, &pgx.Config{
		SchemaName:      c.Schema,
		MigrationsTable: c.MigrationsTable,
	})
	if err != nil {
		return err
	}
	defer driver.Close()
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+c.MigrationsDirPath,
		c.DBName,
		driver,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
