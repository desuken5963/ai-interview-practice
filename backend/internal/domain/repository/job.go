package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

type JobRepository interface {
	// 基本的なCRUD操作
	Create(ctx context.Context, job *entity.JobPosting) error
	FindByID(ctx context.Context, id uint) (*entity.JobPosting, error)
	Update(ctx context.Context, job *entity.JobPosting) error
	Delete(ctx context.Context, id uint) error

	// 一覧取得
	List(ctx context.Context, offset, limit int) ([]*entity.JobPosting, error)

	// 企業IDによる取得
	ListByCompanyID(ctx context.Context, companyID uint, offset, limit int) ([]*entity.JobPosting, error)
}
