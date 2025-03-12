package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

type getCompanyWithJobsRepository struct {
	db *gorm.DB
}

// NewGetCompanyWithJobsRepository は新しいGetCompanyWithJobsRepositoryインスタンスを作成します
func NewGetCompanyWithJobsRepository(db *gorm.DB) repository.GetCompanyWithJobsRepository {
	return &getCompanyWithJobsRepository{db: db}
}

// FindCompanyWithJobs は指定された企業IDの企業情報と求人情報一覧を取得します
func (r *getCompanyWithJobsRepository) FindCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
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
