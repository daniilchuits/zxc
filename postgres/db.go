package postgres

import (
	"database/sql"
	"log"
)

func PostgresTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS files(  
	event_id TEXT NOT NULL, 
	type TEXT NOT NULL, 
	event_time TIMESTAMP NOT NULL, 
	byn INT NOT NULL, 
	currency TEXT NOT NULL,  
	originalEvent TEXT,  
	methods TEXT,  
	last4 TEXT, 
	issuer TEXT 
	)
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
	return err
}
