package company

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		id             string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "正常に企業を削除できる",
			id:             "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "不正なIDパラメータの場合はエラーを返す",
			id:             "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_ID",
					"message": "IDは整数である必要があります",
				},
			},
		},
		{
			name:           "存在しない企業IDの場合は404エラーを返す",
			id:             "999",
			mockError:      errors.New("企業が見つかりません"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			},
		},
		{
			name:           "サーバーエラーの場合は500エラーを返す",
			id:             "1",
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
			if tt.id != "invalid" {
				id := 0
				fmt.Sscanf(tt.id, "%d", &id)
				mockUseCase.On("DeleteCompany", mock.Anything, id).
					Return(tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.DELETE("/api/v1/companies/:id", DeleteCompany(mockUseCase))

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/companies/"+tt.id, nil)
			rec := httptest.NewRecorder()

			// リクエストを実行
			router.ServeHTTP(rec, req)

			// レスポンスを検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// 204 No Contentの場合はボディが空であることを確認
			if tt.expectedStatus == http.StatusNoContent {
				assert.Empty(t, rec.Body.String())
			} else {
				// JSONレスポンスをパース
				var response map[string]interface{}
				if rec.Body.Len() > 0 {
					err := json.NewDecoder(rec.Body).Decode(&response)
					assert.NoError(t, err)

					// 期待されるレスポンスボディを検証
					if tt.expectedBody != nil {
						assert.Equal(t, tt.expectedBody, response)
					}
				}
			}

			// モックが期待通り呼ばれたことを検証
			mockUseCase.AssertExpectations(t)
		})
	}
}
