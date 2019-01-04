//go:generate go-bindata -pkg model -ignore .*\.go .
//go:generate go fmt .
//go:generate goimports -w -l .
package model

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
	"github.com/jmoiron/sqlx"
)

var ErrNoRows = errors.New("no rows")

func Migrate(db *sqlx.DB) error {
	source, err := bindata.WithInstance(bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	}))
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("go-bindata", source, "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err == nil {
		log.Println("Migrated", err)
	}
	if err == migrate.ErrNoChange {
		return nil
	}
	return err
}
