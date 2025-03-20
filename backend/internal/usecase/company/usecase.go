package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

type Usecase interface {
	GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error)
	GetCompanyByID(ctx context.Context, id int) (*entity.Company, error)
	CreateCompany(ctx context.Context, company *entity.Company) error
	UpdateCompany(ctx context.Context, company *entity.Company) error
	DeleteCompany(ctx context.Context, id int) error
}

type usecase struct {
	repo repository.CompanyRepository
}

func NewUsecase(repo repository.CompanyRepository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	return u.repo.GetCompanies(ctx, page, limit)
}

func (u *usecase) GetCompanyByID(ctx context.Context, id int) (*entity.Company, error) {
	return u.repo.GetCompanyByID(ctx, id)
}

func (u *usecase) CreateCompany(ctx context.Context, company *entity.Company) error {
	return u.repo.CreateCompany(ctx, company)
}

func (u *usecase) UpdateCompany(ctx context.Context, company *entity.Company) error {
	return u.repo.UpdateCompany(ctx, company)
}

func (u *usecase) DeleteCompany(ctx context.Context, id int) error {
	return u.repo.DeleteCompany(ctx, id)
}
