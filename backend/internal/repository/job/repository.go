package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// JobRepository は求人情報のリポジトリインターフェースです
type JobRepository interface {
	// FindAll は求人情報の一覧を取得します
	FindAll(ctx context.Context, page, limit int) ([]entity.JobPosting, int64, error)

	// FindByCompanyID は指定された企業IDの求人情報一覧を取得します
	FindByCompanyID(ctx context.Context, companyID int, page, limit int) ([]entity.JobPosting, int64, error)

	// FindByID は指定されたIDの求人情報を取得します
	FindByID(ctx context.Context, id int) (*entity.JobPosting, error)

	// Create は新しい求人情報を作成します
	Create(ctx context.Context, job *entity.JobPosting) error

	// Update は既存の求人情報を更新します
	Update(ctx context.Context, job *entity.JobPosting) error

	// Delete は指定されたIDの求人情報を削除します
	Delete(ctx context.Context, id int) error

	// FindCompanyWithJobs は指定された企業IDの企業情報と求人情報一覧を取得します
	FindCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error)
}

// jobRepository は求人情報のリポジトリ実装です
type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository は求人情報のリポジトリを生成します
func NewJobRepository(db *gorm.DB) repository.JobRepository {
	return &jobRepository{
		db: db,
	}
}

// FindAll は求人情報の一覧を取得します
func (r *jobRepository) FindAll(ctx context.Context, page, limit int) ([]entity.JobPosting, int64, error) {
	var jobs []entity.JobPosting
	var total int64

	// 総件数を取得
	if err := r.db.Model(&entity.JobPosting{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit

	// 求人情報を取得
	if err := r.db.Offset(offset).Limit(limit).Preload("CustomFields").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// FindByCompanyID は指定された企業IDの求人情報一覧を取得します
func (r *jobRepository) FindByCompanyID(ctx context.Context, companyID int, page, limit int) ([]entity.JobPosting, int64, error) {
	var jobs []entity.JobPosting
	var total int64

	// 総件数を取得
	if err := r.db.Model(&entity.JobPosting{}).Where("company_id = ?", companyID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit

	// 求人情報を取得
	if err := r.db.Where("company_id = ?", companyID).Offset(offset).Limit(limit).Preload("CustomFields").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// FindByID は指定されたIDの求人情報を取得します
func (r *jobRepository) FindByID(ctx context.Context, id int) (*entity.JobPosting, error) {
	var job entity.JobPosting

	if err := r.db.Preload("CustomFields").First(&job, id).Error; err != nil {
		return nil, err
	}

	return &job, nil
}

// Create は新しい求人情報を作成します
func (r *jobRepository) Create(ctx context.Context, job *entity.JobPosting) error {
	return r.db.Create(job).Error
}

// Update は既存の求人情報を更新します
func (r *jobRepository) Update(ctx context.Context, job *entity.JobPosting) error {
	// カスタムフィールドを更新
	if len(job.CustomFields) > 0 {
		// 既存のカスタムフィールドを削除
		if err := r.db.Where("job_id = ?", job.ID).Delete(&entity.JobCustomField{}).Error; err != nil {
			return err
		}

		// 新しいカスタムフィールドを作成
		for i := range job.CustomFields {
			job.CustomFields[i].JobID = job.ID
		}
	}

	// 求人情報を更新
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(job).Error
}

// Delete は指定されたIDの求人情報を削除します
func (r *jobRepository) Delete(ctx context.Context, id int) error {
	// カスタムフィールドを削除
	if err := r.db.Where("job_id = ?", id).Delete(&entity.JobCustomField{}).Error; err != nil {
		return err
	}

	// 求人情報を削除
	return r.db.Delete(&entity.JobPosting{}, id).Error
}

// FindCompanyWithJobs は指定された企業IDの企業情報と求人情報一覧を取得します
func (r *jobRepository) FindCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
	var company entity.Company
	var jobs []entity.JobPosting

	if err := r.db.First(&company, companyID).Error; err != nil {
		return nil, nil, err
	}

	if err := r.db.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		return nil, nil, err
	}

	return &company, jobs, nil
}
