package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/miruken1207/warehouse-api/internal/config"
	"github.com/miruken1207/warehouse-api/internal/database"
	"github.com/miruken1207/warehouse-api/internal/handler"
	"github.com/miruken1207/warehouse-api/internal/middleware"
	"github.com/miruken1207/warehouse-api/internal/repository"
	"github.com/miruken1207/warehouse-api/internal/service"

	_ "github.com/miruken1207/warehouse-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Warehouse API
// @version 1.0
// @description  CRUD API for managing warehouses, items and stock.
// @description
// @description  CRUD API для управления складами, товарами и остатками.
// @host localhost:8080
// @BasePath /
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var cfg config.DBConfig
	if err := cfg.Load(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	var srvCfg config.ServerConfig
	srvCfg.Load()
	logger.Info("config loaded")

	db, err := database.NewDB(cfg.DSN())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connected to database")

	if err := repository.MigrationExec(db); err != nil {
		logger.Error(err.Error())
		db.Close()
		os.Exit(1)
	}
	logger.Info("migrations applied")

	warehouseRepository := repository.NewWarehouseRepository(db)
	warehouseService := service.NewWarehouseService(warehouseRepository)
	warehouseHandler := handler.NewWarehouseHandler(warehouseService, logger)

	itemRepository := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepository)
	itemHandler := handler.NewItemHandler(itemService, logger)

	stockRepository := repository.NewStockRepository(db)
	stockService := service.NewStockService(stockRepository)
	stockHandler := handler.NewStockHandler(stockService, logger)

	mux := http.NewServeMux()

	mux.Handle("GET /warehouses", warehouseHandler.GetAll())
	mux.Handle("GET /warehouses/{id}", warehouseHandler.GetWarehouseByID())
	mux.Handle("POST /warehouses", warehouseHandler.CreateWarehouse())
	mux.Handle("DELETE /warehouses/{id}", warehouseHandler.DeleteWarehouseByID())

	mux.Handle("GET /items", itemHandler.GetAll())
	mux.Handle("GET /items/{id}", itemHandler.GetItemByID())
	mux.Handle("POST /items", itemHandler.CreateItem())

	mux.Handle("GET /stock", stockHandler.GetAll())
	mux.Handle("GET /warehouses/{id}/stock", stockHandler.GetStockByWarehouseID())
	mux.Handle("GET /items/{id}/stock", stockHandler.GetStockByItemID())
	mux.Handle("POST /stock", stockHandler.CreateStock())
	mux.Handle("PATCH /stock", stockHandler.UpdateStock())
	mux.Handle("POST /stock/transfer", stockHandler.TransferStock())

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:              srvCfg.Addr(),
		Handler:           middleware.Recover(middleware.Log(mux, logger), logger),
		ReadTimeout:       srvCfg.ReadTimeout,
		ReadHeaderTimeout: srvCfg.ReadHeaderTimeout,
		WriteTimeout:      srvCfg.WriteTimeout,
		IdleTimeout:       srvCfg.IdleTimeout,
	}

	logger.Info("server started", "port", server.Addr)

	displayHost := srvCfg.Host
	if displayHost == "" {
		displayHost = "localhost"
	}
	logger.Info("swagger docs", "url", fmt.Sprintf("http://%s:%s/swagger/index.html", displayHost, srvCfg.Port))

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server.ListenAndServe:", "error", err.Error())
		}
	}()
	<-ctx.Done()

	logger.Info("Gracefully shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), srvCfg.ShutdownTimeout)
	defer cancel()
	server.Shutdown(shutdownCtx)
	db.Close()
}
