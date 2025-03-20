package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

type CompanyRepository interface {
	GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error)
	CreateCompany(ctx context.Context, company *entity.Company) error
	UpdateCompany(ctx context.Context, company *entity.Company) error
	DeleteCompany(ctx context.Context, id int) error
}
