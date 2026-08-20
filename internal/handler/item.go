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

type ItemHandler struct {
	service *service.ItemService
	logger  *slog.Logger
}

func NewItemHandler(s *service.ItemService, l *slog.Logger) *ItemHandler {
	return &ItemHandler{service: s, logger: l}
}

// GetAll godoc
// @Summary      Get all items / Получить список товаров
// @Description  Returns a list of all items.
// @Description
// @Description  Возвращает список всех товаров.
// @Tags         items
// @Produce      json
// @Success      200  {array}   model.Item
// @Failure      500  {object}  model.ErrorResponse
// @Router       /items [get]
func (h *ItemHandler) GetAll() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		items, err := h.service.GetAll(r.Context())
		if err != nil {
			h.logger.Error("failed to get all items", "error", err)
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, &items)
	})
}

// GetItemByID godoc
// @Summary      Get item by ID / Получить товар по ID
// @Description  Returns an item by its identifier.
// @Description
// @Description  Возвращает товар по его идентификатору.
// @Tags         items
// @Produce      json
// @Param        id   path      int  true  "Item ID / ID товара"
// @Success      200  {object}  model.Item
// @Failure      400  {object}  model.ErrorResponse
// @Failure      404  {object}  model.ErrorResponse
// @Failure      500  {object}  model.ErrorResponse
// @Router       /items/{id} [get]
func (h *ItemHandler) GetItemByID() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "id is not valid")
			return
		}

		item, err := h.service.GetItemByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				WriteError(w, http.StatusNotFound, "not found")
				return
			}
			h.logger.Error("failed to get item by id", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		WriteJSON(w, http.StatusOK, &item)
	})
}

// CreateItem godoc
// @Summary      Create item / Создать товар
// @Description  Creates a new item.
// @Description
// @Description  Создаёт новый товар.
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        request  body      model.CreateItemRequest  true  "Item data / Данные товара"
// @Success      201      {object}  model.Item
// @Failure      400      {object}  model.ErrorResponse
// @Failure      500      {object}  model.ErrorResponse
// @Router       /items [post]
func (h *ItemHandler) CreateItem() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var data model.CreateItemRequest
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			WriteError(w, http.StatusBadRequest, "bad request")
			return
		}

		if data.Name == "" {
			WriteError(w, http.StatusBadRequest, "name required")
			return
		}

		item := model.Item{
			Name:     data.Name,
			Category: data.Category,
		}
		if err := h.service.CreateItem(r.Context(), &item); err != nil {
			h.logger.Error("failed to create item", "error", err.Error())
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		WriteJSON(w, http.StatusCreated, &item)
	})
}
