package diploma

func migrate() error {
	_, err := db.Exec(`
CREATE TABLE users
(
	id BIGSERIAL PRIMARY KEY,
	login TEXT NOT NULL,
	password TEXT NOT NULL
)`,
	)
	if err != nil {
		return err
	}

	return nil
}
