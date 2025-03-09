package company

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCompanyUseCase_DeleteCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name          string
		id            int
		mockError     error
		expectedError bool
	}{
		{
			name:          "正常に企業を削除できる",
			id:            1,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "存在しない企業IDではエラーになる",
			id:            999,
			mockError:     errors.New("企業が見つかりません"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの作成
			mockRepo := new(MockCompanyRepository)

			// モックの振る舞いを設定
			mockRepo.On("Delete", mock.Anything, tt.id).
				Return(tt.mockError)

			// テスト対象のユースケースを作成
			useCase := NewCompanyUseCase(mockRepo)

			// テスト実行
			err := useCase.DeleteCompany(context.Background(), tt.id)

			// 検証
			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
			}

			// モックが期待通り呼ばれたことを検証
			mockRepo.AssertExpectations(t)
		})
	}
}
