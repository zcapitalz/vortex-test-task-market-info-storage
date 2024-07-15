package domain

import (
	"log/slog"
	"market-info-storage/internal/utils/slogutils"

	"github.com/pkg/errors"
)

type OrderBookService struct {
	orderBookStorage OrderBookStorage
}

type OrderBookStorage interface {
	SaveOrderBook(exchangeName, pair string, asks, bids []DepthOrder) error
	GetOrderBook(exchangeName, pair string) (bids, asks []DepthOrder, err error)
}

func NewOrderBookService(orderBookStorage OrderBookStorage) *OrderBookService {
	return &OrderBookService{
		orderBookStorage: orderBookStorage,
	}
}

func (s *OrderBookService) SaveOrderBook(exchangeName, pair string, orderBook []DepthOrder) error {
	bids := orderBook[:int(len(orderBook)/2)]
	asks := orderBook[int(len(orderBook)/2):]
	err := s.orderBookStorage.SaveOrderBook(exchangeName, pair, bids, asks)
	if err != nil {
		err = errors.Wrap(err, "save order book")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return err
}

func (s *OrderBookService) GetOrderBook(exchangeName, pair string) (orderBook []DepthOrder, err error) {
	bids, asks, err := s.orderBookStorage.GetOrderBook(exchangeName, pair)
	switch err.(type) {
	case nil:
	case OrderBookNotFound:
		return nil, err
	default:
		err = errors.Wrap(err, "get order book")
		slog.Error("", slogutils.ErrorAttr(err))
	}

	orderBook = append(bids, asks...)
	return orderBook, err
}
