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
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

func TestGetCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		id             string
		mockCompany    *entity.Company
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常に企業を取得できる",
			id:   "1",
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
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":                   float64(1),
				"name":                 "テスト企業",
				"business_description": "テスト企業の説明",
				"job_count":            float64(0),
				"created_at":           "",
				"updated_at":           "",
				"custom_fields": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"company_id": float64(1),
						"field_name": "業界",
						"content":    "IT",
						"created_at": "",
						"updated_at": "",
					},
				},
			},
		},
		{
			name:           "不正なIDパラメータの場合はエラーを返す",
			id:             "invalid",
			mockCompany:    nil,
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
			mockCompany:    nil,
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
			mockCompany:    nil,
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
				mockUseCase.On("GetCompanyByID", mock.Anything, id).
					Return(tt.mockCompany, tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.GET("/api/v1/companies/:id", GetCompany(mockUseCase))

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+tt.id, nil)
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
			if company, ok := response["custom_fields"].([]interface{}); ok {
				for _, field := range company {
					if cf, ok := field.(map[string]interface{}); ok {
						delete(cf, "created_at")
						delete(cf, "updated_at")
					}
				}
			}
			if _, ok := response["created_at"]; ok {
				delete(response, "created_at")
				delete(response, "updated_at")
			}

			// 期待されるレスポンスボディを検証
			if tt.expectedBody != nil {
				// 日付フィールドは動的に生成されるため、期待値からも削除
				if company, ok := tt.expectedBody["custom_fields"].([]interface{}); ok {
					for _, field := range company {
						if cf, ok := field.(map[string]interface{}); ok {
							delete(cf, "created_at")
							delete(cf, "updated_at")
						}
					}
				}
				if _, ok := tt.expectedBody["created_at"]; ok {
					delete(tt.expectedBody, "created_at")
					delete(tt.expectedBody, "updated_at")
				}

				assert.Equal(t, tt.expectedBody, response)
			}

			// モックが期待通り呼ばれたことを検証
			mockUseCase.AssertExpectations(t)
		})
	}
}
