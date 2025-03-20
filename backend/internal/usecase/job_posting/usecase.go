package job_posting

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

type UseCase interface {
	CreateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error)
	UpdateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error)
	DeleteJobPosting(ctx context.Context, id int) error
}

type usecase struct {
	repo repository.JobPostingRepository
}

func NewUseCase(repo repository.JobPostingRepository) UseCase {
	return &usecase{repo: repo}
}

func (u *usecase) CreateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error) {
	return u.repo.CreateJobPosting(ctx, jobPosting)
}

func (u *usecase) UpdateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error) {
	return u.repo.UpdateJobPosting(ctx, jobPosting)
}

func (u *usecase) DeleteJobPosting(ctx context.Context, id int) error {
	return u.repo.DeleteJobPosting(ctx, id)
}
