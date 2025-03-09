package company

import (
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// companyRepository は企業情報のリポジトリ実装です
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository は企業情報のリポジトリを生成します
func NewCompanyRepository(db *gorm.DB) repository.CompanyRepository {
	return &companyRepository{
		db: db,
	}
}
