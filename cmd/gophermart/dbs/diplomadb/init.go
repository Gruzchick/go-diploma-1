package diplomadb

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Init(databaseURI string) error {
	var err error

	DB, err = sql.Open("pgx", databaseURI)
	if err != nil {
		return err
	}

	migrate()

	return nil
}
