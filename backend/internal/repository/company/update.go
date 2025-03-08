package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// Update は既存の企業情報を更新します
func (r *companyRepository) Update(ctx context.Context, company *entity.Company) error {
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

	// 企業情報を更新
	if err := tx.Model(&entity.Company{}).Where("id = ?", company.ID).Updates(map[string]interface{}{
		"name":                 company.Name,
		"business_description": company.BusinessDescription,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 既存のカスタムフィールドを削除
	if err := tx.Where("company_id = ?", company.ID).Delete(&entity.CompanyCustomField{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 新しいカスタムフィールドを作成
	if len(company.CustomFields) > 0 {
		for i := range company.CustomFields {
			company.CustomFields[i].CompanyID = company.ID
		}
		if err := tx.Create(&company.CustomFields).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// トランザクションをコミット
	return tx.Commit().Error
}
