package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

// CreateCompanyUsecase は企業情報作成のユースケースを定義します
type CreateCompanyUsecase interface {
	Execute(ctx context.Context, company *entity.Company) error
}

// GetCompanyUsecase は企業情報取得のユースケースを定義します
type GetCompanyUsecase interface {
	Execute(ctx context.Context, id int) (*entity.Company, error)
}

// GetCompaniesUsecase は企業情報一覧取得のユースケースを定義します
type GetCompaniesUsecase interface {
	Execute(ctx context.Context, page, limit int) (*entity.CompanyResponse, error)
}

// UpdateCompanyUsecase は企業情報更新のユースケースを定義します
type UpdateCompanyUsecase interface {
	Execute(ctx context.Context, company *entity.Company) error
}

// DeleteCompanyUsecase は企業情報削除のユースケースを定義します
type DeleteCompanyUsecase interface {
	Execute(ctx context.Context, id int) error
}

// createCompanyUsecase は企業情報作成のユースケース実装です
type createCompanyUsecase struct {
	repo repository.CompanyRepository
}

// NewCreateCompanyUsecase は新しいCreateCompanyUsecaseインスタンスを作成します
func NewCreateCompanyUsecase(repo repository.CompanyRepository) CreateCompanyUsecase {
	return &createCompanyUsecase{repo: repo}
}

func (u *createCompanyUsecase) Execute(ctx context.Context, company *entity.Company) error {
	return u.repo.Create(ctx, company)
}

// getCompanyUsecase は企業情報取得のユースケース実装です
type getCompanyUsecase struct {
	repo repository.CompanyRepository
}

// NewGetCompanyUsecase は新しいGetCompanyUsecaseインスタンスを作成します
func NewGetCompanyUsecase(repo repository.CompanyRepository) GetCompanyUsecase {
	return &getCompanyUsecase{repo: repo}
}

func (u *getCompanyUsecase) Execute(ctx context.Context, id int) (*entity.Company, error) {
	return u.repo.FindByID(ctx, id)
}

// getCompaniesUsecase は企業情報一覧取得のユースケース実装です
type getCompaniesUsecase struct {
	repo repository.CompanyRepository
}

// NewGetCompaniesUsecase は新しいGetCompaniesUsecaseインスタンスを作成します
func NewGetCompaniesUsecase(repo repository.CompanyRepository) GetCompaniesUsecase {
	return &getCompaniesUsecase{repo: repo}
}

func (u *getCompaniesUsecase) Execute(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	offset := (page - 1) * limit

	companiesPtr, err := u.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	// []*Company から []Company に変換
	companies := make([]entity.Company, len(companiesPtr))
	for i, c := range companiesPtr {
		if c != nil {
			companies[i] = *c
		}
	}

	total, err := u.repo.Count(ctx)
	if err != nil {
		return nil, err
	}

	response := &entity.CompanyResponse{
		Companies: companies,
		Total:     total,
		Page:      page,
		Limit:     limit,
	}

	return response, nil
}

// updateCompanyUsecase は企業情報更新のユースケース実装です
type updateCompanyUsecase struct {
	repo repository.CompanyRepository
}

// NewUpdateCompanyUsecase は新しいUpdateCompanyUsecaseインスタンスを作成します
func NewUpdateCompanyUsecase(repo repository.CompanyRepository) UpdateCompanyUsecase {
	return &updateCompanyUsecase{repo: repo}
}

func (u *updateCompanyUsecase) Execute(ctx context.Context, company *entity.Company) error {
	return u.repo.Update(ctx, company)
}

// deleteCompanyUsecase は企業情報削除のユースケース実装です
type deleteCompanyUsecase struct {
	repo repository.CompanyRepository
}

// NewDeleteCompanyUsecase は新しいDeleteCompanyUsecaseインスタンスを作成します
func NewDeleteCompanyUsecase(repo repository.CompanyRepository) DeleteCompanyUsecase {
	return &deleteCompanyUsecase{repo: repo}
}

func (u *deleteCompanyUsecase) Execute(ctx context.Context, id int) error {
	return u.repo.Delete(ctx, id)
}
