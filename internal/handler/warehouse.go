package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/miruken1207/warehouse-api/internal/service"
)

type WarehouseHandler struct {
	service *service.WarehouseService
	logger  *slog.Logger
}

func NewWarehouseHandler(s *service.WarehouseService, l *slog.Logger) *WarehouseHandler {
	return &WarehouseHandler{service: s, logger: l}
}

func (h *WarehouseHandler) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		warehouses, err := h.service.GetAll(r.Context())
		if err != nil {
			h.logger.Error("failed to get all warehouses", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, warehouses)
	})
}

func (h *WarehouseHandler) GetWarehouseByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "id is not valid")
			return
		}
		warehouse, err := h.service.GetWarehouseByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			}
			h.logger.Error("failed to get warehouse by id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, warehouse)
	})
}




