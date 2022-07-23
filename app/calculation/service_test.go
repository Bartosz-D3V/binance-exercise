package calculation

import (
	"github.com/Bartosz-D3V/binance-exercise/app/transaction"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_TickMeetsTransactionCriteria(t *testing.T) {
	t.Parallel()

	minimumBid := decimal.NewFromFloat(100.0)
	svc := New(decimal.NewFromFloat(90.0), minimumBid)

	tests := []struct {
		name         string
		bestBidPrice decimal.Decimal
		exp          bool
	}{
		{name: "bestBidPrice == minimumBid", bestBidPrice: minimumBid, exp: true},
		{name: "bestBidPrice > minimumBid #1", bestBidPrice: minimumBid.Add(decimal.NewFromFloat(1.0)), exp: true},
		{name: "bestBidPrice > minimumBid #2", bestBidPrice: minimumBid.Add(decimal.NewFromFloat(200.1233)), exp: true},
		{name: "bestBidPrice < minimumBid #1", bestBidPrice: minimumBid.Sub(decimal.NewFromFloat(0.001)), exp: false},
		{name: "bestBidPrice < minimumBid #2", bestBidPrice: minimumBid.Sub(decimal.NewFromFloat(100.2)), exp: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.exp, svc.TickMeetsTransactionCriteria(test.bestBidPrice))
		})
	}
}

func TestService_GetQuantityToSellLeft(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		quantityToSell decimal.Decimal
		logEntries     []transaction.LogEntry
		exp            decimal.Decimal
	}{
		{
			name:           "totalQuantitySold == quantityToSell",
			quantityToSell: decimal.NewFromInt(100),
			logEntries:     []transaction.LogEntry{{Quantity: "50"}, {Quantity: "50"}},
			exp:            decimal.NewFromInt(0),
		},
		{
			name:           "totalQuantitySold < quantityToSell #1",
			quantityToSell: decimal.NewFromFloat(.66655),
			logEntries:     []transaction.LogEntry{{Quantity: ".66644"}},
			exp:            decimal.NewFromFloat(.00011),
		},
		{
			name:           "totalQuantitySold < quantityToSell #2",
			quantityToSell: decimal.NewFromFloat(.66655),
			logEntries:     []transaction.LogEntry{{Quantity: ".66622"}, {Quantity: ".00022"}},
			exp:            decimal.NewFromFloat(.00011),
		},
		{
			name:           "totalQuantitySold < quantityToSell #3",
			quantityToSell: decimal.NewFromFloat(100),
			logEntries:     []transaction.LogEntry{{Quantity: ".99"}},
			exp:            decimal.NewFromFloat(99.01),
		},
		{
			name:           "totalQuantitySold < quantityToSell #4",
			quantityToSell: decimal.NewFromFloat(5.67),
			logEntries:     []transaction.LogEntry{{Quantity: "2.2"}},
			exp:            decimal.NewFromFloat(3.47),
		},
		{
			name:           "totalQuantitySold = quantityToSell = 0",
			quantityToSell: decimal.NewFromInt(0),
			logEntries:     []transaction.LogEntry{{Quantity: "0"}},
			exp:            decimal.NewFromInt(0),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			svc := New(test.quantityToSell, decimal.NewFromFloat(15.66))
			res := svc.GetQuantityToSellLeft(test.logEntries)

			assert.True(t, test.exp.Equal(res))
		})
	}
}

func TestService_GetQuantityToSellPerTick(t *testing.T) {
	t.Parallel()

	svc := New(decimal.NewFromFloat(90.0), decimal.NewFromFloat(100.0))

	tests := []struct {
		name              string
		sellAllowanceLeft decimal.Decimal
		bestBidQuantity   decimal.Decimal
		exp               decimal.Decimal
	}{
		{
			name:              "bestBid = 0 and sellAllowance > 0",
			sellAllowanceLeft: decimal.NewFromFloat(50.1),
			bestBidQuantity:   decimal.NewFromFloat(0),
			exp:               decimal.NewFromFloat(0.0),
		},
		{
			name:              "bestBid > 0 and sellAllowance > 0 #1",
			sellAllowanceLeft: decimal.NewFromFloat(50.1),
			bestBidQuantity:   decimal.NewFromFloat(.1),
			exp:               decimal.NewFromFloat(.1),
		},
		{
			name:              "bestBid > 0 and sellAllowance > 0 #2",
			sellAllowanceLeft: decimal.NewFromFloat(50),
			bestBidQuantity:   decimal.NewFromFloat(5),
			exp:               decimal.NewFromFloat(5),
		},
		{
			name:              "bestBid > 0 and sellAllowance > 0 #3",
			sellAllowanceLeft: decimal.NewFromFloat(10),
			bestBidQuantity:   decimal.NewFromFloat(50.1),
			exp:               decimal.NewFromFloat(10),
		},
		{
			name:              "bestBid > 0 0 and sellAllowance = 0 #1",
			sellAllowanceLeft: decimal.NewFromFloat(0),
			bestBidQuantity:   decimal.NewFromFloat(50.1),
			exp:               decimal.NewFromFloat(0),
		},
		{
			name:              "bestBid = sellAllowance",
			sellAllowanceLeft: decimal.NewFromFloat(15.667),
			bestBidQuantity:   decimal.NewFromFloat(15.667),
			exp:               decimal.NewFromFloat(15.667),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res := svc.GetQuantityToSellPerTick(test.bestBidQuantity, test.sellAllowanceLeft)
			assert.True(t, test.exp.Equal(res))
		})
	}
}
