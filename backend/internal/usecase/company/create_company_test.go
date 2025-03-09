package company

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

func TestCompanyUseCase_CreateCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name          string
		company       *entity.Company
		mockError     error
		expectedError bool
	}{
		{
			name: "正常に企業を作成できる",
			company: &entity.Company{
				Name:                "新規企業",
				BusinessDescription: stringPtr("新規企業の説明"),
				CustomFields: []entity.CompanyCustomField{
					{
						FieldName: "業界",
						Content:   "IT",
					},
				},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name: "リポジトリでエラーが発生した場合はエラーを返す",
			company: &entity.Company{
				Name:                "エラー企業",
				BusinessDescription: stringPtr("エラー企業の説明"),
			},
			mockError:     errors.New("データベースエラー"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの作成
			mockRepo := new(MockCompanyRepository)

			// モックの振る舞いを設定
			mockRepo.On("Create", mock.Anything, tt.company).
				Return(tt.mockError)

			// テスト対象のユースケースを作成
			useCase := NewCompanyUseCase(mockRepo)

			// テスト実行
			err := useCase.CreateCompany(context.Background(), tt.company)

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
