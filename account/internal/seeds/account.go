package seeds

import "github.com/0xRichardL/temporal-practice/account/internal/models"

var ACCOUNTS = []models.Account{
	/// Normal account.
	{ID: "1", Balance: 1000},
	/// Normal account.
	{ID: "2", Balance: 1000},
	/// Empty account.
	{ID: "3", Balance: 0},
	/// Big balance account, reached the maximum of int64.
	{ID: "4", Balance: 9223372036854775807},
}
