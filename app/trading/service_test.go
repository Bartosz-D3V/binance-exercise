package trading

import (
	"context"
	"fmt"
	"github.com/Bartosz-D3V/binance-exercise/app/binance"
	"github.com/Bartosz-D3V/binance-exercise/app/calculation"
	"github.com/Bartosz-D3V/binance-exercise/app/mock"
	"github.com/Bartosz-D3V/binance-exercise/app/transaction"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	mockRepo       *mock.MockRepository
	calcSvc        calculation.Service
	transactionSvc Service

	quantityToSell decimal.Decimal
	minimumBid     decimal.Decimal
	stream         chan binance.BookTickerResponse
	out            chan []transaction.LogEntry
}

func TestRunTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())

	suite.quantityToSell = decimal.NewFromFloat(100.1)
	suite.minimumBid = decimal.NewFromFloat(50.2)
	suite.stream = make(chan binance.BookTickerResponse, 5)
	suite.out = make(chan []transaction.LogEntry, 5)

	suite.mockRepo = mock.NewMockRepository(ctrl)
	suite.calcSvc = calculation.New(suite.quantityToSell, suite.minimumBid)
	suite.transactionSvc = New(suite.calcSvc, suite.mockRepo)
}

func (suite *ServiceTestSuite) TestService_ProcessStockTick_ShouldReturnIfNoSellAllowanceLeft() {
	logEntries := []transaction.LogEntry{
		{Quantity: "10.12"},
		{Quantity: "5.33"},
		{Quantity: "34.73"},
		{Quantity: "49.92"},
	}

	suite.stream <- binance.BookTickerResponse{}
	close(suite.stream)

	suite.mockRepo.EXPECT().GetAll(context.TODO()).Times(1).Return(logEntries)
	suite.mockRepo.EXPECT().Save(context.TODO(), gomock.Any()).Times(0)

	suite.transactionSvc.ProcessStockTick(context.TODO(), suite.stream, suite.out)

	res := <-suite.out
	suite.Equal(logEntries, res)
}

func (suite *ServiceTestSuite) TestService_ProcessStockTick_ShouldIgnoreTicksThatDoesNotMeetPriceCriteria() {
	logEntries := []transaction.LogEntry{
		{Quantity: "10.12"},
	}

	bestBidQuantity := "12.112"
	bestBidPrice := suite.minimumBid.Sub(decimal.NewFromInt(10)).String()
	suite.stream <- binance.BookTickerResponse{BestBidQuantity: bestBidQuantity, BestBidPrice: bestBidPrice}
	close(suite.stream)

	suite.mockRepo.EXPECT().GetAll(context.TODO()).Return(logEntries)
	suite.mockRepo.EXPECT().Save(context.TODO(), gomock.Any()).MaxTimes(0)

	suite.transactionSvc.ProcessStockTick(context.TODO(), suite.stream, suite.out)
	suite.Equal(0, len(suite.out))
}

func (suite *ServiceTestSuite) TestService_ProcessStockTick_ShouldMakeTransactionsThatMeetPriceCriteria() {
	logEntries := []transaction.LogEntry{
		{Quantity: "10.12"},
	}
	bestBidPrice := suite.minimumBid.Add(decimal.NewFromInt(10)).String()
	bestBidQuantity := "89.98"
	suite.stream <- binance.BookTickerResponse{BestBidQuantity: bestBidQuantity, BestBidPrice: bestBidPrice}
	close(suite.stream)

	suite.mockRepo.EXPECT().GetAll(context.TODO()).Return(logEntries)
	suite.mockRepo.EXPECT().Save(context.TODO(), logEntryMatcher{bestBidQuantity, bestBidPrice}).Times(1)

	suite.transactionSvc.ProcessStockTick(context.TODO(), suite.stream, suite.out)
}

type logEntryMatcher struct {
	bestBidQuantity string
	bestBidPrice    string
}

func (e logEntryMatcher) Matches(x interface{}) bool {
	gotLogEntry, ok := x.(transaction.LogEntry)
	return ok && e.bestBidQuantity == gotLogEntry.Quantity &&
		e.bestBidPrice == gotLogEntry.Price
}

func (e logEntryMatcher) String() string {
	return fmt.Sprintf("logEntryMatcher bestBidQuantity=%s and bestBidPrice=%s", e.bestBidQuantity, e.bestBidPrice)
}
