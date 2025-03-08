package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// CreateCompany は新しい企業情報を作成します
func (u *companyUseCase) CreateCompany(ctx context.Context, company *entity.Company) error {
	return u.companyRepo.Create(ctx, company)
}
