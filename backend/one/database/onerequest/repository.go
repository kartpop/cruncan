package onerequest

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type OneRequest struct {
	ReqID     string          `gorm:"type:text;primary_key;" json:"req_id"`
	UserID    string          `gorm:"type:text" json:"user_id"`
	CreatedAt time.Time       `gorm:"type:timestamptz not null" json:"created_at"`
	UpdatedAt time.Time       `gorm:"type:timestamptz not null" json:"updated_at"`
	Req       json.RawMessage `gorm:"type:jsonb not null" json:"req"`
}

func (OneRequest) TableName() string {
	return "request"
}

type Repository interface {
	Create(ctx context.Context, req *OneRequest) error
	Get(ctx context.Context, reqId string) (*OneRequest, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) WithTracing() *TracingRepository {
	return NewTracingRepository(r)
}

func (r *RepositoryImpl) Create(ctx context.Context, req *OneRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *RepositoryImpl) Get(ctx context.Context, reqId string) (*OneRequest, error) {
	var oneReq OneRequest
	err := r.db.WithContext(ctx).Where("req_id = ?", reqId).First(&oneReq).Error
	return &oneReq, err
}
