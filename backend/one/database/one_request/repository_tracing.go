package onerequest

import (
	"context"
	"fmt"

	"github.com/kartpop/cruncan/backend/pkg/otel"
	otelContext "github.com/kartpop/cruncan/backend/pkg/otel/context"
)

type TracingRepository struct {
	repo Repository
}

func NewTracingRepository(repo Repository) *TracingRepository {
	return &TracingRepository{
		repo: repo,
	}
}

func (r *TracingRepository) Create(ctx context.Context, req *OneRequest) error {
	tracer, _ := otelContext.Tracer(ctx)
	ctx, span := tracer.Start(ctx, "oneRequest.Create")
	defer span.End()

	err := r.repo.Create(ctx, req)
	if err != nil {
		otel.SetSpanErrorWithMessage(span, err, fmt.Sprintf("failed to create one request: %v", err))
	} else {
		otel.SetSpanOk(span)
	}

	return err
}
