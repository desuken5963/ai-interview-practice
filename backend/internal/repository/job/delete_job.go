package job

import (
	"context"

	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

type deleteJobRepository struct {
	db *gorm.DB
}

// NewDeleteJobRepository は新しいDeleteJobRepositoryインスタンスを作成します
func NewDeleteJobRepository(db *gorm.DB) repository.DeleteJobRepository {
	return &deleteJobRepository{db: db}
}

// Delete は指定されたIDの求人情報を削除します
func (r *deleteJobRepository) Delete(ctx context.Context, id int) error {
	return r.db.Delete("job_postings", id).Error
}
