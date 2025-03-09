package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// FindByID は指定されたIDの求人情報を取得します
func (r *jobRepository) FindByID(ctx context.Context, id int) (*entity.JobPosting, error) {
	var job entity.JobPosting

	if err := r.db.Preload("CustomFields").First(&job, id).Error; err != nil {
		return nil, err
	}

	return &job, nil
}
