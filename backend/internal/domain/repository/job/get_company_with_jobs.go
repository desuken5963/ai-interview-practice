package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetCompanyWithJobsRepository は企業情報と求人情報一覧を取得するためのリポジトリインターフェースです
type GetCompanyWithJobsRepository interface {
	FindCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error)
}
