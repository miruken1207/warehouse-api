package handler

import (
	"log/slog"
	"net/http"
	"strconv"

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
