package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// CompanyRepository は企業情報の永続化を担当するインターフェースです
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
