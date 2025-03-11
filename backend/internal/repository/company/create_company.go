package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type createCompanyRepository struct {
	db *gorm.DB
}

// NewCreateCompanyRepository は新しいCreateCompanyRepositoryインスタンスを作成します
func NewCreateCompanyRepository(db *gorm.DB) repository.CreateCompanyRepository {
	return &createCompanyRepository{db: db}
}

// Create は新しい企業情報を作成します
func (r *createCompanyRepository) Create(ctx context.Context, company *entity.Company) error {
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
