package job

import (
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/job"
	"gorm.io/gorm"
)

// jobRepository は求人情報のリポジトリ実装です
type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository は求人情報のリポジトリを生成します
func NewJobRepository(db *gorm.DB) repository.JobRepository {
	return &jobRepository{
		db: db,
	}
}
