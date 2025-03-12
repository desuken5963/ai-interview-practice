package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
)

// GetJobsByCompanyIDUsecase は企業IDに紐づく求人情報一覧を取得するためのインターフェースです
type GetJobsByCompanyIDUsecase interface {
	Execute(ctx context.Context, companyID int, page, limit int) (*entity.JobResponse, error)
}

type getJobsByCompanyIDUsecase struct {
	repo repository.GetJobsByCompanyRepository
}

// NewGetJobsByCompanyIDUsecase は新しいGetJobsByCompanyIDUsecaseインスタンスを作成します
func NewGetJobsByCompanyIDUsecase(repo repository.GetJobsByCompanyRepository) GetJobsByCompanyIDUsecase {
	return &getJobsByCompanyIDUsecase{repo: repo}
}

// Execute は企業IDに紐づく求人情報一覧を取得します
func (u *getJobsByCompanyIDUsecase) Execute(ctx context.Context, companyID int, page, limit int) (*entity.JobResponse, error) {
	// ページとリミットのデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// リポジトリから求人情報を取得
	jobs, total, err := u.repo.FindByCompanyID(ctx, companyID, page, limit)
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
