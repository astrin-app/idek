package db

import (
	"database/sql"
	"fmt"

	_ "github.com/libsql/libsql-client-go/libsql"
)

func CreateDB(url string) (*sql.DB, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TODOS (
		ID INTEGER PRIMARY KEY,
		TASK TEXT NOT NULL,
		COMPLETED BOOLEAN DEFAULT FALSE
	);
	
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
