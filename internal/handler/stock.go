package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/miruken1207/warehouse-api/internal/model"
	"github.com/miruken1207/warehouse-api/internal/repository"
	"github.com/miruken1207/warehouse-api/internal/service"
)

type StockHandler struct {
	service *service.StockService
	logger  *slog.Logger
}

func NewStockHandler(s *service.StockService, l *slog.Logger) *StockHandler {
	return &StockHandler{service: s, logger: l}
}

// GetAll godoc
// @Summary      Get all stock records / Получить все остатки
// @Description  Returns a list of all stock records across warehouses.
// @Description
// @Description  Возвращает список всех остатков товаров по складам.
// @Tags         stock
// @Produce      json
// @Success      200  {array}   model.Stock
// @Failure      500  {object}  model.ErrorResponse
// @Router       /stock [get]
func (h *StockHandler) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stock, err := h.service.GetAll(r.Context())
		if err != nil {
			h.logger.Error("failed to get all stock", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, stock)
	})
}

// GetStockByWarehouseID godoc
// @Summary      Get stock by warehouse / Получить остатки по складу
// @Description  Returns stock records for a warehouse by its identifier.
// @Description
// @Description  Возвращает остатки товаров на складе по его идентификатору.
// @Tags         stock
// @Produce      json
// @Param        id   path      int  true  "Warehouse ID / ID склада"
// @Success      200  {array}   model.Stock
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /warehouses/{id}/stock [get]
func (h *StockHandler) GetStockByWarehouseID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "id is not valid")
			return
		}
		stock, err := h.service.GetStockByWarehouseID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			}
			h.logger.Error("failed to get stock by warehouse_id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, stock)
	})
}

// GetStockByItemID godoc
// @Summary      Get stock by item / Получить остатки по товару
// @Description  Returns stock records for an item across all warehouses by its identifier.
// @Description
// @Description  Возвращает остатки товара по всем складам по идентификатору товара.
// @Tags         stock
// @Produce      json
// @Param        id   path      int  true  "Item ID / ID товара"
// @Success      200  {array}   model.Stock
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /items/{id}/stock [get]
func (h *StockHandler) GetStockByItemID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "id is not valid")
			return
		}

		stock, err := h.service.GetStockByItemID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			}
			h.logger.Error("failed to get stock by item_id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		WriteJSON(w, http.StatusOK, &stock)
	})
}

// CreateStock godoc
// @Summary      Create stock record / Создать остаток товара
// @Description  Creates a stock record for an item in a warehouse.
// @Description
// @Description  Создаёт запись об остатке товара на складе.
// @Tags         stock
// @Accept       json
// @Produce      json
// @Param        request  body      model.CreateStockRequest  true  "Stock data / Данные остатка"
// @Success      201      {object}  model.Stock
// @Failure      400      {object}  model.ErrorResponse
// @Failure      500      {object}  model.ErrorResponse
// @Router       /stock [post]
func (h *StockHandler) CreateStock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var data model.CreateStockRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			WriteError(w, http.StatusBadRequest, "bad request")
			return
		}

		if data.WarehouseID <= 0 || data.ItemID <= 0 {
			WriteError(w, http.StatusBadRequest, "warehouse_id and item_id are required")
			return
		}

		if data.Quantity < 0 {
			WriteError(w, http.StatusBadRequest, "quantity must not be negative")
			return
		}

		stock := model.Stock{
			WarehouseID: data.WarehouseID,
			ItemID:      data.ItemID,
			Quantity:    data.Quantity,
		}
		if err := h.service.CreateStock(r.Context(), &stock); err != nil {
			h.logger.Error("failed to create stock", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		WriteJSON(w, http.StatusCreated, &stock)
	})
}

// UpdateStock godoc
// @Summary      Adjust stock quantity / Изменить остаток товара
// @Description  Adjusts item quantity in a warehouse by a delta (can be negative).
// @Description
// @Description  Изменяет количество товара на складе на величину delta (может быть отрицательной).
// @Tags         stock
// @Accept       json
// @Param        request  body  model.UpdateStockRequest  true  "Stock adjustment data / Данные изменения остатка"
// @Success      200  "OK"
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse  "stock record not found / остаток не найден"
// @Failure      422  {object}  model.ErrorResponse  "insufficient stock quantity / недостаточно товара на складе"
// @Failure      500  {object}  model.ErrorResponse
// @Router       /stock [patch]
func (h *StockHandler) UpdateStock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var updateReq model.UpdateStockRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			WriteError(w, http.StatusBadRequest, "bad request")
			return
		}

		if err := h.service.UpdateStock(r.Context(), &updateReq); err != nil {
			switch err {
			case repository.ErrStockNotFound:
				WriteError(w, http.StatusNotFound, "not found")
				return
			case repository.ErrInsufficientStock:
				WriteError(w, http.StatusUnprocessableEntity, "unprocessable entity")
				return
			default:
				h.logger.Error("failed to update stock", "err", err.Error())
				WriteError(w, http.StatusInternalServerError, "internal server error")
			}
		}

		w.WriteHeader(http.StatusOK)
	})
}

// TransferStock godoc
// @Summary      Transfer stock between warehouses / Переместить товар между складами
// @Description  Moves a given quantity of an item from one warehouse to another.
// @Description
// @Description  Переносит указанное количество товара с одного склада на другой.
// @Tags         stock
// @Accept       json
// @Param        request  body  model.TransferStockRequest  true  "Transfer data / Данные перемещения"
// @Success      200  "OK"
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse  "stock record not found / остаток не найден"
// @Failure      422  {object}  model.ErrorResponse  "insufficient stock quantity / недостаточно товара на складе"
// @Failure      500  {object}  model.ErrorResponse
// @Router       /stock/transfer [post]
func (h *StockHandler) TransferStock() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var transferReq model.TransferStockRequest
		if err := json.NewDecoder(r.Body).Decode(&transferReq); err != nil {
			WriteError(w, http.StatusBadRequest, "bad request")
			return
		}

		if err := h.service.TransferStock(r.Context(), &transferReq); err != nil {
			if errors.Is(err, repository.ErrStockNotFound) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			}
			if errors.Is(err, repository.ErrInsufficientStock) {
				WriteError(w, http.StatusUnprocessableEntity, "unprocessable entity")
				return
			}
			h.logger.Error("StockHandler.TransferStock", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
