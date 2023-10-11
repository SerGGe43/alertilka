package bot

const (
	Start         = "/start"
	Help          = "help"
	PriceByTicker = "Price by ticker"
	NewAlert      = "New alert"
)

type Bot interface {
	SendMenu(chatId int64) error
	SendHelp(chatId int64) error
	SendTickerRequest(chatId int64) error
	SendTickerPrices(chatID int64, prices string) error
}
