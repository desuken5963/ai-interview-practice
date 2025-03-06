package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// companyRepository は企業リポジトリの実装です
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository は企業リポジトリの新しいインスタンスを作成します
func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

// FindAll は企業情報の一覧を取得します
func (r *companyRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error) {
	var companies []entity.Company
	var total int64

	// ページネーションの設定
	offset := (page - 1) * limit

	// 企業の総数を取得
	if err := r.db.Model(&entity.Company{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 企業情報を取得
	if err := r.db.Offset(offset).Limit(limit).Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	// 各企業のカスタムフィールドを取得
	for i := range companies {
		if err := r.db.Where("company_id = ?", companies[i].ID).Find(&companies[i].CustomFields).Error; err != nil {
			return nil, 0, err
		}

		// 求人数を取得
		var jobCount int64
		if err := r.db.Model(&entity.JobPosting{}).Where("company_id = ?", companies[i].ID).Count(&jobCount).Error; err != nil {
			return nil, 0, err
		}
		companies[i].JobCount = int(jobCount)
	}

	return companies, total, nil
}

// FindByID は指定されたIDの企業情報を取得します
func (r *companyRepository) FindByID(ctx context.Context, id int) (*entity.Company, error) {
	var company entity.Company

	if err := r.db.First(&company, id).Error; err != nil {
		return nil, err
	}

	// カスタムフィールドを取得
	if err := r.db.Where("company_id = ?", company.ID).Find(&company.CustomFields).Error; err != nil {
		return nil, err
	}

	// 求人数を取得
	var jobCount int64
	if err := r.db.Model(&entity.JobPosting{}).Where("company_id = ?", company.ID).Count(&jobCount).Error; err != nil {
		return nil, err
	}
	company.JobCount = int(jobCount)

	return &company, nil
}

// Create は新しい企業情報を作成します
func (r *companyRepository) Create(ctx context.Context, company *entity.Company) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 企業情報を作成
		if err := tx.Create(company).Error; err != nil {
			return err
		}

		// カスタムフィールドを作成
		for i := range company.CustomFields {
			company.CustomFields[i].CompanyID = company.ID
			if err := tx.Create(&company.CustomFields[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Update は既存の企業情報を更新します
func (r *companyRepository) Update(ctx context.Context, company *entity.Company) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 企業情報を更新
		if err := tx.Model(company).Updates(map[string]interface{}{
			"name":                 company.Name,
			"business_description": company.BusinessDescription,
		}).Error; err != nil {
			return err
		}

		// 既存のカスタムフィールドを削除
		if err := tx.Where("company_id = ?", company.ID).Delete(&entity.CompanyCustomField{}).Error; err != nil {
			return err
		}

		// 新しいカスタムフィールドを作成
		for i := range company.CustomFields {
			company.CustomFields[i].CompanyID = company.ID
			if err := tx.Create(&company.CustomFields[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete は指定されたIDの企業情報を削除します
func (r *companyRepository) Delete(ctx context.Context, id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 関連するカスタムフィールドは外部キー制約によって自動的に削除されます
		if err := tx.Delete(&entity.Company{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}
