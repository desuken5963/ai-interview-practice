package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// Delete は指定されたIDの求人情報を削除します
func (r *jobRepository) Delete(ctx context.Context, id int) error {
	// カスタムフィールドを削除
	if err := r.db.Where("job_id = ?", id).Delete(&entity.JobCustomField{}).Error; err != nil {
		return err
	}

	// 求人情報を削除
	return r.db.Delete(&entity.JobPosting{}, id).Error
}
