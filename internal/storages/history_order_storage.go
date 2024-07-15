package storages

import (
	"context"
	"market-info-storage/internal/domain"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/pkg/errors"
)

type HistoryOrderStorage struct {
	db driver.Conn
}

func NewHistoryOrderStorage(db driver.Conn) *HistoryOrderStorage {
	return &HistoryOrderStorage{
		db: db,
	}
}

func (s *HistoryOrderStorage) SaveHistoryOrder(order *domain.HistoryOrder) error {
	err := s.db.Exec(context.Background(), `
		INSERT INTO history_orders (
			client_name,
			exchange_name,
			label,
			pair,
			side,
			type,
			base_qty,
			price,
			algorithm_name_placed,
			lowest_sell_prc,
			highest_buy_prc,
			commission_quote_qty,
			time_placed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		order.ClientName, order.ExchangeName, order.Label, order.Pair, order.Side, order.Type, order.BaseQty, order.Price, order.AlgorithmNamePlaced, order.LowestSellPrc, order.HighestBuyPrc, order.CommissionQuoteQty, order.TimePlaced)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *HistoryOrderStorage) GetHistoryOrdersByClient(client *domain.Client) ([]domain.HistoryOrder, error) {
	rows, err := s.db.Query(context.Background(), `
		SELECT
			client_name,
			exchange_name,
			label,
			pair,
			side,
			type,
			base_qty,
			price,
			algorithm_name_placed,
			lowest_sell_prc,
			highest_buy_prc,
			commission_quote_qty,
			time_placed
		FROM history_orders
		WHERE
			client_name = ? AND
			exchange_name = ? AND
			label = ? AND
			pair = ?`,
		client.ClientName, client.ExchangeName, client.Label, client.Pair)
	if err != nil {
		return nil, errors.Wrap(err, "execute query")
	}

	var historyOrders []domain.HistoryOrder
	for rows.Next() {
		var historyOrder domain.HistoryOrder
		err := rows.Scan(
			&historyOrder.ClientName,
			&historyOrder.ExchangeName,
			&historyOrder.Label,
			&historyOrder.Pair,
			&historyOrder.Side,
			&historyOrder.Type,
			&historyOrder.BaseQty,
			&historyOrder.Price,
			&historyOrder.AlgorithmNamePlaced,
			&historyOrder.LowestSellPrc,
			&historyOrder.HighestBuyPrc,
			&historyOrder.CommissionQuoteQty,
			&historyOrder.TimePlaced)
		if err != nil {
			return nil, errors.Wrap(err, "scan values")
		}
		historyOrders = append(historyOrders, historyOrder)
	}

	return historyOrders, nil
}
