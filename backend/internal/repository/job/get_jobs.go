package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

type getJobsRepository struct {
	db *gorm.DB
}

// NewGetJobsRepository は新しいGetJobsRepositoryインスタンスを作成します
func NewGetJobsRepository(db *gorm.DB) repository.GetJobsRepository {
	return &getJobsRepository{db: db}
}

// FindAll は求人情報の一覧を取得します
func (r *getJobsRepository) FindAll(ctx context.Context, page, limit int) ([]entity.JobPosting, int64, error) {
	var jobs []entity.JobPosting
	var total int64

	offset := (page - 1) * limit

	// 総件数を取得
	if err := r.db.Model(&entity.JobPosting{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 求人情報を取得
	if err := r.db.Offset(offset).Limit(limit).Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}
