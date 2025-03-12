package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// GetJobUsecase は求人情報を取得するためのインターフェースです
type GetJobUsecase interface {
	Execute(ctx context.Context, id int) (*entity.JobPosting, error)
}

type getJobUsecase struct {
	repo repository.GetJobRepository
}

// NewGetJobUsecase は新しいGetJobUsecaseインスタンスを作成します
func NewGetJobUsecase(repo repository.GetJobRepository) GetJobUsecase {
	return &getJobUsecase{repo: repo}
}

// Execute は求人情報を取得します
func (u *getJobUsecase) Execute(ctx context.Context, id int) (*entity.JobPosting, error) {
	return u.repo.FindByID(ctx, id)
}
