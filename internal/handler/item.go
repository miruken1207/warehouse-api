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
