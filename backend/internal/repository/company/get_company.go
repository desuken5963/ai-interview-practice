package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type getCompanyRepository struct {
	db *gorm.DB
}

// NewGetCompanyRepository は新しいGetCompanyRepositoryインスタンスを作成します
func NewGetCompanyRepository(db *gorm.DB) repository.GetCompanyRepository {
	return &getCompanyRepository{db: db}
}

// FindByID は指定されたIDの企業情報を取得します
func (r *getCompanyRepository) FindByID(ctx context.Context, id int) (*entity.Company, error) {
	var company entity.Company
	err := r.db.First(&company, id).Error
	if err != nil {
		return nil, err
	}

	// カスタムフィールドを取得
	var customFields []entity.CompanyCustomField
	err = r.db.Where("company_id = ?", id).Find(&customFields).Error
	if err != nil {
		return nil, err
	}

	company.CustomFields = customFields

	return &company, nil
}
