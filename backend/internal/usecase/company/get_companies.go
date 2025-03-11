package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	companyRepo "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/company"
)

// GetCompaniesUsecase は企業情報の一覧を取得するためのインターフェースです
type GetCompaniesUsecase interface {
	Execute(ctx context.Context, page, limit int) (*entity.CompanyResponse, error)
}

type getCompaniesUsecase struct {
	repo companyRepo.GetCompaniesRepository
}

// NewGetCompaniesUsecase は新しいGetCompaniesUsecaseインスタンスを作成します
func NewGetCompaniesUsecase(repo companyRepo.GetCompaniesRepository) GetCompaniesUsecase {
	return &getCompaniesUsecase{repo: repo}
}

// Execute は企業情報の一覧を取得します
func (u *getCompaniesUsecase) Execute(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	// ページとリミットのデフォルト値を設定
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// リポジトリから企業情報を取得
	companies, total, err := u.repo.FindAll(ctx, page, limit)
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

// GetCompanies は企業情報の一覧を取得します
// 後方互換性のために残しています
func (u *companyUseCase) GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	usecase := NewGetCompaniesUsecase(u.getCompaniesRepo)
	return usecase.Execute(ctx, page, limit)
}

// companyUseCase は企業情報に関するユースケースの実装です
type companyUseCase struct {
	createCompanyRepo companyRepo.CreateCompanyRepository
	updateCompanyRepo companyRepo.UpdateCompanyRepository
	deleteCompanyRepo companyRepo.DeleteCompanyRepository
	getCompanyRepo    companyRepo.GetCompanyRepository
	getCompaniesRepo  companyRepo.GetCompaniesRepository
}

// NewCompanyUseCase は企業ユースケースの新しいインスタンスを作成します
func NewCompanyUseCase(
	createCompanyRepo companyRepo.CreateCompanyRepository,
	updateCompanyRepo companyRepo.UpdateCompanyRepository,
	deleteCompanyRepo companyRepo.DeleteCompanyRepository,
	getCompanyRepo companyRepo.GetCompanyRepository,
	getCompaniesRepo companyRepo.GetCompaniesRepository,
) CompanyUseCase {
	return &companyUseCase{
		createCompanyRepo: createCompanyRepo,
		updateCompanyRepo: updateCompanyRepo,
		deleteCompanyRepo: deleteCompanyRepo,
		getCompanyRepo:    getCompanyRepo,
		getCompaniesRepo:  getCompaniesRepo,
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
