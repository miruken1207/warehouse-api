package handler

import (
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
			h.logger.Error("failed to get stock by warehouse_id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, stock)
	})
}

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
			h.logger.Error("failed to get stock by item_id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		WriteJSON(w, http.StatusOK, &stock)
	})
}

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
			h.logger.Error("StockHandler.TransferStock: %w", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		
		w.WriteHeader(http.StatusOK)
	})
}
