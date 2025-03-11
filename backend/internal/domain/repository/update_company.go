package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// UpdateCompanyRepository は企業情報を更新するためのリポジトリインターフェースです
type UpdateCompanyRepository interface {
	Update(ctx context.Context, company *entity.Company) error
}
