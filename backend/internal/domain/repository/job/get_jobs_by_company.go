package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetJobsByCompanyRepository は企業IDに基づいて求人情報一覧を取得するためのリポジトリインターフェースです
type GetJobsByCompanyRepository interface {
	FindByCompanyID(ctx context.Context, companyID int, page, limit int) ([]entity.JobPosting, int64, error)
}
