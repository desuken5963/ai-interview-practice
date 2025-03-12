package job_posting

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

type JobPostingUsecase interface {
	Create(ctx context.Context, jobPosting *entity.JobPosting) error
	Get(ctx context.Context, id int) (*entity.JobPosting, error)
	List(ctx context.Context, offset, limit int) ([]*entity.JobPosting, error)
	Update(ctx context.Context, jobPosting *entity.JobPosting) error
	Delete(ctx context.Context, id int) error
	ListByCompanyID(ctx context.Context, companyID, offset, limit int) (*entity.JobResponse, error)
	GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, error)
}

type jobPostingUsecase struct {
	repo repository.JobPostingRepository
}

func NewJobPostingUsecase(repo repository.JobPostingRepository) JobPostingUsecase {
	return &jobPostingUsecase{repo: repo}
}

func (u *jobPostingUsecase) Create(ctx context.Context, jobPosting *entity.JobPosting) error {
	return u.repo.Create(ctx, jobPosting)
}

func (u *jobPostingUsecase) Get(ctx context.Context, id int) (*entity.JobPosting, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *jobPostingUsecase) List(ctx context.Context, offset, limit int) ([]*entity.JobPosting, error) {
	return u.repo.List(ctx, offset, limit)
}

func (u *jobPostingUsecase) Update(ctx context.Context, jobPosting *entity.JobPosting) error {
	return u.repo.Update(ctx, jobPosting)
}

func (u *jobPostingUsecase) Delete(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

func (u *jobPostingUsecase) ListByCompanyID(ctx context.Context, companyID, offset, limit int) (*entity.JobResponse, error) {
	jobsPtr, err := u.repo.ListByCompanyID(ctx, companyID, offset, limit)
	if err != nil {
		return nil, err
	}

	// []*JobPosting から []JobPosting に変換
	jobs := make([]entity.JobPosting, len(jobsPtr))
	for i, j := range jobsPtr {
		if j != nil {
			jobs[i] = *j
		}
	}

	total, err := u.repo.CountByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}

	page := 1
	if limit > 0 {
		page = (offset / limit) + 1
	}

	response := &entity.JobResponse{
		Jobs:  jobs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}

func (u *jobPostingUsecase) GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, error) {
	return u.repo.FindWithJobs(ctx, companyID)
}
