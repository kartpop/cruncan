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
}

type RepositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *RepositoryImpl {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) Create(ctx context.Context, req *OneRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}
