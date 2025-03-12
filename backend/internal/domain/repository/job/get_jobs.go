package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetJobsRepository は求人情報一覧を取得するためのリポジトリインターフェースです
type GetJobsRepository interface {
	FindAll(ctx context.Context, page, limit int) ([]entity.JobPosting, int64, error)
}
