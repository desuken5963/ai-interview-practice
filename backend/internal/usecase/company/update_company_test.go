package company

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

func TestCompanyUseCase_UpdateCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name          string
		company       *entity.Company
		mockError     error
		expectedError bool
	}{
		{
			name: "正常に企業を更新できる",
			company: &entity.Company{
				ID:                  1,
				Name:                "更新企業",
				BusinessDescription: stringPtr("更新企業の説明"),
				CustomFields: []entity.CompanyCustomField{
					{
						ID:        1,
						CompanyID: 1,
						FieldName: "業界",
						Content:   "更新後のIT",
					},
				},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "リポジトリでエラーが発生した場合はエラーを返す",
			company: &entity.Company{
				ID:                  999,
				Name:                "存在しない企業",
				BusinessDescription: stringPtr("存在しない企業の説明"),
			},
			mockError:     errors.New("企業が見つかりません"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの作成
			mockRepo := new(MockCompanyRepository)

			// モックの振る舞いを設定
			mockRepo.On("Update", mock.Anything, tt.company).
				Return(tt.mockError)

			// テスト対象のユースケースを作成
			useCase := NewCompanyUseCase(mockRepo)

			// テスト実行
			err := useCase.UpdateCompany(context.Background(), tt.company)

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
