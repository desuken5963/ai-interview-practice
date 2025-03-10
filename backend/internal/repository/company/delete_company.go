package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type deleteCompanyRepository struct {
	db *gorm.DB
}

// NewDeleteCompanyRepository は新しいDeleteCompanyRepositoryインスタンスを作成します
func NewDeleteCompanyRepository(db *gorm.DB) repository.DeleteCompanyRepository {
	return &deleteCompanyRepository{db: db}
}

// Delete は指定されたIDの企業情報を削除します
func (r *deleteCompanyRepository) Delete(ctx context.Context, id int) error {
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

	// 関連する求人情報のカスタムフィールドを削除
	if err := tx.Exec("DELETE jcf FROM job_custom_fields jcf JOIN job_postings jp ON jcf.job_id = jp.id WHERE jp.company_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 関連する求人情報を削除
	if err := tx.Where("company_id = ?", id).Delete(&entity.JobPosting{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 企業のカスタムフィールドを削除
	if err := tx.Where("company_id = ?", id).Delete(&entity.CompanyCustomField{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 企業情報を削除
	if err := tx.Delete(&entity.Company{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// トランザクションをコミット
	return tx.Commit().Error
}

// 後方互換性のための実装
func (r *companyRepository) Delete(ctx context.Context, id int) error {
	repo := NewDeleteCompanyRepository(r.db)
	return repo.Delete(ctx, id)
}
