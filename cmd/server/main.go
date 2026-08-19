package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/miruken1207/warehouse-api/internal/config"
	"github.com/miruken1207/warehouse-api/internal/handler"
	"github.com/miruken1207/warehouse-api/internal/repository"
	"github.com/miruken1207/warehouse-api/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var cfg config.DBConfig
	if err := cfg.Load(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("config loaded")

	db, err := repository.NewDB(cfg.DSN())
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

	mux.Handle("GET /warehouses/{id}/stock", stockHandler.GetStockByWarehouseID())
	mux.Handle("GET /items/{id}/stock", stockHandler.GetStockByItemID())

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info("server started", "port", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server.ListenAndServe:", "error", err.Error())
		}
	}()
	<-ctx.Done()

	logger.Info("Gracefully shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)
	db.Close()
}
