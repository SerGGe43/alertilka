package postgres

import (
	"database/sql"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

type Indicator struct {
	db *sql.DB
}

func NewIndicator(db *sql.DB) *Indicator {
	return &Indicator{
		db: db,
	}
}

func (i *Indicator) Add(indicator domain.Indicator) (int64, error) {
	query := `INSERT INTO indicator (alertID, value) VALUES ($1, $2) RETURNING id`
	var id int64
	err := i.db.QueryRow(query, indicator.AlertID, indicator.Value).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't add indicator: %w", err)
	}
	return id, nil
}

func (i *Indicator) AddID(id domain.TrackingType) error {
	query := `UPDATE indicator
				SET indicatorid = $1
				WHERE indicatorid IS NULL`
	err := i.db.QueryRow(query, id)
	if err != nil {
		return fmt.Errorf("can't add indicatorID: %w", err)
	}
	return nil
}

func (i *Indicator) AddValue(value float64) error {
	query := `UPDATE indicator
				SET value = $1
				WHERE value = -1`
	err := i.db.QueryRow(query, value)
	if err != nil {
		return fmt.Errorf("can't add indicator value: %w", err)
	}
	return nil
}

func (i *Indicator) GetByID(id int64) (*domain.Indicator, error) {
	query := `SELECT * FROM indicator WHERE id = $1`
	indicator := domain.Indicator{}
	err := i.db.QueryRow(query, id).Scan(&indicator.Id, &indicator.AlertID,
		&indicator.IndicatorID, &indicator.Value)
	if err != nil {
		return nil, fmt.Errorf("can't get indicator by id: %w", err)
	}
	return &indicator, nil
}

func (i *Indicator) GetByAlertID(alertID int64) ([]domain.Indicator, error) {
	query := `SELECT * FROM indicator WHERE alertid = $1`
	rows, err := i.db.Query(query, alertID)
	if err != nil {
		return nil, fmt.Errorf("can't get rows GetByAlertId: %w", err)
	}
	defer rows.Close()
	var indicators []domain.Indicator
	for rows.Next() {
		var indicator domain.Indicator
		err = rows.Scan(&indicator.Id, &indicator.AlertID,
			&indicator.IndicatorID, &indicator.Value)
		if err != nil {
			return nil, fmt.Errorf("rows cycle error: %w", err)
		}
		indicators = append(indicators, indicator)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return indicators, nil
}
