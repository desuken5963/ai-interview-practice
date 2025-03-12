package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// CreateJobRepository は求人情報を作成するためのリポジトリインターフェースです
type CreateJobRepository interface {
	Create(ctx context.Context, job *entity.JobPosting) error
}
