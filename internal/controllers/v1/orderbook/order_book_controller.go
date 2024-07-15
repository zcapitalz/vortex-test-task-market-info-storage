package orderbookcontroller

import (
	"market-info-storage/internal/domain"

	"github.com/gin-gonic/gin"
)

type OrderBookController struct {
	orderBookService OrderBookService
}

//go:generate mockery --name OrderBookService --filename order_book_service.go
type OrderBookService interface {
	SaveOrderBook(exchange_name, pair string, orderBook []domain.DepthOrder) error
	GetOrderBook(exchange_name, pair string) ([]domain.DepthOrder, error)
}

func NewOrderBookController(orderBookService OrderBookService) *OrderBookController {
	return &OrderBookController{
		orderBookService: orderBookService,
	}
}

func (c *OrderBookController) RegisterRoutes(engine *gin.Engine) {
	orderBookGroup := engine.Group("/api/v1/exchanges/:exchange/pairs/:pair/order-book")
	orderBookGroup.PUT("", c.saveOrderBook)
	orderBookGroup.GET("", c.getOrderBook)
}
