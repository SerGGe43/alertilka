package postgres

import (
	"database/sql"
	"fmt"
	"github.com/SerGGe43/alertilka/internal/domain"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{
		db: db,
	}
}

func (u User) GetById(id int64) (*domain.User, error) {
	query := `SELECT (id, name, chatID) FROM "user" WHERE id = $1`
	user := domain.User{}
	err := u.db.QueryRow(query, id).Scan(&user.Id, &user.Name, &user.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Can't get user by ID %w", err)
	}
	return &user, nil
}

func (u User) GetByChatId(chatID int64) (*domain.User, error) {
	query := `SELECT (id, name, chatID) FROM "user" WHERE chatID = $1`
	user := domain.User{}
	err := u.db.QueryRow(query, chatID).Scan(&user.Id, &user.Name, &user.ChatId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Can't get user by chatID %w", err)
	}
	return &user, nil
}

func (u User) Add(user domain.User) (int64, error) {
	query := `INSERT INTO "user"(name, chatID) VALUES ($1, $2) RETURNING id`
	var id int64
	err := u.db.QueryRow(query, user.Name, user.ChatId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't add user %w", err)
	}
	return id, nil
}
