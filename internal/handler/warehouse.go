package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/miruken1207/warehouse-api/internal/model"
	"github.com/miruken1207/warehouse-api/internal/service"
)

type WarehouseHandler struct {
	service *service.WarehouseService
	logger  *slog.Logger
}

func NewWarehouseHandler(s *service.WarehouseService, l *slog.Logger) *WarehouseHandler {
	return &WarehouseHandler{service: s, logger: l}
}

// GetAll godoc
// @Summary      Get all warehouses / Получить список складов
// @Description  Returns a list of all warehouses.
// @Description
// @Description  Возвращает список всех складов.
// @Tags         warehouses
// @Produce      json
// @Success      200  {array}   model.Warehouse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /warehouses [get]
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

// GetWarehouseByID godoc
// @Summary      Get warehouse by ID / Получить склад по ID
// @Description  Returns a warehouse by its identifier.
// @Description
// @Description  Возвращает склад по его идентификатору.
// @Tags         warehouses
// @Produce      json
// @Param        id   path      int  true  "Warehouse ID / ID склада"
// @Success      200  {object}  model.Warehouse
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /warehouses/{id} [get]
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

// CreateWarehouse godoc
// @Summary      Create warehouse / Создать склад
// @Description  Creates a new warehouse.
// @Description
// @Description  Создаёт новый склад.
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        request  body      model.CreateWarehouseRequest  true  "Warehouse data / Данные склада"
// @Success      201      {object}  model.Warehouse
// @Failure      400      {object}  model.ErrorResponse
// @Failure      500      {object}  model.ErrorResponse
// @Router       /warehouses [post]
func (h *WarehouseHandler) CreateWarehouse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var data model.CreateWarehouseRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			WriteError(w, http.StatusBadRequest, "bad request")
			return
		}

		if data.Name == "" || data.Location == "" {
			WriteError(w, http.StatusBadRequest, "name and location are required")
			return
		}

		warehouse := model.Warehouse{
			Name:     data.Name,
			Location: data.Location,
		}
		if err := h.service.CreateWarehouse(r.Context(), &warehouse); err != nil {
			h.logger.Error("failed to create warehouse", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		WriteJSON(w, http.StatusCreated, &warehouse)
	})
}

// DeleteWarehouseByID godoc
// @Summary      Delete warehouse / Удалить склад
// @Description  Deletes a warehouse by its identifier.
// @Description
// @Description  Удаляет склад по его идентификатору.
// @Tags         warehouses
// @Param        id   path  int  true  "Warehouse ID / ID склада"
// @Success      204  "No Content"
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /warehouses/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouseByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "id is not valid")
			return
		}

		if err := h.service.DeleteWarehouseByID(r.Context(), id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			}
			h.logger.Error("failed to delete warehouse by id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
