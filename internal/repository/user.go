package repository

import "github.com/SerGGe43/alertilka/internal/domain"

type User interface {
	GetById(id int64) (*domain.User, error)
	GetByChatId(chatId int64) (*domain.User, error)
	Add(user domain.User) (int64, error)
	GetState(chatID int64) (int64, error)
	SetState(chatID, state int64) error
}
