package orderbookcontroller

import (
	"errors"
	"market-info-storage/internal/controllers/httputils"
	"market-info-storage/internal/domain"

	"github.com/gin-gonic/gin"
)

type saveOrderBookRequestQuery struct {
	ExchangeName string `uri:"exchange" binding:"required"`
	Pair         string `uri:"pair" binding:"required"`
}

type saveOrderBookRequestBody struct {
	OrderBook []domain.DepthOrder `json:"order_book" binding:"required"`
}

// saveOrderBook godoc
// @Summary Save Order Book
// @Description Saves an order book for a specific exchange and pair.
// @Tags OrderBook
// @Accept json
// @Param exchange path string true "Exchange name"
// @Param pair path string true "Currency Pair"
// @Param orderBook body saveOrderBookRequestBody true "Order Book data"
// @Success 200
// @Failure 400 {object} httputils.HTTPError "Invalid request body"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /exchanges/{exchange}/pairs/{pair}/order-book [put]
func (c *OrderBookController) saveOrderBook(ctx *gin.Context) {
	var reqURI saveOrderBookRequestQuery
	err := ctx.BindUri(&reqURI)
	if err != nil {
		httputils.BindURIError(ctx, err)
		return
	}
	var reqBody saveOrderBookRequestBody
	err = ctx.BindJSON(&reqBody)
	if err != nil {
		httputils.BindJSONBodyError(ctx, err)
		return
	}
	err = validateSaveOrderBookRequestBody(&reqBody)
	if err != nil {
		httputils.BindJSONBodyError(ctx, err)
		return
	}

	err = c.orderBookService.SaveOrderBook(reqURI.ExchangeName, reqURI.Pair, reqBody.OrderBook)
	if err != nil {
		httputils.InternalError(ctx)
		return
	}
}

func validateSaveOrderBookRequestBody(reqBody *saveOrderBookRequestBody) error {
	if len(reqBody.OrderBook)%2 != 0 {
		return errors.New("order book should have even number of entries")
	}
	return nil
}
