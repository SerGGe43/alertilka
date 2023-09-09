package repository

import "github.com/SerGGe43/alertilka/internal/domain"

type Alert interface {
	GetByID(id int) (*domain.Alert, error)
	GetByUserID(id int) ([]domain.Alert, error)
}
