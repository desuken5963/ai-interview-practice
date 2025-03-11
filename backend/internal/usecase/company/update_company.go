package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

// UpdateCompanyUsecase は企業情報を更新するためのインターフェースです
type UpdateCompanyUsecase interface {
	Execute(ctx context.Context, company *entity.Company) error
}

type updateCompanyUsecase struct {
	repo repository.UpdateCompanyRepository
}

// NewUpdateCompanyUsecase は新しいUpdateCompanyUsecaseインスタンスを作成します
func NewUpdateCompanyUsecase(repo repository.UpdateCompanyRepository) UpdateCompanyUsecase {
	return &updateCompanyUsecase{repo: repo}
}

// Execute は企業情報を更新します
func (u *updateCompanyUsecase) Execute(ctx context.Context, company *entity.Company) error {
	return u.repo.Update(ctx, company)
}

// UpdateCompany は既存の企業情報を更新します
// 後方互換性のために残しています
func (u *companyUseCase) UpdateCompany(ctx context.Context, company *entity.Company) error {
	usecase := NewUpdateCompanyUsecase(u.companyRepo)
	return usecase.Execute(ctx, company)
}
