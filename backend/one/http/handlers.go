package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/kartpop/cruncan/backend/pkg/model"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler() *Handler {
	return &Handler{
		logger: slog.Default(),
	}
}

// Post is a handler for POST /one
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req model.OneRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.logger.ErrorContext(ctx, fmt.Sprintf("failed to parse OneRequest json: %s, error: %v", body, err))
		http.Error(w, fmt.Sprintf("failed to parse OneRequest json, error: %v", err), http.StatusBadRequest)
		return
	}

	h.logger.Info("received request", slog.String("body", string(body)))
}