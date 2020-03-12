package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func initDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	db.Exec(`
CREATE TABLE IF NOT EXISTS greetings (
	id serial PRIMARY KEY,
	greetedPerson text UNIQUE,
	greeting text
)
	`)
}

func saveGreeting(greetedPerson string, greeting string) error {
	_, err := db.Exec(`
INSERT INTO greetings (greetedPerson, greeting)
VALUES ($1, $2)
ON CONFLICT (greetedPerson) 
DO UPDATE 
SET greeting = $2`,
		greetedPerson, greeting)

	return err
}

func getGreeting(greetedPerson string) (*string, error) {
	rows, err := db.Query("SELECT greeting FROM greetings WHERE greetedPerson = $1", greetedPerson)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	greeting := ""
	for rows.Next() {
		err := rows.Scan(&greeting)
		if err != nil {
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &greeting, nil
}
