package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	repository "github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository/company"
	"gorm.io/gorm"
)

type getCompaniesRepository struct {
	db *gorm.DB
}

// NewGetCompaniesRepository は新しいGetCompaniesRepositoryインスタンスを作成します
func NewGetCompaniesRepository(db *gorm.DB) repository.GetCompaniesRepository {
	return &getCompaniesRepository{db: db}
}

// FindAll は企業情報の一覧を取得します
func (r *getCompaniesRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error) {
	var companies []entity.Company
	var total int64

	// 総数を取得
	if err := r.db.Model(&entity.Company{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit
	if err := r.db.Offset(offset).Limit(limit).Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	// 各企業のカスタムフィールドを取得
	for i := range companies {
		if err := r.db.Where("company_id = ?", companies[i].ID).Find(&companies[i].CustomFields).Error; err != nil {
			return nil, 0, err
		}
	}

	return companies, total, nil
}
