package domain

const (
	MAIN_MENU           = 0
	NEW_ALERT           = 1
	ADD_INDICATOR_ID    = 2
	PRICE_BY_TICKER     = 3
	ADD_NAME_TO_ALERT   = 4
	ADD_INDICATOR_VALUE = 5
)

type User struct {
	Id     int64
	Name   string
	ChatId int64
	State  int64
}
