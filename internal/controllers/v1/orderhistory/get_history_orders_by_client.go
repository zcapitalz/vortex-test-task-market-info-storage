package orderhistorycontroller

import (
	"market-info-storage/internal/controllers/httputils"
	"market-info-storage/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getOrderHistoryResponse struct {
	HistoryOrders []domain.HistoryOrder `json:"orders"`
}

// getOrderHistory godoc
// @Summary Get order history for a client
// @Description Returns a list of orders for the specified client.
// @Tags OrderHistory
// @Accept json
// @Produce json
// @Param client-name query string true "Client name"
// @Param exchange query string true "Exchange name"
// @Param label query string true "Label"
// @Param pair query string true "Currency pair"
// @Success 200 {object} getOrderHistoryResponse
// @Failure 400 {object} httputils.HTTPError "Invalid client data"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /order-history [get]
func (c *OrderHistoryController) getHistoryOrdersByClient(ctx *gin.Context) {
	var client domain.Client
	err := ctx.BindQuery(&client)
	if err != nil {
		httputils.BindQueryError(ctx, err)
	}

	orders, err := c.orderHistoryService.GetHistoryOrdersByClient(&client)
	if err != nil {
		httputils.InternalError(ctx)
		return
	}
	if orders == nil {
		orders = []domain.HistoryOrder{}
	}

	ctx.JSON(http.StatusOK, getOrderHistoryResponse{
		HistoryOrders: orders,
	})
}
