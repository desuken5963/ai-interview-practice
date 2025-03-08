package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/repository"
	"gorm.io/gorm"
)

// CompanyRepository は企業情報のリポジトリインターフェースです
type CompanyRepository interface {
	// FindAll は企業情報の一覧を取得します
	FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error)

	// FindByID は指定されたIDの企業情報を取得します
	FindByID(ctx context.Context, id int) (*entity.Company, error)

	// Create は新しい企業情報を作成します
	Create(ctx context.Context, company *entity.Company) error

	// Update は既存の企業情報を更新します
	Update(ctx context.Context, company *entity.Company) error

	// Delete は指定されたIDの企業情報を削除します
	Delete(ctx context.Context, id int) error
}

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
