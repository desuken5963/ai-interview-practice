package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/company"
)

// GetCompanyUsecase は企業情報を取得するためのインターフェースです
type GetCompanyUsecase interface {
	Execute(ctx context.Context, id int) (*entity.Company, error)
}

type getCompanyUsecase struct {
	repo repository.GetCompanyRepository
}

// NewGetCompanyUsecase は新しいGetCompanyUsecaseインスタンスを作成します
func NewGetCompanyUsecase(repo repository.GetCompanyRepository) GetCompanyUsecase {
	return &getCompanyUsecase{repo: repo}
}

// Execute は指定されたIDの企業情報を取得します
func (u *getCompanyUsecase) Execute(ctx context.Context, id int) (*entity.Company, error) {
	return u.repo.FindByID(ctx, id)
}

// GetCompanyByID は指定されたIDの企業情報を取得します
// 後方互換性のために残しています
func (u *companyUseCase) GetCompanyByID(ctx context.Context, id int) (*entity.Company, error) {
	usecase := NewGetCompanyUsecase(u.getCompanyRepo)
	return usecase.Execute(ctx, id)
}
