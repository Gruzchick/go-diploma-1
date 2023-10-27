package diplomadb

import "fmt"

func migrate() {
	_, err := DB.Exec(`
		CREATE TABLE users
		(
			id BIGSERIAL PRIMARY KEY,
			login TEXT NOT NULL,
			password TEXT NOT NULL
		)`,
	)
	if err != nil {
		fmt.Println(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE orders
		(
			id BIGSERIAL PRIMARY KEY,
			userId BIGSERIAL
		)`,
	)
	if err != nil {
		fmt.Println(err)
	}

	_, err = DB.Exec(`
		CREATE TABLE withdrawals
		(
			id BIGSERIAL PRIMARY KEY,
			userId BIGSERIAL,
			sum FLOAT8,
			orderId TEXT
		)`,
	)
	if err != nil {
		fmt.Println(err)
	}
}
