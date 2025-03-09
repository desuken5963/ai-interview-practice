package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// FindByCompanyID は指定された企業IDの求人情報一覧を取得します
func (r *jobRepository) FindByCompanyID(ctx context.Context, companyID int, page, limit int) ([]entity.JobPosting, int64, error) {
	var jobs []entity.JobPosting
	var total int64

	// 総件数を取得
	if err := r.db.Model(&entity.JobPosting{}).Where("company_id = ?", companyID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit

	// 求人情報を取得
	if err := r.db.Where("company_id = ?", companyID).Offset(offset).Limit(limit).Preload("CustomFields").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}
