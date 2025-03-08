package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

// GetCompanies は企業情報の一覧を取得します
func (u *companyUseCase) GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	// ページとリミットのデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// リポジトリから企業情報を取得
	companies, total, err := u.companyRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	// レスポンスを作成
	response := &entity.CompanyResponse{
		Companies: companies,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}

	return response, nil
}

// companyUseCase は企業情報に関するユースケースの実装です
type companyUseCase struct {
	companyRepo repository.CompanyRepository
}

// NewCompanyUseCase は企業ユースケースの新しいインスタンスを作成します
func NewCompanyUseCase(companyRepo repository.CompanyRepository) CompanyUseCase {
	return &companyUseCase{
		companyRepo: companyRepo,
	}
}

// CompanyUseCase は企業情報に関するユースケースを定義するインターフェースです
type CompanyUseCase interface {
	// GetCompanies は企業情報の一覧を取得します
	GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error)

	// GetCompanyByID は指定されたIDの企業情報を取得します
	GetCompanyByID(ctx context.Context, id int) (*entity.Company, error)

	// CreateCompany は新しい企業情報を作成します
	CreateCompany(ctx context.Context, company *entity.Company) error

	// UpdateCompany は既存の企業情報を更新します
	UpdateCompany(ctx context.Context, company *entity.Company) error

	// DeleteCompany は指定されたIDの企業情報を削除します
	DeleteCompany(ctx context.Context, id int) error
}
