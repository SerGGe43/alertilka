package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upTables, downTables)
}

func upTables(tx *sql.Tx) error {
	query := `
			DO $$
				BEGIN
    				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tracking_type1') THEN
        				CREATE TYPE tracking_type1 AS ENUM ('more than price', 'lower than price', 'more than MA');
    				END IF;
				END
			$$;
			CREATE TABLE IF NOT EXISTS "User" (
			    id SERIAL PRIMARY KEY,
			    name varchar(32),
			    chatID bigint 
			);
			CREATE TABLE IF NOT EXISTS alert (
		    id SERIAL PRIMARY KEY,
		    ticker varchar(4),
		    name varchar(32),
		    userID integer REFERENCES "User"(id) ON DELETE CASCADE 
			);
			CREATE TABLE IF NOT EXISTS indicator (
			id SERIAL PRIMARY KEY,
			alertID int REFERENCES alert(id) ON DELETE CASCADE,
			indicatorID tracking_type1,
			value int
			)`
	_, err := tx.Exec(query)
	if err != nil {
		fmt.Printf("can't exec 0001 migration (up) %w", err)
	}
	return nil
}

func downTables(tx *sql.Tx) error {
	query := `DROP TABLE indicator;
				DROP TABLE alert;
				DROP TABLE "User"
				DROP TYPE tracking_type1`
	_, err := tx.Exec(query)
	if err != nil {
		fmt.Printf("can't exec 0001 migration (down) %w", err)
	}
	return nil
}
