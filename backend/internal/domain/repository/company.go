package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// CompanyRepository は企業情報のリポジトリインターフェースです
type CompanyRepository interface {
	// GetCompanies は企業一覧を取得します
	GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error)

	// GetCompanyByID は指定されたIDの企業情報を取得します
	GetCompanyByID(ctx context.Context, id int) (*entity.Company, error)

	// CreateCompany は新規の企業情報を作成します
	CreateCompany(ctx context.Context, company *entity.Company) error

	// UpdateCompany は既存の企業情報を更新します
	UpdateCompany(ctx context.Context, company *entity.Company) error

	// DeleteCompany は指定されたIDの企業情報を削除します
	DeleteCompany(ctx context.Context, id int) error
}
