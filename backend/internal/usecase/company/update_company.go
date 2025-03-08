package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// UpdateCompany は既存の企業情報を更新します
func (u *companyUseCase) UpdateCompany(ctx context.Context, company *entity.Company) error {
	return u.companyRepo.Update(ctx, company)
}
