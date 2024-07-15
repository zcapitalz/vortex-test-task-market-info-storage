package app

import (
	"context"
	"fmt"
	"log/slog"
	_ "market-info-storage/api/v1"
	"market-info-storage/internal/config"
	orderbookcontroller "market-info-storage/internal/controllers/v1/orderbook"
	orderhistorycontroller "market-info-storage/internal/controllers/v1/orderhistory"
	"market-info-storage/internal/db/clickhouse"
	"market-info-storage/internal/db/postgres"
	"market-info-storage/internal/domain"
	"market-info-storage/internal/storages"
	"market-info-storage/internal/utils/slogutils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title           Market info storage
// @version         1.0
// @description     API to store and retreive market data

// @BasePath  /api/v1
func Run(cfg config.Config) {
	logger := mustNewLogger(cfg.Env)
	slog.SetDefault(logger)

	slog.Info("Setting up server dependencies")

	postgresClient, err := postgres.NewClient(cfg.Postgres)
	if err != nil {
		slog.Error("initialize Postgres client", slogutils.ErrorAttr(err))
		return
	}
	clickhouseClient, err := clickhouse.NewClient(cfg.ClickHouse)
	if err != nil {
		slog.Error("initialize ClickHouse client", slogutils.ErrorAttr(err))
		return
	}

	orderBookStorage := storages.NewOrderBookStorage(postgresClient)
	historyOrderStorage := storages.NewHistoryOrderStorage(clickhouseClient)

	orderBookService := domain.NewOrderBookService(orderBookStorage)
	orderHistoryService := domain.NewOrderHistoryService(historyOrderStorage)

	orderBookController := orderbookcontroller.NewOrderBookController(orderBookService)
	orderHistoryController := orderhistorycontroller.NewOrderHistoryController(orderHistoryService)

	switch cfg.Env {
	case config.EnvLocal:
		gin.SetMode(gin.DebugMode)
	case config.EnvProd:
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(sloggin.New(logger))
	engine.Use(gin.Recovery())
	engine.GET("api/v1/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))
	orderBookController.RegisterRoutes(engine)
	orderHistoryController.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:    cfg.HTTPServer.IpAddress + ":" + cfg.HTTPServer.Port,
		Handler: engine.Handler(),
	}

	slog.Info("Starting server ...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen: %s\n", slogutils.ErrorAttr(err))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server Shutdown:", slogutils.ErrorAttr(err))
		os.Exit(1)
	}

	select {
	case <-ctx.Done():
		slog.Info("timeout of 5 seconds.")
	}
	slog.Info("Server exiting")
}

func mustNewLogger(env config.Env) (logger *slog.Logger) {
	switch env {
	case config.EnvLocal:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case config.EnvDev, config.EnvProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		panic(fmt.Errorf("unknown env: %v", env))
	}

	return logger
}
