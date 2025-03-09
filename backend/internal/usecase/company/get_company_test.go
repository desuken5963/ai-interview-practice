package company

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

func TestCompanyUseCase_GetCompanyByID(t *testing.T) {
	// テストケース
	tests := []struct {
		name          string
		id            int
		mockCompany   *entity.Company
		mockError     error
		expectedError bool
	}{
		{
			name: "正常に企業を取得できる",
			id:   1,
			mockCompany: &entity.Company{
				ID:                  1,
				Name:                "テスト企業",
				BusinessDescription: stringPtr("テスト企業の説明"),
				CustomFields: []entity.CompanyCustomField{
					{
						ID:        1,
						CompanyID: 1,
						FieldName: "業界",
						Content:   "IT",
					},
				},
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "存在しない企業IDではエラーになる",
			id:            999,
			mockCompany:   nil,
			mockError:     errors.New("企業が見つかりません"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの作成
			mockRepo := new(MockCompanyRepository)

			// モックの振る舞いを設定
			mockRepo.On("FindByID", mock.Anything, tt.id).
				Return(tt.mockCompany, tt.mockError)

			// テスト対象のユースケースを作成
			useCase := NewCompanyUseCase(mockRepo)

			// テスト実行
			company, err := useCase.GetCompanyByID(context.Background(), tt.id)

			// 検証
			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err)
				assert.Nil(t, company)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, company)
				assert.Equal(t, tt.mockCompany.ID, company.ID)
				assert.Equal(t, tt.mockCompany.Name, company.Name)
				assert.Equal(t, tt.mockCompany.BusinessDescription, company.BusinessDescription)
				assert.Equal(t, len(tt.mockCompany.CustomFields), len(company.CustomFields))
				if len(tt.mockCompany.CustomFields) > 0 {
					assert.Equal(t, tt.mockCompany.CustomFields[0].FieldName, company.CustomFields[0].FieldName)
					assert.Equal(t, tt.mockCompany.CustomFields[0].Content, company.CustomFields[0].Content)
				}
			}

			// モックが期待通り呼ばれたことを検証
			mockRepo.AssertExpectations(t)
		})
	}
}
