package domain

import (
	"log/slog"
	"market-info-storage/internal/utils/slogutils"

	"github.com/pkg/errors"
)

type OrderHistoryService struct {
	orderHistoryStorage OrderHistoryStorage
}

type OrderHistoryStorage interface {
	SaveHistoryOrder(order *HistoryOrder) error
	GetHistoryOrdersByClient(client *Client) ([]HistoryOrder, error)
}

func NewOrderHistoryService(orderHistoryStorage OrderHistoryStorage) *OrderHistoryService {
	return &OrderHistoryService{
		orderHistoryStorage: orderHistoryStorage,
	}
}

func (s *OrderHistoryService) SaveHistoryOrder(order *HistoryOrder) error {
	err := s.orderHistoryStorage.SaveHistoryOrder(order)
	if err != nil {
		err = errors.Wrap(err, "save order")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return err
}

func (s *OrderHistoryService) GetHistoryOrdersByClient(client *Client) ([]HistoryOrder, error) {
	orderHistory, err := s.orderHistoryStorage.GetHistoryOrdersByClient(client)
	if err != nil {
		err = errors.Wrap(err, "get order history")
		slog.Error("", slogutils.ErrorAttr(err))
	}
	return orderHistory, err
}
