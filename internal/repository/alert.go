package repository

import "github.com/SerGGe43/alertilka/internal/domain"

type Alert interface {
	GetByID(id int64) (*domain.Alert, error)
	GetByUserID(id int64) ([]domain.Alert, error)
	Add(alert domain.Alert) (int64, error)
}
