package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

// CreateCompanyUsecase は企業情報を作成するためのインターフェースです
type CreateCompanyUsecase interface {
	Execute(ctx context.Context, company *entity.Company) error
}

type createCompanyUsecase struct {
	repo repository.CreateCompanyRepository
}

// NewCreateCompanyUsecase は新しいCreateCompanyUsecaseインスタンスを作成します
func NewCreateCompanyUsecase(repo repository.CreateCompanyRepository) CreateCompanyUsecase {
	return &createCompanyUsecase{repo: repo}
}

// Execute は企業情報を作成します
func (u *createCompanyUsecase) Execute(ctx context.Context, company *entity.Company) error {
	return u.repo.Create(ctx, company)
}

// CreateCompany は新しい企業情報を作成します
// 後方互換性のために残しています
func (u *companyUseCase) CreateCompany(ctx context.Context, company *entity.Company) error {
	usecase := NewCreateCompanyUsecase(u.createCompanyRepo)
	return usecase.Execute(ctx, company)
}
