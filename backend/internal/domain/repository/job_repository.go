package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// JobRepository は求人情報の永続化を担当するインターフェースです
type JobRepository interface {
	// FindAll は求人情報の一覧を取得します
	FindAll(ctx context.Context, page, limit int) ([]entity.JobPosting, int64, error)

	// FindByCompanyID は指定された企業IDの求人情報一覧を取得します
	FindByCompanyID(ctx context.Context, companyID int, page, limit int) ([]entity.JobPosting, int64, error)

	// FindByID は指定されたIDの求人情報を取得します
	FindByID(ctx context.Context, id int) (*entity.JobPosting, error)

	// Create は新しい求人情報を作成します
	Create(ctx context.Context, job *entity.JobPosting) error

	// Update は既存の求人情報を更新します
	Update(ctx context.Context, job *entity.JobPosting) error

	// Delete は指定されたIDの求人情報を削除します
	Delete(ctx context.Context, id int) error
}
