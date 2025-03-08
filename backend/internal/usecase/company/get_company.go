package company

import (
	"context"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetCompanyByID は指定されたIDの企業情報を取得します
func (u *companyUseCase) GetCompanyByID(ctx context.Context, id int) (*entity.Company, error) {
	return u.companyRepo.FindByID(ctx, id)
}
