package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

type getJobRepository struct {
	db *gorm.DB
}

// NewGetJobRepository は新しいGetJobRepositoryインスタンスを作成します
func NewGetJobRepository(db *gorm.DB) repository.GetJobRepository {
	return &getJobRepository{db: db}
}

// FindByID は指定されたIDの求人情報を取得します
func (r *getJobRepository) FindByID(ctx context.Context, id int) (*entity.JobPosting, error) {
	var job entity.JobPosting
	if err := r.db.Preload("CustomFields").First(&job, id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
