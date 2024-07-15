package orderhistorycontroller

import (
	"market-info-storage/internal/controllers"
	"market-info-storage/internal/domain"

	"github.com/gin-gonic/gin"
)

type OrderHistoryController struct {
	orderHistoryService OrderHistoryService
}

//go:generate mockery --name OrderHistoryService --filename order_history_service.go
type OrderHistoryService interface {
	SaveHistoryOrder(historyOrder *domain.HistoryOrder) error
	GetHistoryOrdersByClient(client *domain.Client) ([]domain.HistoryOrder, error)
}

func NewOrderHistoryController(orderHistoryService OrderHistoryService) controllers.Controller {
	return &OrderHistoryController{
		orderHistoryService: orderHistoryService,
	}
}

func (c *OrderHistoryController) RegisterRoutes(engine *gin.Engine) {
	orderHistoryGroup := engine.Group("/api/v1/order-history")
	orderHistoryGroup.POST("", c.saveHistoryOrder)
	orderHistoryGroup.GET("", c.getHistoryOrdersByClient)
}
