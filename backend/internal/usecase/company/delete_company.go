package company

import (
	"context"
)

// DeleteCompany は指定されたIDの企業情報を削除します
func (u *companyUseCase) DeleteCompany(ctx context.Context, id int) error {
	return u.companyRepo.Delete(ctx, id)
}
