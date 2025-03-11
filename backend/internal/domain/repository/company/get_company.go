package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetCompanyRepository は企業情報を取得するためのリポジトリインターフェースです
type GetCompanyRepository interface {
	FindByID(ctx context.Context, id int) (*entity.Company, error)
}
