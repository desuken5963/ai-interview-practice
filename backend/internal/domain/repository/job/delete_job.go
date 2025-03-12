package repository

import (
	"context"
)

// DeleteJobRepository は求人情報を削除するためのリポジトリインターフェースです
type DeleteJobRepository interface {
	Delete(ctx context.Context, id int) error
}
