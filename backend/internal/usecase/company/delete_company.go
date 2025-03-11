package company

import (
	"context"

	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/company"
)

// DeleteCompanyUsecase は企業情報を削除するためのインターフェースです
type DeleteCompanyUsecase interface {
	Execute(ctx context.Context, id int) error
}

type deleteCompanyUsecase struct {
	repo repository.DeleteCompanyRepository
}

// NewDeleteCompanyUsecase は新しいDeleteCompanyUsecaseインスタンスを作成します
func NewDeleteCompanyUsecase(repo repository.DeleteCompanyRepository) DeleteCompanyUsecase {
	return &deleteCompanyUsecase{repo: repo}
}

// Execute は指定されたIDの企業情報を削除します
func (u *deleteCompanyUsecase) Execute(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}

// DeleteCompany は指定されたIDの企業情報を削除します
// 後方互換性のために残しています
func (u *companyUseCase) DeleteCompany(ctx context.Context, id int) error {
	usecase := NewDeleteCompanyUsecase(u.deleteCompanyRepo)
	return usecase.Execute(ctx, id)
}
