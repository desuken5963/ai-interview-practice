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

// CreateCompanyMockUseCase はテスト用のモックです
type CreateCompanyMockUseCase struct {
	mock.Mock
}

func (m *CreateCompanyMockUseCase) Execute(ctx context.Context, company *entity.Company) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func TestCreateCompany(t *testing.T) {
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
		requestBody    interface{}
		mockSetup      func(*CreateCompanyMockUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "正常系: 企業情報の作成",
			requestBody: gin.H{
				"name":                 "テスト企業",
				"business_description": "テスト企業の説明",
				"custom_fields": []gin.H{
					{
						"field_name": "業界",
						"content":    "IT",
					},
				},
			},
			mockSetup: func(m *CreateCompanyMockUseCase) {
				m.On("Execute", mock.Anything, mock.MatchedBy(func(c *entity.Company) bool {
					return c.Name == "テスト企業" &&
						*c.BusinessDescription == "テスト企業の説明" &&
						len(c.CustomFields) == 1 &&
						c.CustomFields[0].FieldName == "業界" &&
						c.CustomFields[0].Content == "IT"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   mockCompany,
		},
		{
			name: "異常系: 必須フィールドの欠落",
			requestBody: gin.H{
				"business_description": "テスト企業の説明",
			},
			mockSetup: func(m *CreateCompanyMockUseCase) {
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
			name: "異常系: サーバーエラー",
			requestBody: gin.H{
				"name":                 "テスト企業",
				"business_description": "テスト企業の説明",
			},
			mockSetup: func(m *CreateCompanyMockUseCase) {
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
			mockUC := new(CreateCompanyMockUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewCreateCompanyHandler(mockUC)

			// リクエストボディの作成
			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// テスト用のGinコンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/companies", bytes.NewBuffer(body))

			// ハンドラーの実行
			handler.Handle(c)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
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
