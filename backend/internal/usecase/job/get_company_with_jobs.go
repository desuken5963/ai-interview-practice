package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	companyRepo "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/company"
	jobRepo "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// GetCompanyWithJobsUsecase は企業情報と求人情報を取得するためのインターフェースです
type GetCompanyWithJobsUsecase interface {
	Execute(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error)
}

type getCompanyWithJobsUsecase struct {
	jobRepo        jobRepo.GetJobsByCompanyRepository
	getCompanyRepo companyRepo.GetCompanyRepository
}

// NewGetCompanyWithJobsUsecase は新しいGetCompanyWithJobsUsecaseインスタンスを作成します
func NewGetCompanyWithJobsUsecase(jobRepo jobRepo.GetJobsByCompanyRepository, getCompanyRepo companyRepo.GetCompanyRepository) GetCompanyWithJobsUsecase {
	return &getCompanyWithJobsUsecase{
		jobRepo:        jobRepo,
		getCompanyRepo: getCompanyRepo,
	}
}

// Execute は企業情報と求人情報を取得します
func (u *getCompanyWithJobsUsecase) Execute(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
	company, err := u.getCompanyRepo.FindByID(ctx, companyID)
	if err != nil {
		return nil, nil, err
	}

	jobs, _, err := u.jobRepo.FindByCompanyID(ctx, companyID, 1, 100)
	if err != nil {
		return nil, nil, err
	}

	return company, jobs, nil
}
