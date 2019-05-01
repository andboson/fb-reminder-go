package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func InitDB(c *Config) (*sql.DB, error) {
	var db *sql.DB
	connLine := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.PgUser, c.PgPasswd, c.PgAddress, c.PgDb)
	db, err := sql.Open("postgres", connLine)
	if err != nil {
		fmt.Printf("[conn] %s", connLine)
		return nil, err
	}

	return db, db.Ping()
}
