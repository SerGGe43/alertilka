package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func NewConnection(dbInfo string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping db: %w", err)
	}
	return db, nil
}
