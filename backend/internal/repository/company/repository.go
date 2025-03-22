package company

import (
	"context"
	"time"

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

func (r *companyRepository) CreateCompany(ctx context.Context, company *entity.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) UpdateCompany(ctx context.Context, company *entity.Company) error {
	// 既存の企業情報を取得
	var existingCompany entity.Company
	if err := r.db.First(&existingCompany, company.ID).Error; err != nil {
		return err
	}

	// 更新対象のフィールドを設定
	updates := map[string]interface{}{
		"name":                 company.Name,
		"business_description": company.BusinessDescription,
		"updated_at":           time.Now(),
	}

	// トランザクションを開始
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 企業情報を更新
		if err := tx.Model(&entity.Company{}).Where("id = ?", company.ID).Updates(updates).Error; err != nil {
			return err
		}

		// 既存のカスタムフィールドを削除
		if err := tx.Where("company_id = ?", company.ID).Delete(&entity.CompanyCustomField{}).Error; err != nil {
			return err
		}

		// 新しいカスタムフィールドを作成
		if len(company.CustomFields) > 0 {
			for i := range company.CustomFields {
				company.CustomFields[i].CompanyID = company.ID
				company.CustomFields[i].CreatedAt = time.Now()
				company.CustomFields[i].UpdatedAt = time.Now()
			}
			if err := tx.Create(&company.CustomFields).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *companyRepository) DeleteCompany(ctx context.Context, id int) error {
	return r.db.Delete(&entity.Company{}, id).Error
}
