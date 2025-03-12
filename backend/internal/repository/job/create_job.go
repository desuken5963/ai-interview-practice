package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

type createJobRepository struct {
	db *gorm.DB
}

// NewCreateJobRepository は新しいCreateJobRepositoryインスタンスを作成します
func NewCreateJobRepository(db *gorm.DB) repository.CreateJobRepository {
	return &createJobRepository{db: db}
}

// Create は新しい求人情報を作成します
func (r *createJobRepository) Create(ctx context.Context, job *entity.JobPosting) error {
	return r.db.Create(job).Error
}
