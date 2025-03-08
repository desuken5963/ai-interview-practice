package company

import (
	"context"
	"errors"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"gorm.io/gorm"
)

// FindByID は指定されたIDの企業情報を取得します
func (r *companyRepository) FindByID(ctx context.Context, id int) (*entity.Company, error) {
	var company entity.Company

	// 企業情報を取得
	if err := r.db.First(&company, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 見つからない場合はnilを返す
		}
		return nil, err
	}

	// カスタムフィールドを取得
	if err := r.db.Model(&company).Association("CustomFields").Find(&company.CustomFields); err != nil {
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
