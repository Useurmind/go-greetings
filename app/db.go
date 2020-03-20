package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var dbInitialized bool
var dataSource string

func setDataSource(ds string) {
	dataSource = ds
}

func getDB() (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if !dbInitialized {
		db.Exec(`
CREATE TABLE IF NOT EXISTS greetings (
	id serial PRIMARY KEY,
	greetedPerson text UNIQUE,
	greeting text
)
		`)
	}

	return db, nil;
}

func saveGreeting(greetedPerson string, greeting string) error {
	db, err := getDB()
	if err != nil {
		return err
	}

	_, err = db.Exec(`
INSERT INTO greetings (greetedPerson, greeting)
VALUES ($1, $2)
ON CONFLICT (greetedPerson) 
DO UPDATE 
SET greeting = $2`,
		greetedPerson, greeting)

	return err
}

func getGreeting(greetedPerson string) (*string, error) {
	db, err := getDB()
	if err != nil {
		return nil, err
	}

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
