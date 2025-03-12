package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

type updateJobRepository struct {
	db *gorm.DB
}

// NewUpdateJobRepository は新しいUpdateJobRepositoryインスタンスを作成します
func NewUpdateJobRepository(db *gorm.DB) repository.UpdateJobRepository {
	return &updateJobRepository{db: db}
}

// Update は既存の求人情報を更新します
func (r *updateJobRepository) Update(ctx context.Context, job *entity.JobPosting) error {
	return r.db.Model(job).Updates(map[string]interface{}{
		"title":       job.Title,
		"description": job.Description,
	}).Error
}
