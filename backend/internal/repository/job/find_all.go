package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// FindAll は求人情報の一覧を取得します
func (r *jobRepository) FindAll(ctx context.Context, page, limit int) ([]entity.JobPosting, int64, error) {
	var jobs []entity.JobPosting
	var total int64

	// 総件数を取得
	if err := r.db.Model(&entity.JobPosting{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit

	// 求人情報を取得
	if err := r.db.Offset(offset).Limit(limit).Preload("CustomFields").Find(&jobs).Error; err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}
