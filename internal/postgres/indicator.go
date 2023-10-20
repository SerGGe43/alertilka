package postgres

import (
	"database/sql"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

type Indicator struct {
	db *sql.DB
}

func (i *Indicator) Add(indicator domain.Indicator) (int64, error) {
	query := `INSERT INTO indicator (alertID, indicatorid, value) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := i.db.QueryRow(query, indicator.AlertID, indicator.IndicatorID, indicator.Value).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't add indicator: %w", err)
	}
	return id, nil
}

func (i *Indicator) GetById(id int64) (*domain.Indicator, error) {
	query := `SELECT * FROM indicator WHERE id = $1`
	indicator := domain.Indicator{}
	err := i.db.QueryRow(query, id).Scan(&indicator.Id, &indicator.AlertID,
		&indicator.IndicatorID, &indicator.Value)
	if err != nil {
		return nil, fmt.Errorf("can't get indicator by id: %w", err)
	}
	return &indicator, nil
}

func (i *Indicator) GetByAlertId(alertID int64) ([]domain.Indicator, error) {
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
