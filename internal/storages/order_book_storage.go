package storages

import (
	"database/sql"
	"fmt"
	"log/slog"
	"market-info-storage/internal/domain"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type OrderBookStorage struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func NewOrderBookStorage(db *sqlx.DB) *OrderBookStorage {
	return &OrderBookStorage{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (s *OrderBookStorage) SaveOrderBook(exchangeName string, pair string, bids, asks []domain.DepthOrder) error {
	bidsExpr := sq.Expr(buildDepthOrderArrayExprSQL(bids), flattenDepthOrders(bids)...)
	asksExpr := sq.Expr(buildDepthOrderArrayExprSQL(asks), flattenDepthOrders(asks)...)
	builder := s.builder.
		Insert("order_books").
		Columns("exchange, pair, bids, asks").
		Values(exchangeName, pair, bidsExpr, asksExpr).
		Suffix(`ON CONFLICT (exchange, pair)
				DO UPDATE SET bids = ?, asks = ?`,
			bidsExpr, asksExpr)

	query, args, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "execute query")
	}

	return nil
}

func (s *OrderBookStorage) GetOrderBook(exchangeName string, pair string) (bids, asks []domain.DepthOrder, err error) {
	builder := s.builder.
		Select("bids, asks").
		From("order_books").
		Where(sq.And{sq.Eq{"exchange": exchangeName}, sq.Eq{"pair": pair}})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, nil, errors.Wrap(err, "build query")
	}
	slog.Debug(fmt.Sprintf("SQL query: %s", query))

	var bidsArray, asksArray DepthOrders
	err = s.db.QueryRowx(query, args...).Scan(&bidsArray, &asksArray)
	switch err {
	case nil:
	case sql.ErrNoRows:
		return nil, nil, domain.OrderBookNotFound{}
	default:
		return nil, nil, errors.Wrap(err, "execute query")
	}

	return bidsArray, asksArray, nil
}

func buildDepthOrderArrayExprSQL(depthOrders []domain.DepthOrder) string {
	expStrSlice := make([]string, 0, len(depthOrders))
	for i := 0; i < len(depthOrders); i++ {
		expStrSlice = append(expStrSlice, "ROW(?, ?)::depth_order")
	}

	return fmt.Sprintf("ARRAY[%s]", strings.Join(expStrSlice, ", "))
}

func flattenDepthOrders(depthOrders []domain.DepthOrder) []any {
	flatten := make([]any, 0, len(depthOrders))
	for _, depthOrder := range depthOrders {
		flatten = append(flatten, depthOrder.Price, depthOrder.BaseQty)
	}
	return flatten
}
