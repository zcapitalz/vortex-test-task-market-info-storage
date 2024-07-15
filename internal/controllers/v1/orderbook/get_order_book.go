package orderbookcontroller

import (
	"market-info-storage/internal/controllers/httputils"
	"market-info-storage/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getOrderBookRequest struct {
	ExchangeName string `uri:"exchange" binding:"required"`
	Pair         string `uri:"pair" binding:"required"`
}

type getOrderBookResponse struct {
	OrderBook []domain.DepthOrder `json:"order_book"`
}

// getOrderBook godoc
// @Summary Get Order Book
// @Description Retrieves the order book for a specific exchange and pair.
// @Tags OrderBook
// @Produce json
// @Param exchange path string true "Exchange name"
// @Param pair path string true "Currency Pair"
// @Success 200 {object} getOrderBookResponse "Order Book data"
// @Failure 400 {object} httputils.HTTPError "Invalid request body"
// @Failure 404 {object} httputils.HTTPError "Order Book not found"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /exchanges/{exchange}/pairs/{pair}/order-book [get]
func (c *OrderBookController) getOrderBook(ctx *gin.Context) {
	var req getOrderBookRequest
	err := ctx.BindUri(&req)
	if err != nil {
		httputils.BindURIError(ctx, err)
		return
	}

	orderBook, err := c.orderBookService.GetOrderBook(req.ExchangeName, req.Pair)
	switch err.(type) {
	case nil:
	case domain.OrderBookNotFound:
		httputils.NotFoundError(ctx, err)
		return
	default:
		httputils.InternalError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, getOrderBookResponse{
		OrderBook: orderBook,
	})
}
