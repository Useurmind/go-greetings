package pgsql

import (	
	"database/sql"
	_ "github.com/lib/pq"
)

type PGSqlDBContext struct {
	db *sql.DB
}

func NewPGSqlDBContext(dataSource string) (*PGSqlDBContext, error) {
	db, err := getDB(dataSource)
	if err != nil {
		return nil, err
	}

	ctx := &PGSqlDBContext{
		db: db,
	}

	err = ctx.initDB()
	if err != nil {
		return nil, err
	}
	
	return ctx, nil
}

func getDB(dataSource string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (ctx *PGSqlDBContext) initDB() error {
	_, err := ctx.db.Exec(`
CREATE TABLE IF NOT EXISTS greetings (
	id serial PRIMARY KEY,
	greetedPerson text UNIQUE,
	greeting text
)
		`)

	return err
}

func (ctx *PGSqlDBContext) SaveGreeting(greetedPerson string, greeting string) error {
	_, err := ctx.db.Exec(`
INSERT INTO greetings (greetedPerson, greeting)
VALUES ($1, $2)
ON CONFLICT (greetedPerson) 
DO UPDATE 
SET greeting = $2`,
		greetedPerson, greeting)

	return err
}

func (ctx *PGSqlDBContext) GetGreeting(greetedPerson string) (*string, error) {
	rows, err := ctx.db.Query("SELECT greeting FROM greetings WHERE greetedPerson = $1", greetedPerson)
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
