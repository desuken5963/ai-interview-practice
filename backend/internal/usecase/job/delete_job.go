package job

import (
	"context"

	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// DeleteJobUsecase は求人情報を削除するためのインターフェースです
type DeleteJobUsecase interface {
	Execute(ctx context.Context, id int) error
}

type deleteJobUsecase struct {
	repo repository.DeleteJobRepository
}

// NewDeleteJobUsecase は新しいDeleteJobUsecaseインスタンスを作成します
func NewDeleteJobUsecase(repo repository.DeleteJobRepository) DeleteJobUsecase {
	return &deleteJobUsecase{repo: repo}
}

// Execute は求人情報を削除します
func (u *deleteJobUsecase) Execute(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}
