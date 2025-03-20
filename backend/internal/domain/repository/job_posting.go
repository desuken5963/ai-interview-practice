package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

type JobPostingRepository interface {
	CreateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error)
	UpdateJobPosting(ctx context.Context, jobPosting *entity.JobPosting) (*entity.JobPosting, error)
	DeleteJobPosting(ctx context.Context, id int) error
}
