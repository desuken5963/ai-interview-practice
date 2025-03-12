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
	"github.com/takanoakira/ai-interview-practice/backend/test"
)

// GetCompanyMockUseCase はテスト用のモックです
type GetCompanyMockUseCase struct {
	mock.Mock
}

func (m *GetCompanyMockUseCase) Execute(ctx context.Context, id int) (*entity.Company, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Company), args.Error(1)
}

func TestGetCompany(t *testing.T) {
	// テスト用の企業データ
	now := time.Now()
	mockCompany := &entity.Company{
		ID:                  1,
		Name:                "テスト企業",
		BusinessDescription: test.StringPtr("テスト企業の説明"),
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
	}

	// テストケース
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*GetCompanyMockUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "正常系: 企業情報の取得",
			id:   "1",
			mockSetup: func(m *GetCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 1).Return(mockCompany, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   mockCompany,
		},
		{
			name: "異常系: 無効なID",
			id:   "invalid",
			mockSetup: func(m *GetCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 0).Return(nil, errors.New("invalid id"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "IDは整数である必要があります",
				},
			},
		},
		{
			name: "異常系: 企業が存在しない",
			id:   "999",
			mockSetup: func(m *GetCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 999).Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "企業が見つかりません",
				},
			},
		},
		{
			name: "異常系: サーバーエラー",
			id:   "1",
			mockSetup: func(m *GetCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 1).Return(nil, errors.New("server error"))
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
			mockUC := new(GetCompanyMockUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewGetCompanyHandler(mockUC)

			// テスト用のGinコンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			c.Request = httptest.NewRequest("GET", "/api/v1/companies/"+tt.id, nil)

			// ハンドラーの実行
			handler.Handle(c)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response entity.Company
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody.(*entity.Company), &response)
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
