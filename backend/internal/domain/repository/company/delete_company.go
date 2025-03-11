package repository

import (
	"context"
)

// DeleteCompanyRepository は企業情報を削除するためのリポジトリインターフェースです
type DeleteCompanyRepository interface {
	Delete(ctx context.Context, id int) error
}
