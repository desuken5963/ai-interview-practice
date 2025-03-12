package job_posting

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"gorm.io/gorm"
)

// JobPostingRepository は求人情報に関する全てのリポジトリ操作を定義するインターフェースです
type JobPostingRepository interface {
	Create(ctx context.Context, jobPosting *entity.JobPosting) error
	FindByID(ctx context.Context, id int) (*entity.JobPosting, error)
	List(ctx context.Context, offset, limit int) ([]*entity.JobPosting, error)
	Update(ctx context.Context, jobPosting *entity.JobPosting) error
	Delete(ctx context.Context, id int) error
	ListByCompanyID(ctx context.Context, companyID, offset, limit int) ([]*entity.JobPosting, error)
	CountByCompanyID(ctx context.Context, companyID int) (int, error)
	FindWithJobs(ctx context.Context, companyID int) (*entity.Company, error)
}

type jobPostingRepository struct {
	db *gorm.DB
}

func NewJobPostingRepository(db *gorm.DB) JobPostingRepository {
	return &jobPostingRepository{db: db}
}

func (r *jobPostingRepository) Create(ctx context.Context, jobPosting *entity.JobPosting) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// カスタムフィールドを一時保存
		customFields := jobPosting.CustomFields
		jobPosting.CustomFields = nil

		// 求人情報を保存
		if err := tx.Create(jobPosting).Error; err != nil {
			return err
		}

		// カスタムフィールドを保存
		for i := range customFields {
			customFields[i].ID = 0
			customFields[i].JobID = jobPosting.ID
			if err := tx.Create(&customFields[i]).Error; err != nil {
				return err
			}
		}

		// 保存したカスタムフィールドを設定
		jobPosting.CustomFields = customFields
		return nil
	})
}

func (r *jobPostingRepository) FindByID(ctx context.Context, id int) (*entity.JobPosting, error) {
	var jobPosting entity.JobPosting
	if err := r.db.WithContext(ctx).
		Preload("CustomFields").
		First(&jobPosting, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &jobPosting, nil
}

func (r *jobPostingRepository) List(ctx context.Context, offset, limit int) ([]*entity.JobPosting, error) {
	var jobs []*entity.JobPosting
	if err := r.db.WithContext(ctx).
		Preload("CustomFields").
		Offset(offset).
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobPostingRepository) Update(ctx context.Context, jobPosting *entity.JobPosting) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 既存のカスタムフィールドを削除
		if err := tx.Where("job_id = ?", jobPosting.ID).
			Delete(&entity.JobCustomField{}).Error; err != nil {
			return err
		}

		// カスタムフィールドを一時保存
		customFields := jobPosting.CustomFields
		jobPosting.CustomFields = nil

		// 求人情報を更新（必要なフィールドのみ）
		if err := tx.Model(&entity.JobPosting{}).
			Where("id = ?", jobPosting.ID).
			Updates(map[string]interface{}{
				"company_id":  jobPosting.CompanyID,
				"title":       jobPosting.Title,
				"description": jobPosting.Description,
			}).Error; err != nil {
			return err
		}

		// カスタムフィールドを保存
		for i := range customFields {
			customFields[i].ID = 0
			customFields[i].JobID = jobPosting.ID
			if err := tx.Create(&customFields[i]).Error; err != nil {
				return err
			}
		}

		// 保存したカスタムフィールドを設定
		jobPosting.CustomFields = customFields

		// 更新後のデータを取得
		if err := tx.First(jobPosting, jobPosting.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *jobPostingRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("job_id = ?", id).Delete(&entity.JobCustomField{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&entity.JobPosting{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *jobPostingRepository) ListByCompanyID(ctx context.Context, companyID, offset, limit int) ([]*entity.JobPosting, error) {
	var jobs []*entity.JobPosting
	if err := r.db.WithContext(ctx).
		Preload("CustomFields").
		Where("company_id = ?", companyID).
		Offset(offset).
		Limit(limit).
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobPostingRepository) CountByCompanyID(ctx context.Context, companyID int) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.JobPosting{}).
		Where("company_id = ?", companyID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *jobPostingRepository) FindWithJobs(ctx context.Context, companyID int) (*entity.Company, error) {
	var company entity.Company
	if err := r.db.WithContext(ctx).
		Preload("Jobs").
		Preload("Jobs.CustomFields").
		First(&company, companyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}
