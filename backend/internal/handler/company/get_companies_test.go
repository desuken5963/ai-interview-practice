package company

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// モックユースケースの定義
type MockCompanyUseCase struct {
	mock.Mock
}

func (m *MockCompanyUseCase) GetCompanies(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CompanyResponse), args.Error(1)
}

func (m *MockCompanyUseCase) GetCompanyByID(ctx context.Context, id int) (*entity.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Company), args.Error(1)
}

func (m *MockCompanyUseCase) CreateCompany(ctx context.Context, company *entity.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockCompanyUseCase) UpdateCompany(ctx context.Context, company *entity.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *MockCompanyUseCase) DeleteCompany(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetCompanies(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		queryParams    string
		mockResponse   *entity.CompanyResponse
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "正常に企業一覧を取得できる",
			queryParams: "?page=1&limit=10",
			mockResponse: &entity.CompanyResponse{
				Companies: []entity.Company{
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
				Total: 2,
				Page:  1,
				Limit: 10,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"companies": []interface{}{
					map[string]interface{}{
						"id":                   float64(1),
						"name":                 "テスト企業1",
						"business_description": "テスト企業1の説明",
						"custom_fields":        nil,
						"job_count":            float64(0),
						"created_at":           "",
						"updated_at":           "",
					},
					map[string]interface{}{
						"id":                   float64(2),
						"name":                 "テスト企業2",
						"business_description": "テスト企業2の説明",
						"custom_fields":        nil,
						"job_count":            float64(0),
						"created_at":           "",
						"updated_at":           "",
					},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "不正なページパラメータの場合はエラーを返す",
			queryParams:    "?page=invalid&limit=10",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_PAGE",
					"message": "ページは1以上の整数である必要があります",
				},
			},
		},
		{
			name:           "不正なリミットパラメータの場合はエラーを返す",
			queryParams:    "?page=1&limit=invalid",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_LIMIT",
					"message": "リミットは1から100の間の整数である必要があります",
				},
			},
		},
		{
			name:           "サーバーエラーの場合は500エラーを返す",
			queryParams:    "?page=1&limit=10",
			mockResponse:   nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Ginのテストモードを設定
			gin.SetMode(gin.TestMode)

			// モックユースケースの作成
			mockUseCase := new(MockCompanyUseCase)

			// 正常なパラメータの場合のみモックの振る舞いを設定
			if tt.mockResponse != nil || tt.mockError != nil {
				mockUseCase.On("GetCompanies", mock.Anything, 1, 10).
					Return(tt.mockResponse, tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.GET("/api/v1/companies", GetCompanies(mockUseCase))

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			// リクエストを実行
			router.ServeHTTP(rec, req)

			// レスポンスを検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// JSONレスポンスをパース
			var response map[string]interface{}
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// 日付フィールドは動的に生成されるため、テスト対象から除外
			if companies, ok := response["companies"].([]interface{}); ok {
				for _, company := range companies {
					if c, ok := company.(map[string]interface{}); ok {
						delete(c, "created_at")
						delete(c, "updated_at")
					}
				}
			}

			// 期待されるレスポンスボディを検証
			if tt.expectedBody != nil {
				// 日付フィールドは動的に生成されるため、期待値からも削除
				if companies, ok := tt.expectedBody["companies"].([]interface{}); ok {
					for _, company := range companies {
						if c, ok := company.(map[string]interface{}); ok {
							delete(c, "created_at")
							delete(c, "updated_at")
						}
					}
				}

				assert.Equal(t, tt.expectedBody, response)
			}

			// モックが期待通り呼ばれたことを検証
			mockUseCase.AssertExpectations(t)
		})
	}
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}
