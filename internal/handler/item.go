package handler

import (
	"log/slog"
	"net/http"

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
			h.logger.Error("failed to get all items: %w", err)
			WriteError(w, http.StatusInternalServerError, "internal server error")
			return
		}
		WriteJSON(w, http.StatusOK, &items)
	})
}
