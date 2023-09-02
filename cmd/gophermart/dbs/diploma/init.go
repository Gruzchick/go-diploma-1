package diploma

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func Init(databaseURI string) error {
	var err error // TODO ASK: Как сделать чтобы не объявлять переменную

	db, err = sql.Open("pgx", databaseURI)
	if err != nil {
		return err
	}

	defer db.Close() // TODO ASK: Узнать когда это закрывать

	if err = migrate(); err != nil {
		return err
	}

	return nil
}
