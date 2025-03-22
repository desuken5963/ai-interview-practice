package job_posting

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

type jobPostingRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository.JobPostingRepository {
	return &jobPostingRepository{db: db}
}

func (r *jobPostingRepository) CreateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error) {
	if err := r.db.WithContext(ctx).Create(jobPosting).Error; err != nil {
		return nil, err
	}
	return jobPosting, nil
}

func (r *jobPostingRepository) UpdateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error) {
	// 既存の求人情報を取得
	var existingJobPosting entity.JobPosting
	if err := r.db.First(&existingJobPosting, jobPosting.ID).Error; err != nil {
		return nil, err
	}

	// 更新対象のフィールドを設定
	updates := map[string]interface{}{
		"company_id":  jobPosting.CompanyID,
		"title":       jobPosting.Title,
		"description": jobPosting.Description,
		"updated_at":  time.Now(),
	}

	// トランザクションを開始
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 求人情報を更新
		if err := tx.Model(&entity.JobPosting{}).Where("id = ?", jobPosting.ID).Updates(updates).Error; err != nil {
			return err
		}

		// 既存のカスタムフィールドを削除
		if err := tx.Where("job_id = ?", jobPosting.ID).Delete(&entity.JobCustomField{}).Error; err != nil {
			return err
		}

		// 新しいカスタムフィールドを作成
		if len(jobPosting.CustomFields) > 0 {
			for i := range jobPosting.CustomFields {
				jobPosting.CustomFields[i].JobID = jobPosting.ID
				jobPosting.CustomFields[i].CreatedAt = time.Now()
				jobPosting.CustomFields[i].UpdatedAt = time.Now()
			}
			if err := tx.Create(&jobPosting.CustomFields).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 更新後の求人情報を取得して返す
	var updatedJobPosting entity.JobPosting
	if err := r.db.Preload("CustomFields").First(&updatedJobPosting, jobPosting.ID).Error; err != nil {
		return nil, err
	}

	return &updatedJobPosting, nil
}

func (r *jobPostingRepository) DeleteJobPosting(ctx context.Context, id int) error {
	// カスタムフィールドは外部キー制約で自動的に削除される
	return r.db.WithContext(ctx).Delete(&entity.JobPosting{}, id).Error
}
