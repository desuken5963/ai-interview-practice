package company

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetCompaniesMockUseCase はテスト用のモックです
type GetCompaniesMockUseCase struct {
	mock.Mock
}

func (m *GetCompaniesMockUseCase) Execute(ctx context.Context, page, limit int) (*entity.CompanyResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CompanyResponse), args.Error(1)
}

func TestGetCompanies(t *testing.T) {
	// テスト用の企業データ
	now := time.Now()
	mockCompanies := []entity.Company{
		{
			ID:                  1,
			Name:                "テスト企業1",
			BusinessDescription: stringPtr("テスト企業1の説明"),
			CustomFields: []entity.CompanyCustomField{
				{
					ID:        1,
					CompanyID: 1,
					FieldName: "業界",
					Content:   "IT",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:                  2,
			Name:                "テスト企業2",
			BusinessDescription: stringPtr("テスト企業2の説明"),
			CustomFields: []entity.CompanyCustomField{
				{
					ID:        2,
					CompanyID: 2,
					FieldName: "業界",
					Content:   "金融",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockResponse := &entity.CompanyResponse{
		Companies: mockCompanies,
		Total:     2,
		Page:      1,
		Limit:     10,
	}

	// テストケース
	tests := []struct {
		name           string
		query          string
		mockSetup      func(*GetCompaniesMockUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "正常系: 企業一覧の取得",
			query: "page=1&limit=10",
			mockSetup: func(m *GetCompaniesMockUseCase) {
				m.On("Execute", mock.Anything, 1, 10).Return(mockResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   mockResponse,
		},
		{
			name:  "異常系: 無効なページ番号",
			query: "page=0&limit=10",
			mockSetup: func(m *GetCompaniesMockUseCase) {
				m.On("Execute", mock.Anything, 0, 10).Return(nil, errors.New("invalid page"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": gin.H{
					"code":    "INVALID_PAGE",
					"message": "ページは1以上の整数である必要があります",
				},
			},
		},
		{
			name:  "異常系: 無効なリミット",
			query: "page=1&limit=101",
			mockSetup: func(m *GetCompaniesMockUseCase) {
				m.On("Execute", mock.Anything, 1, 101).Return(nil, errors.New("invalid limit"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": gin.H{
					"code":    "INVALID_LIMIT",
					"message": "リミットは1から100の間の整数である必要があります",
				},
			},
		},
		{
			name:  "異常系: サーバーエラー",
			query: "page=1&limit=10",
			mockSetup: func(m *GetCompaniesMockUseCase) {
				m.On("Execute", mock.Anything, 1, 10).Return(nil, errors.New("server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックの設定
			mockUC := new(GetCompaniesMockUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewGetCompaniesHandler(mockUC)

			// テスト用のGinコンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/v1/companies?"+tt.query, nil)

			// ハンドラーの実行
			handler.Handle(c)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.CompanyResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(*entity.CompanyResponse), &response)
			} else {
				var response gin.H
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			// モックの検証
			mockUC.AssertExpectations(t)
		})
	}
}

// stringPtr は文字列のポインタを返すヘルパー関数です
func stringPtr(s string) *string {
	return &s
}
