package repository

import "github.com/SerGGe43/alertilka/internal/domain"

type Indicator interface {
	Add(indicator domain.Indicator) (int64, error)
	AddID(id domain.TrackingType) error
	GetByID(id int64) (*domain.Indicator, error)
	GetByAlertID(id int64) ([]domain.Indicator, error)
}
