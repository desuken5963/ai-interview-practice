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
)

// DeleteCompanyMockUseCase はテスト用のモックです
type DeleteCompanyMockUseCase struct {
	mock.Mock
}

func (m *DeleteCompanyMockUseCase) Execute(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestDeleteCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*DeleteCompanyMockUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "正常系: 企業情報の削除",
			id:   "1",
			mockSetup: func(m *DeleteCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 1).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name: "異常系: 無効なID",
			id:   "invalid",
			mockSetup: func(m *DeleteCompanyMockUseCase) {
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
			name: "異常系: 企業が存在しない",
			id:   "999",
			mockSetup: func(m *DeleteCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 999).Return(errors.New("not found"))
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
			mockSetup: func(m *DeleteCompanyMockUseCase) {
				m.On("Execute", mock.Anything, 1).Return(errors.New("server error"))
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
			mockUC := new(DeleteCompanyMockUseCase)
			tt.mockSetup(mockUC)

			// ハンドラーの作成
			handler := NewDeleteCompanyHandler(mockUC)

			// テスト用のGinコンテキストの作成
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			c.Request = httptest.NewRequest("DELETE", "/api/v1/companies/"+tt.id, nil)

			// ハンドラーの実行
			handler.Handle(c)

			// アサーション
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
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
