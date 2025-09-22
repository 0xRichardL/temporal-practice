package models

type Account struct {
	ID      string
	Balance int64
}

func (Account) TableName() string {
	return "accounts"
}
