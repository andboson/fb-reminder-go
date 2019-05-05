package migrations

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	bindata "github.com/golang-migrate/migrate/source/go_bindata"
)

//go:generate sh -c "cd ../migrations && $GOPATH/bin/go-bindata -pkg migrations ."

func Migrate(db *sql.DB) error {

	s := bindata.Resource(AssetNames(),
		func(name string) ([]byte, error) {
			return Asset(name)
		})

	sourceDriver, err := bindata.WithInstance(s)
	if err != nil {
		log.Printf("err sourceDriver init: %s", err)
		return err
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Printf("err dbDriver init: %s", err)
		return err
	}

	m, err := migrate.NewWithInstance("go-bindata", sourceDriver, "postgres", dbDriver)
	if err != nil {
		log.Printf("err init instancet: %s", err)
		return err
	}
	err = m.Up()

	switch err {
	case migrate.ErrNoChange:
		return nil
	default:
		return err
	}
}
