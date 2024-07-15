package orderbookcontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"market-info-storage/internal/controllers/v1/orderbook/mocks"
	"market-info-storage/internal/domain"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveOrderBook(t *testing.T) {
	exchange := "bybit"
	pair := "MATIC_USDT"
	orderBook := []domain.DepthOrder{
		{Price: 0.53, BaseQty: 1.5},
		{Price: 0.54, BaseQty: 1.1},
	}

	service := mocks.NewOrderBookService(t)
	service.On("SaveOrderBook", exchange, pair, mock.Anything).Return(nil)
	controller := NewOrderBookController(service)

	url := fmt.Sprintf("/api/v1/exchanges/%s/pairs/%s/order-book", exchange, pair)
	reqBodyReader := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyReader).Encode(saveOrderBookRequestBody{OrderBook: orderBook})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, url, reqBodyReader)

	w := httptest.NewRecorder()
	router := gin.Default()
	controller.RegisterRoutes(router)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("response body: %s", w.Body.String()))
}

func TestSaveOrderBookWrongFormat(t *testing.T) {
	testCases := []struct {
		name     string
		exchange string
		pair     string
		body     string
	}{
		{
			name:     "AbsentBaseQty",
			exchange: "bybit",
			pair:     "MATIC_USDT",
			body: `{
				{Price: 0.53},
				{Price: 0.54}
			}`,
		},
		{
			name:     "WrongFieldName",
			exchange: "bybit",
			pair:     "MATIC_USDT",
			body: `{
				{PPPPPPPrice: 0.53, BaseQty: 1.5},
				{PPPPPPPrice: 0.54, BaseQty: 1.1},
			}`,
		},
		{
			name:     "EmptyExchange",
			exchange: "",
			pair:     "MATIC_USDT",
			body: `{
				{Price: 0.53, BaseQty: 1.5},
				{Price: 0.54, BaseQty: 1.1},
			}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := mocks.NewOrderBookService(t)
			controller := NewOrderBookController(service)

			url := fmt.Sprintf("/api/v1/exchanges/%s/pairs/%s/order-book", tc.exchange, tc.pair)
			reqBodyReader := strings.NewReader(tc.body)
			req := httptest.NewRequest(http.MethodPut, url, reqBodyReader)

			w := httptest.NewRecorder()
			router := gin.Default()
			controller.RegisterRoutes(router)
			router.ServeHTTP(w, req)

			require.Equal(t, http.StatusBadRequest, w.Code, fmt.Sprintf("response body: %s", w.Body.String()))
		})
	}
}

func TestSaveOrderBookWithOddNumberOfOrders(t *testing.T) {
	exchange := "bybit"
	pair := "MATIC_USDT"
	orderBook := []domain.DepthOrder{
		{Price: 0.53, BaseQty: 1.5},
		{Price: 0.54, BaseQty: 1.1},
		{Price: 0.55, BaseQty: 0.9},
	}

	service := mocks.NewOrderBookService(t)
	controller := NewOrderBookController(service)

	url := fmt.Sprintf("/api/v1/exchanges/%s/pairs/%s/order-book", exchange, pair)
	reqBodyReader := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyReader).Encode(saveOrderBookRequestBody{OrderBook: orderBook})
	require.NoError(t, err)
	req := httptest.NewRequest(http.MethodPut, url, reqBodyReader)

	w := httptest.NewRecorder()
	router := gin.Default()
	controller.RegisterRoutes(router)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code, fmt.Sprintf("response body: %s", w.Body.String()))
}

func TestGetOrderBook(t *testing.T) {
	exchange := "binance"
	pair := "SOL_USDT"
	orderBook := []domain.DepthOrder{
		{Price: 0.53, BaseQty: 1.5},
		{Price: 0.54, BaseQty: 1.1},
	}

	service := mocks.NewOrderBookService(t)
	service.On("GetOrderBook", "binance", "SOL_USDT").Return(orderBook, nil)
	controller := NewOrderBookController(service)

	url := fmt.Sprintf("/api/v1/exchanges/%s/pairs/%s/order-book", exchange, pair)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	w := httptest.NewRecorder()
	router := gin.Default()
	controller.RegisterRoutes(router)
	router.ServeHTTP(w, req)

	var respBody getOrderBookResponse
	err := json.Unmarshal(w.Body.Bytes(), &respBody)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("response body: %s", w.Body.String()))
	require.True(t, reflect.DeepEqual(orderBook, respBody.OrderBook))
}

func TestGetNonExistentOrderBook(t *testing.T) {
	exchange := "binance"
	pair := "SOL_USDT"

	service := mocks.NewOrderBookService(t)
	service.On("GetOrderBook", "binance", "SOL_USDT").Return(nil, domain.OrderBookNotFound{})
	controller := NewOrderBookController(service)

	url := fmt.Sprintf("/api/v1/exchanges/%s/pairs/%s/order-book", exchange, pair)
	req := httptest.NewRequest(http.MethodGet, url, nil)

	w := httptest.NewRecorder()
	router := gin.Default()
	controller.RegisterRoutes(router)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusNotFound, w.Code)
}
