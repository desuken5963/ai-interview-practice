package company

import (
	"context"

	"gorm.io/gorm"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
)

type companyRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	var companies []entity.Company
	var total int64

	offset := (page - 1) * limit

	// 総件数を取得
	if err := r.db.Model(&entity.Company{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 企業情報を取得（関連するカスタムフィールドと求人情報も含む）
	if err := r.db.Preload("CustomFields").
		Preload("JobPostings").
		Preload("JobPostings.CustomFields").
		Offset(offset).
		Limit(limit).
		Find(&companies).Error; err != nil {
		return nil, err
	}

	return &entity.CompanyResponse{
		Companies: companies,
		Total:     int(total),
		Page:      page,
		Limit:     limit,
	}, nil
}

func (r *companyRepository) GetCompanyByID(ctx context.Context, id int) (*entity.Company, error) {
	var company entity.Company
	if err := r.db.Preload("CustomFields").
		Preload("JobPostings").
		Preload("JobPostings.CustomFields").
		First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) CreateCompany(ctx context.Context, company *entity.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) UpdateCompany(ctx context.Context, company *entity.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) DeleteCompany(ctx context.Context, id int) error {
	return r.db.Delete(&entity.Company{}, id).Error
}
