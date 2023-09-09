package repository

import "github.com/SerGGe43/alertilka/internal/domain"

type Indicator interface {
	GetByID(id int) (*domain.Indicator, error)
	GetByAlertID(id int) ([]domain.Indicator, error)
}
