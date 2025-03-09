package job

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"gorm.io/gorm"
)

// Update は既存の求人情報を更新します
func (r *jobRepository) Update(ctx context.Context, job *entity.JobPosting) error {
	// カスタムフィールドを更新
	if len(job.CustomFields) > 0 {
		// 既存のカスタムフィールドを削除
		if err := r.db.Where("job_id = ?", job.ID).Delete(&entity.JobCustomField{}).Error; err != nil {
			return err
		}

		// 新しいカスタムフィールドを作成
		for i := range job.CustomFields {
			job.CustomFields[i].JobID = job.ID
		}
	}

	// 求人情報を更新
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(job).Error
}
