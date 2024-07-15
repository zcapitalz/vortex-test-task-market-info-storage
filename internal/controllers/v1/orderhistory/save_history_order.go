package orderhistorycontroller

import (
	"market-info-storage/internal/controllers/httputils"
	"market-info-storage/internal/domain"
	"time"

	"github.com/gin-gonic/gin"
)

type saveOrderRequestBody struct {
	HistoryOrder HistoryOrderToSave `json:"historyOrder" binding:"required"`
}

type saveOrderRequestQuery struct {
	ClientName   string `form:"client-name" binding:"required"`
	ExchangeName string `form:"exchange" binding:"required"`
	Label        string `form:"label" binding:"required"`
	Pair         string `form:"pair" binding:"required"`
}

type HistoryOrderToSave struct {
	Side                string    `json:"side" binding:"required"`
	Type                string    `json:"type" binding:"required"`
	BaseQty             *float64  `json:"baseQty" binding:"required"`
	Price               *float64  `json:"price" binding:"required"`
	AlgorithmNamePlaced string    `json:"algorithmNamePlaced" binding:"required"`
	LowestSellPrc       *float64  `json:"lowestSellPrc" binding:"required"`
	HighestBuyPrc       *float64  `json:"highestBuyPrc" binding:"required"`
	CommissionQuoteQty  *float64  `json:"commissionQuoteQty" binding:"required"`
	TimePlaced          time.Time `json:"timePlaced" binding:"required"`
}

// saveOrder godoc
// @Summary Save order
// @Description Saves an order.
// @Tags OrderHistory
// @Accept json
// @Param client-name query string true "Client name"
// @Param exchange query string true "Exchange name"
// @Param label query string true "Label"
// @Param pair query string true "Currency pair"
// @Param history-order body saveOrderRequestBody true "History order"
// @Success 200
// @Failure 400 {object} httputils.HTTPError "Invalid request body"
// @Failure 500 {object} httputils.HTTPError "Internal server error"
// @Router /order-history [post]
func (c *OrderHistoryController) saveHistoryOrder(ctx *gin.Context) {
	var reqQuery saveOrderRequestQuery
	err := ctx.BindQuery(&reqQuery)
	if err != nil {
		httputils.BindQueryError(ctx, err)
	}
	var reqBody saveOrderRequestBody
	err = ctx.BindJSON(&reqBody)
	if err != nil {
		httputils.BindJSONBodyError(ctx, err)
		return
	}

	historyOrder := &domain.HistoryOrder{
		ClientName:          reqQuery.ClientName,
		ExchangeName:        reqQuery.ExchangeName,
		Label:               reqQuery.Label,
		Pair:                reqQuery.Pair,
		Side:                reqBody.HistoryOrder.Side,
		Type:                reqBody.HistoryOrder.Type,
		BaseQty:             *reqBody.HistoryOrder.BaseQty,
		Price:               *reqBody.HistoryOrder.Price,
		AlgorithmNamePlaced: reqBody.HistoryOrder.AlgorithmNamePlaced,
		LowestSellPrc:       *reqBody.HistoryOrder.LowestSellPrc,
		HighestBuyPrc:       *reqBody.HistoryOrder.HighestBuyPrc,
		CommissionQuoteQty:  *reqBody.HistoryOrder.CommissionQuoteQty,
		TimePlaced:          reqBody.HistoryOrder.TimePlaced,
	}
	err = c.orderHistoryService.SaveHistoryOrder(historyOrder)
	if err != nil {
		httputils.InternalError(ctx)
		return
	}
}
