package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// JobPostingRepository は求人情報に関する全てのリポジトリ操作を定義するインターフェースです
type JobPostingRepository interface {
	// Create は新しい求人情報を作成します
	Create(ctx context.Context, jobPosting *entity.JobPosting) error
	// FindByID は指定されたIDの求人情報を取得します
	FindByID(ctx context.Context, id int) (*entity.JobPosting, error)
	// List は求人情報の一覧を取得します
	List(ctx context.Context, offset, limit int) ([]*entity.JobPosting, error)
	// Update は既存の求人情報を更新します
	Update(ctx context.Context, jobPosting *entity.JobPosting) error
	// Delete は指定されたIDの求人情報を削除します
	Delete(ctx context.Context, id int) error
	// ListByCompanyID は指定された企業IDの求人情報一覧を取得します
	ListByCompanyID(ctx context.Context, companyID, offset, limit int) ([]*entity.JobPosting, error)
	// CountByCompanyID は指定された企業IDの求人情報の総数を取得します
	CountByCompanyID(ctx context.Context, companyID int) (int, error)
	// FindWithJobs は企業情報と関連する求人情報を取得します
	FindWithJobs(ctx context.Context, companyID int) (*entity.Company, error)
}
