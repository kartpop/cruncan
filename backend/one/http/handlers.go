package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	onerequest "github.com/kartpop/cruncan/backend/one/database/one_request"
	"github.com/kartpop/cruncan/backend/pkg/id"
	"github.com/kartpop/cruncan/backend/pkg/model"
)

type Handler struct {
	repo      onerequest.Repository
	idService id.Service
	logger    *slog.Logger
}

func NewHandler(repo onerequest.Repository, idService id.Service) *Handler {
	return &Handler{
		repo:      repo,
		idService: idService,
		logger:    slog.Default(),
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

	err = h.repo.Create(ctx, &onerequest.OneRequest{
		ReqID:  h.idService.GenerateID(),
		UserID: req.UserID,
		Req:    body,
	})
	if err != nil {
		h.logger.ErrorContext(ctx, fmt.Sprintf("failed to save request to database, error: %v", err))
		http.Error(w, fmt.Sprintf("failed to parse OneRequest json, error: %v", err), http.StatusBadRequest)
		return
	}

	h.logger.Info(fmt.Sprintf("successfully handled request for user: %s", req.UserID))
}
