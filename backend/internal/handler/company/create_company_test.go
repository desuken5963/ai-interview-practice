package company

import (
	"bytes"
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

func TestCreateCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常に企業を作成できる",
			requestBody: map[string]interface{}{
				"name":                 "新規企業",
				"business_description": "新規企業の説明",
				"custom_fields": []map[string]interface{}{
					{
						"field_name": "業界",
						"content":    "IT",
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":                   float64(1),
				"name":                 "新規企業",
				"business_description": "新規企業の説明",
				"job_count":            float64(0),
				"created_at":           "",
				"updated_at":           "",
				"custom_fields": []interface{}{
					map[string]interface{}{
						"id":         float64(0),
						"company_id": float64(0),
						"field_name": "業界",
						"content":    "IT",
						"created_at": "",
						"updated_at": "",
					},
				},
			},
		},
		{
			name: "必須フィールドがない場合はエラーを返す",
			requestBody: map[string]interface{}{
				"business_description": "名前のない企業の説明",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "VALIDATION_ERROR",
					"message": "バリデーションエラーが発生しました",
					"details": []interface{}{
						map[string]interface{}{
							"field":   "name",
							"message": "企業名は必須です",
						},
					},
				},
			},
		},
		{
			name: "サーバーエラーの場合は500エラーを返す",
			requestBody: map[string]interface{}{
				"name":                 "エラー企業",
				"business_description": "エラー企業の説明",
			},
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

			// 正常なリクエストの場合のみモックの振る舞いを設定
			if tt.requestBody["name"] != nil && tt.requestBody["name"] != "" {
				mockUseCase.On("CreateCompany", mock.Anything, mock.AnythingOfType("*entity.Company")).
					Run(func(args mock.Arguments) {
						// モックに渡された企業オブジェクトを検証
						company := args.Get(1).(*entity.Company)
						assert.Equal(t, tt.requestBody["name"], company.Name)
						if tt.requestBody["business_description"] != nil {
							assert.Equal(t, tt.requestBody["business_description"], *company.BusinessDescription)
						}

						// IDを設定（作成成功をシミュレート）
						company.ID = 1
					}).
					Return(tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.POST("/api/v1/companies", CreateCompany(mockUseCase))

			// リクエストボディをJSON化
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// リクエストを実行
			router.ServeHTTP(rec, req)

			// レスポンスを検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

			// JSONレスポンスをパース
			var response map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &response)
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
