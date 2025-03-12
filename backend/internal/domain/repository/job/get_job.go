package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetJobRepository は求人情報を取得するためのリポジトリインターフェースです
type GetJobRepository interface {
	FindByID(ctx context.Context, id int) (*entity.JobPosting, error)
}
