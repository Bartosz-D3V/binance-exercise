package transaction

import "github.com/shopspring/decimal"

type Report struct {
	TotalQuantitySold decimal.Decimal
	TotalMoneyEarned  decimal.Decimal
}
