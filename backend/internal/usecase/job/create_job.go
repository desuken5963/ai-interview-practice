package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// CreateJobUsecase は求人情報を作成するためのインターフェースです
type CreateJobUsecase interface {
	Execute(ctx context.Context, job *entity.JobPosting) error
}

type createJobUsecase struct {
	repo repository.CreateJobRepository
}

// NewCreateJobUsecase は新しいCreateJobUsecaseインスタンスを作成します
func NewCreateJobUsecase(repo repository.CreateJobRepository) CreateJobUsecase {
	return &createJobUsecase{repo: repo}
}

// Execute は求人情報を作成します
func (u *createJobUsecase) Execute(ctx context.Context, job *entity.JobPosting) error {
	return u.repo.Create(ctx, job)
}
