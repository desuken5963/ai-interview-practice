package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// Create は新しい求人情報を作成します
func (r *jobRepository) Create(ctx context.Context, job *entity.JobPosting) error {
	return r.db.Create(job).Error
}
