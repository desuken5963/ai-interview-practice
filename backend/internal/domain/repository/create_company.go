package repository

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// CreateCompanyRepository は企業情報を作成するためのリポジトリインターフェースです
type CreateCompanyRepository interface {
	Create(ctx context.Context, company *entity.Company) error
}
