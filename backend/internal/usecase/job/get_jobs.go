package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// GetJobsUsecase は求人情報一覧を取得するためのインターフェースです
type GetJobsUsecase interface {
	Execute(ctx context.Context, page, limit int) (*entity.JobResponse, error)
}

type getJobsUsecase struct {
	repo repository.GetJobsRepository
}

// NewGetJobsUsecase は新しいGetJobsUsecaseインスタンスを作成します
func NewGetJobsUsecase(repo repository.GetJobsRepository) GetJobsUsecase {
	return &getJobsUsecase{repo: repo}
}

// Execute は求人情報一覧を取得します
func (u *getJobsUsecase) Execute(ctx context.Context, page, limit int) (*entity.JobResponse, error) {
	// ページとリミットのデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// リポジトリから求人情報を取得
	jobs, total, err := u.repo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	// レスポンスを作成
	response := &entity.JobResponse{
		Jobs:  jobs,
		Total: total,
		Page:  page,
		Limit: limit,
	}

	return response, nil
}
