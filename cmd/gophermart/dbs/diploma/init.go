package diploma

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Init(databaseURI string) error {
	var err error // TODO ASK: Как сделать чтобы не объявлять переменную

	DB, err = sql.Open("pgx", databaseURI)
	if err != nil {
		return err
	}

	//defer DB.Close() // TODO ASK: Узнать когда это закрывать

	migrate()

	return nil
}
