package extractor

import "github.com/shopspring/decimal"

// json to return
type Transaction struct {
	Origin 		string 	`json:"origin"`
	PostingDate	string	`json:"posting_date"`
	ISOPostingDate string `json:"posting_date_iso"`
	Date 		string 	`json:"transaction_date"`
	ISODate 	string 	`json:"transaction_date_iso"`
	Action 		string 	`json:"transaction_action"`
	Amount 		decimal.Decimal `json:"amount"`
	Balance 	decimal.Decimal `json:"balance"`
	Beneficiary string 	`json:"transaction_beneficiary"`
	Description string 	`json:"description"`
	Method 		string 	`json:"transaction_method"`
}

type Data struct {
	Source string `json:"source"`
	AccountNumber string `json:"account_number"`
	Year string `json:"year"`
	Month string `json:"month"`
	Account string `json:"account"`
	AccountType string `json:"account_type"`
	StartingBalance decimal.Decimal `json:"starting_balance"`

	EndingBalance decimal.Decimal `json:"ending_balance"`
	ParsedEndingBalance decimal.Decimal `json:"parsed_ending_balance"`

	TotalDebit	decimal.Decimal `json:"total_debit"`
	TotalCredit	decimal.Decimal `json:"total_credit"`

	Transactions []Transaction `json:"transactions"`
}

func (d *Data) SetAccountType(account string, acc_type string) {
	d.Account = account
	d.AccountType = acc_type
}

func (d *Data) SetYearAndMonth(y string, m string) {
	d.Year = y
	d.Month = m
}

func (d *Data) AddTransactions(t []Transaction) {
	d.Transactions = t
}

func (d *Data) SetStartingBalance(v decimal.Decimal) {
	d.StartingBalance = v
}

func (d *Data) SetEndingBalance(v decimal.Decimal) {
	d.EndingBalance = v
}

func (d *Data) SetParsedEndingBalance(v decimal.Decimal) {
	d.ParsedEndingBalance = v
}

func (d *Data) EndingBalanceMatches() (bool) {
	return d.ParsedEndingBalance.Equal(d.EndingBalance)
}

func (d *Data) SetTotalCredit(v decimal.Decimal) {
	d.TotalCredit = v
}

func (d *Data) SetTotalDebit(v decimal.Decimal) {
	d.TotalDebit = v
}