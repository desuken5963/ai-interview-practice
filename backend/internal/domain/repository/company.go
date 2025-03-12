package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

type CompanyRepository interface {
	// 基本的なCRUD操作
	Create(ctx context.Context, company *entity.Company) error
	FindByID(ctx context.Context, id uint) (*entity.Company, error)
	Update(ctx context.Context, company *entity.Company) error
	Delete(ctx context.Context, id uint) error

	// 一覧取得
	List(ctx context.Context, offset, limit int) ([]*entity.Company, error)

	// 関連データを含む取得
	FindWithJobs(ctx context.Context, id uint) (*entity.Company, error)
}
