package postgres

import (
	"database/sql"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

type Alert struct {
	db *sql.DB
}

func NewAlert(db *sql.DB) *Alert {

	return &Alert{
		db: db,
	}
}

func (a *Alert) Add(alert domain.Alert) (int64, error) {
	query := `INSERT INTO alert (ticker, userID) VALUES ($1, $2) RETURNING id`
	var id int64
	err := a.db.QueryRow(query, alert.Ticker, alert.UserID).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("can't add alert: %w", err)
	}
	return id, nil
}

func (a *Alert) GetByID(id int64) (*domain.Alert, error) {
	query := `SELECT * FROM alert WHERE id = $1`
	alert := domain.Alert{}
	err := a.db.QueryRow(query, id).Scan(&alert.Id, &alert.Ticker, &alert.Name, &alert.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Can't get alert by ID %w", err)
	}
	return &alert, nil
}

func (a *Alert) GetByUserID(userID int64) ([]domain.Alert, error) {
	query := `SELECT * FROM alert WHERE userid = $1`
	var alerts []domain.Alert
	rows, err := a.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("can't initialize rows alert getByUserId: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var alert domain.Alert
		err = rows.Scan(&alert.Id, &alert.Ticker, &alert.Name, &alert.UserID)
		if err != nil {
			return nil, fmt.Errorf("rows next cycle error: %w", err)
		}
		alerts = append(alerts, alert)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return alerts, nil
}

func (a *Alert) AddName(name string, userid int64) error {
	query := `UPDATE alert
				SET name = $1
				WHERE name IS NULL and userid = $2`
	err := a.db.QueryRow(query, name, userid)
	fmt.Println(name)
	if err != nil {
		fmt.Errorf("can't set name to ticker: %w")
	}
	return nil
}
