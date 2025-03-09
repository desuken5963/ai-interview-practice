package company

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// モックリポジトリの定義
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) Create(ctx context.Context, company *entity.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockCompanyRepository) FindByID(ctx context.Context, id int) (*entity.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Company), args.Error(1)
}

func (m *MockCompanyRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Company, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]entity.Company), args.Get(1).(int64), args.Error(2)
}

func (m *MockCompanyRepository) Update(ctx context.Context, company *entity.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockCompanyRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCompanyUseCase_GetCompanies(t *testing.T) {
	// テストケース
	tests := []struct {
		name          string
		page          int
		limit         int
		mockCompanies []entity.Company
		mockTotal     int64
		mockError     error
		expectedPage  int
		expectedLimit int
	}{
		{
			name:  "正常に企業一覧を取得できる",
			page:  1,
			limit: 10,
			mockCompanies: []entity.Company{
				{
					ID:                  1,
					Name:                "テスト企業1",
					BusinessDescription: stringPtr("テスト企業1の説明"),
				},
				{
					ID:                  2,
					Name:                "テスト企業2",
					BusinessDescription: stringPtr("テスト企業2の説明"),
				},
			},
			mockTotal:     2,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "ページが0以下の場合はデフォルト値が使用される",
			page:          0,
			limit:         10,
			mockCompanies: []entity.Company{},
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  1, // デフォルト値
			expectedLimit: 10,
		},
		{
			name:          "リミットが0以下の場合はデフォルト値が使用される",
			page:          1,
			limit:         0,
			mockCompanies: []entity.Company{},
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10, // デフォルト値
		},
		{
			name:          "リミットが100より大きい場合はデフォルト値が使用される",
			page:          1,
			limit:         101,
			mockCompanies: []entity.Company{},
			mockTotal:     0,
			mockError:     nil,
			expectedPage:  1,
			expectedLimit: 10, // デフォルト値
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックリポジトリの作成
			mockRepo := new(MockCompanyRepository)

			// モックの振る舞いを設定
			mockRepo.On("FindAll", mock.Anything, tt.expectedPage, tt.expectedLimit).
				Return(tt.mockCompanies, tt.mockTotal, tt.mockError)

			// テスト対象のユースケースを作成
			useCase := NewCompanyUseCase(mockRepo)

			// テスト実行
			result, err := useCase.GetCompanies(context.Background(), tt.page, tt.limit)

			// 検証
			if tt.mockError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.mockError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockCompanies, result.Companies)
				assert.Equal(t, tt.mockTotal, result.Total)
				assert.Equal(t, tt.expectedPage, result.Page)
				assert.Equal(t, tt.expectedLimit, result.Limit)
			}

			// モックが期待通り呼ばれたことを検証
			mockRepo.AssertExpectations(t)
		})
	}
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}
