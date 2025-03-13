package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository は新しいCompanyRepositoryインスタンスを作成します
func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepository{db: db}
}

// Create は新しい企業情報を作成します
func (r *companyRepository) Create(ctx context.Context, company *entity.Company) error {
	// トランザクションを開始
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ロールバック用のdefer
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 企業情報を作成
	if err := tx.Create(company).Error; err != nil {
		tx.Rollback()
		return err
	}

	// カスタムフィールドを作成
	if len(company.CustomFields) > 0 {
		for i := range company.CustomFields {
			company.CustomFields[i].CompanyID = company.ID
			company.CustomFields[i].ID = 0 // IDを明示的にゼロ値に設定して自動生成させる
		}
		if err := tx.Create(&company.CustomFields).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// トランザクションをコミット
	return tx.Commit().Error
}

// FindByID は指定されたIDの企業情報を取得します
func (r *companyRepository) FindByID(ctx context.Context, id int) (*entity.Company, error) {
	var company entity.Company

	// 企業情報と求人数を一緒に取得
	err := r.db.WithContext(ctx).
		Select("companies.*, COALESCE(COUNT(DISTINCT job_postings.id), 0) as job_count").
		Joins("LEFT JOIN job_postings ON companies.id = job_postings.company_id").
		Group("companies.id, companies.name, companies.business_description, companies.created_at, companies.updated_at").
		First(&company, id).Error
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

// Update は既存の企業情報を更新します
func (r *companyRepository) Update(ctx context.Context, company *entity.Company) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 既存の企業情報を取得
		var existing entity.Company
		if err := tx.First(&existing, company.ID).Error; err != nil {
			return err
		}

		// 企業情報を更新
		if err := tx.Model(&existing).Updates(map[string]interface{}{
			"name":                 company.Name,
			"business_description": company.BusinessDescription,
		}).Error; err != nil {
			return err
		}

		// カスタムフィールドを更新
		if err := tx.Where("company_id = ?", company.ID).Delete(&entity.CompanyCustomField{}).Error; err != nil {
			return err
		}

		if len(company.CustomFields) > 0 {
			for i := range company.CustomFields {
				company.CustomFields[i].CompanyID = company.ID
				company.CustomFields[i].ID = 0 // 新規作成のためIDをリセット
			}
			if err := tx.Create(&company.CustomFields).Error; err != nil {
				return err
			}
		}

		// 更新後の企業情報を取得（求人数も含める）
		if err := tx.
			Select("companies.*, COALESCE(COUNT(DISTINCT job_postings.id), 0) as job_count").
			Joins("LEFT JOIN job_postings ON companies.id = job_postings.company_id").
			Group("companies.id, companies.name, companies.business_description, companies.created_at, companies.updated_at").
			Preload("CustomFields").
			First(company, company.ID).Error; err != nil {
			return err
		}

		return nil
	})
}

// Delete は企業情報を削除します
func (r *companyRepository) Delete(ctx context.Context, id int) error {
	return r.db.Delete(&entity.Company{}, id).Error
}

// List は企業情報の一覧を取得します
func (r *companyRepository) List(ctx context.Context, offset, limit int) ([]*entity.Company, error) {
	var companies []*entity.Company

	// 企業情報と求人数を一緒に取得
	if err := r.db.WithContext(ctx).
		Select("companies.*, COALESCE(COUNT(DISTINCT job_postings.id), 0) as job_count").
		Joins("LEFT JOIN job_postings ON companies.id = job_postings.company_id").
		Group("companies.id, companies.name, companies.business_description, companies.created_at, companies.updated_at").
		Order("companies.id ASC").
		Offset(offset).
		Limit(limit).
		Find(&companies).Error; err != nil {
		return nil, err
	}

	// 各企業のカスタムフィールドを取得
	for i := range companies {
		if err := r.db.Where("company_id = ?", companies[i].ID).Find(&companies[i].CustomFields).Error; err != nil {
			return nil, err
		}
	}

	return companies, nil
}

// Count は企業情報の総数を取得します
func (r *companyRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Company{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindWithJobs は指定されたIDの企業情報と求人情報を取得します
func (r *companyRepository) FindWithJobs(ctx context.Context, id int) (*entity.Company, error) {
	var company entity.Company
	if err := r.db.WithContext(ctx).
		Select("companies.*, COALESCE(COUNT(DISTINCT job_postings.id), 0) as job_count").
		Joins("LEFT JOIN job_postings ON companies.id = job_postings.company_id").
		Group("companies.id, companies.name, companies.business_description, companies.created_at, companies.updated_at").
		Preload("Jobs").
		Preload("Jobs.CustomFields").
		First(&company, id).Error; err != nil {
		return nil, err
	}
	return &company, nil
}
