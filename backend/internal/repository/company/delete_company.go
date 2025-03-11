package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/company"
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

	// 企業情報が存在するか確認
	var count int64
	if err := tx.Model(&entity.Company{}).Where("id = ?", id).Count(&count).Error; err != nil {
		tx.Rollback()
		return err
	}

	if count == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}

	// カスタムフィールドを削除
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
