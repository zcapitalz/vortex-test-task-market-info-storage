package orderhistorycontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"market-info-storage/internal/controllers/v1/orderhistory/mocks"
	"market-info-storage/internal/domain"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestSaveHistoryOrder(t *testing.T) {
	historyOrder := newHistoryOrderToSave()
	historyOrder.Side = "buy"
	historyOrder.Type = "market"
	*historyOrder.BaseQty = 10.0
	*historyOrder.Price = 100.0
	historyOrder.AlgorithmNamePlaced = "MyAlgorithm"
	*historyOrder.LowestSellPrc = 99.0
	*historyOrder.HighestBuyPrc = 101.0
	*historyOrder.CommissionQuoteQty = 0.1
	historyOrder.TimePlaced = time.Now().Truncate(time.Nanosecond)
	reqQuery := saveOrderRequestQuery{
		ClientName:   "John Doe",
		ExchangeName: "binance",
		Label:        "My Order",
		Pair:         "BTCUSDT",
	}

	service := mocks.NewOrderHistoryService(t)
	service.On("SaveHistoryOrder", &domain.HistoryOrder{
		ClientName:          reqQuery.ClientName,
		ExchangeName:        reqQuery.ExchangeName,
		Label:               reqQuery.Label,
		Pair:                reqQuery.Pair,
		Side:                historyOrder.Side,
		Type:                historyOrder.Type,
		BaseQty:             *historyOrder.BaseQty,
		Price:               *historyOrder.Price,
		AlgorithmNamePlaced: historyOrder.AlgorithmNamePlaced,
		LowestSellPrc:       *historyOrder.LowestSellPrc,
		HighestBuyPrc:       *historyOrder.HighestBuyPrc,
		CommissionQuoteQty:  *historyOrder.CommissionQuoteQty,
		TimePlaced:          historyOrder.TimePlaced,
	}).Return(nil)
	controller := NewOrderHistoryController(service)

	reqBodyReader := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyReader).Encode(saveOrderRequestBody{HistoryOrder: *historyOrder})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/order-history", reqBodyReader)
	q := req.URL.Query()
	q.Add("client-name", reqQuery.ClientName)
	q.Add("exchange", reqQuery.ExchangeName)
	q.Add("label", reqQuery.Label)
	q.Add("pair", reqQuery.Pair)
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	router := gin.Default()
	controller.RegisterRoutes(router)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("response body: %s", w.Body.String()))
}

func TestGetOrderHistory(t *testing.T) {
	client := domain.Client{
		ClientName:   "John Doe",
		ExchangeName: "binance",
		Label:        "My Order",
		Pair:         "BTCUSDT",
	}
	historyOrders := []domain.HistoryOrder{
		{
			ClientName:          client.ClientName,
			ExchangeName:        client.ExchangeName,
			Label:               client.Label,
			Pair:                client.Pair,
			Side:                "buy",
			Type:                "market",
			BaseQty:             10.0,
			Price:               100.0,
			AlgorithmNamePlaced: "MyAlgorithm",
			LowestSellPrc:       99.0,
			HighestBuyPrc:       101.0,
			CommissionQuoteQty:  0.1,
			TimePlaced:          time.Now().Truncate(time.Nanosecond),
		},
		{
			ClientName:          client.ClientName,
			ExchangeName:        client.ExchangeName,
			Label:               client.Label,
			Pair:                client.Pair,
			Side:                "sell",
			Type:                "limit",
			BaseQty:             99.0,
			Price:               900.0,
			AlgorithmNamePlaced: "OtherAlgorithm",
			LowestSellPrc:       99.0,
			HighestBuyPrc:       901.0,
			CommissionQuoteQty:  0.9,
			TimePlaced:          time.Now().Truncate(time.Nanosecond),
		},
	}

	service := mocks.NewOrderHistoryService(t)
	service.On("GetHistoryOrdersByClient", &domain.Client{
		ClientName:   client.ClientName,
		ExchangeName: client.ExchangeName,
		Label:        client.Label,
		Pair:         client.Pair,
	}).Return(historyOrders, nil)
	controller := NewOrderHistoryController(service)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/order-history", nil)
	q := req.URL.Query()
	q.Add("client-name", client.ClientName)
	q.Add("exchange", client.ExchangeName)
	q.Add("label", client.Label)
	q.Add("pair", client.Pair)
	req.URL.RawQuery = q.Encode()

	w := httptest.NewRecorder()
	router := gin.Default()
	controller.RegisterRoutes(router)
	router.ServeHTTP(w, req)

	var respBody getOrderHistoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("response body: %s", w.Body.String()))
	require.True(t, reflect.DeepEqual(historyOrders, respBody.HistoryOrders))
}

func newHistoryOrderToSave() *HistoryOrderToSave {
	return &HistoryOrderToSave{
		BaseQty:            new(float64),
		Price:              new(float64),
		LowestSellPrc:      new(float64),
		HighestBuyPrc:      new(float64),
		CommissionQuoteQty: new(float64),
	}
}
