package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// UpdateJobUsecase は求人情報を更新するためのインターフェースです
type UpdateJobUsecase interface {
	Execute(ctx context.Context, job *entity.JobPosting) error
}

type updateJobUsecase struct {
	repo repository.UpdateJobRepository
}

// NewUpdateJobUsecase は新しいUpdateJobUsecaseインスタンスを作成します
func NewUpdateJobUsecase(repo repository.UpdateJobRepository) UpdateJobUsecase {
	return &updateJobUsecase{repo: repo}
}

// Execute は求人情報を更新します
func (u *updateJobUsecase) Execute(ctx context.Context, job *entity.JobPosting) error {
	return u.repo.Update(ctx, job)
}
