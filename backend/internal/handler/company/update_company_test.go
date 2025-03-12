package company

import (
	"bytes"
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

// UpdateCompanyMockUseCase はテスト用のモックです
type UpdateCompanyMockUseCase struct {
	mock.Mock
}

func (m *UpdateCompanyMockUseCase) Execute(ctx context.Context, company *entity.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func TestUpdateCompany(t *testing.T) {
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
		requestBody    interface{}
		mockSetup      func(*UpdateCompanyMockUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "正常系: 企業情報の更新",
			id:   "1",
			requestBody: gin.H{
				"name":                 "更新後の企業名",
				"business_description": "更新後の説明",
				"custom_fields": []gin.H{
					{
						"field_name": "業界",
						"content":    "金融",
					},
				},
			},
			mockSetup: func(m *UpdateCompanyMockUseCase) {
				m.On("Execute", mock.Anything, mock.MatchedBy(func(c *entity.Company) bool {
					return c.ID == 1 &&
						c.Name == "更新後の企業名" &&
						*c.BusinessDescription == "更新後の説明" &&
						len(c.CustomFields) == 1 &&
						c.CustomFields[0].FieldName == "業界" &&
						c.CustomFields[0].Content == "金融"
				})).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   mockCompany,
		},
		{
			name: "異常系: 無効なID",
			id:   "invalid",
			requestBody: gin.H{
				"name": "更新後の企業名",
			},
			mockSetup: func(m *UpdateCompanyMockUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).Return(errors.New("invalid id"))
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
			name: "異常系: 必須フィールドの欠落",
			id:   "1",
			requestBody: gin.H{
				"business_description": "更新後の説明",
			},
			mockSetup: func(m *UpdateCompanyMockUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).Return(errors.New("name is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "企業名は必須です",
				},
			},
		},
		{
			name: "異常系: 企業が存在しない",
			id:   "999",
			requestBody: gin.H{
				"name": "更新後の企業名",
			},
			mockSetup: func(m *UpdateCompanyMockUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).Return(errors.New("not found"))
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
			requestBody: gin.H{
				"name": "更新後の企業名",
			},
			mockSetup: func(m *UpdateCompanyMockUseCase) {
				m.On("Execute", mock.Anything, mock.Anything).Return(errors.New("server error"))
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
			mockUC := new(UpdateCompanyMockUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewUpdateCompanyHandler(mockUC)

			// リクエストボディの作成
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// テスト用のGinコンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			c.Request = httptest.NewRequest("PUT", "/api/v1/companies/"+tt.id, bytes.NewBuffer(body))

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
