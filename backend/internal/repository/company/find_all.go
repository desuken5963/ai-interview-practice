package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// FindAll は企業情報の一覧を取得します
func (r *companyRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error) {
	var companies []entity.Company
	var total int64

	// オフセットを計算
	offset := (page - 1) * limit

	// 総件数を取得
	if err := r.db.Model(&entity.Company{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 企業情報を取得
	if err := r.db.Offset(offset).Limit(limit).Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	// カスタムフィールドを取得
	for i := range companies {
		if err := r.db.Model(&companies[i]).Association("CustomFields").Find(&companies[i].CustomFields); err != nil {
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
