package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// UpdateJobRepository は求人情報を更新するためのリポジトリインターフェースです
type UpdateJobRepository interface {
	Update(ctx context.Context, job *entity.JobPosting) error
}
