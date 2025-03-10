package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetCompaniesRepository は企業情報の一覧を取得するためのリポジトリインターフェースです
type GetCompaniesRepository interface {
	FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error)
}
