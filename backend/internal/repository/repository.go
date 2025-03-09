package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// CompanyRepository は企業情報のリポジトリのインターフェースです
type CompanyRepository interface {
	FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error)
	FindByID(ctx context.Context, id int) (*entity.Company, error)
	Create(ctx context.Context, company *entity.Company) error
	Update(ctx context.Context, company *entity.Company) error
	Delete(ctx context.Context, id int) error
}

// JobRepository は求人情報のリポジトリのインターフェースです
type JobRepository interface {
	FindAll(ctx context.Context) ([]entity.JobPosting, int64, error)
	FindByCompanyID(ctx context.Context, companyID int, page, limit int) ([]entity.JobPosting, int64, error)
	FindByID(ctx context.Context, id int) (*entity.JobPosting, error)
	Create(ctx context.Context, job *entity.JobPosting) error
	Update(ctx context.Context, job *entity.JobPosting) error
	Delete(ctx context.Context, id int) error
	FindCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error)
}
