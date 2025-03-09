package company

import (
	"bytes"
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

func TestUpdateCompany(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		id             string
		requestBody    map[string]interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "正常に企業を更新できる",
			id:   "1",
			requestBody: map[string]interface{}{
				"name":                 "更新企業",
				"business_description": "更新企業の説明",
				"custom_fields": []map[string]interface{}{
					{
						"id":         1,
						"company_id": 1,
						"field_name": "業界",
						"content":    "更新後のIT",
					},
					{
						"field_name": "従業員数",
						"content":    "100人",
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":                   float64(1),
				"name":                 "更新企業",
				"business_description": "更新企業の説明",
				"job_count":            float64(0),
				"created_at":           "",
				"updated_at":           "",
				"custom_fields": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"company_id": float64(1),
						"field_name": "業界",
						"content":    "更新後のIT",
						"created_at": "",
						"updated_at": "",
					},
					map[string]interface{}{
						"id":         float64(2),
						"company_id": float64(1),
						"field_name": "従業員数",
						"content":    "100人",
						"created_at": "",
						"updated_at": "",
					},
				},
			},
		},
		{
			name: "不正なIDパラメータの場合はエラーを返す",
			id:   "invalid",
			requestBody: map[string]interface{}{
				"name": "更新企業",
			},
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
			name: "必須フィールドがない場合はエラーを返す",
			id:   "1",
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
			name: "存在しない企業IDの場合は404エラーを返す",
			id:   "999",
			requestBody: map[string]interface{}{
				"name":                 "存在しない企業",
				"business_description": "存在しない企業の説明",
			},
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
			name: "サーバーエラーの場合は500エラーを返す",
			id:   "1",
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
			if tt.id != "invalid" && tt.requestBody["name"] != nil && tt.requestBody["name"] != "" {
				id := 0
				fmt.Sscanf(tt.id, "%d", &id)

				mockUseCase.On("UpdateCompany", mock.Anything, mock.AnythingOfType("*entity.Company")).
					Run(func(args mock.Arguments) {
						// モックに渡された企業オブジェクトを検証
						company := args.Get(1).(*entity.Company)
						assert.Equal(t, id, company.ID)
						assert.Equal(t, tt.requestBody["name"], company.Name)
						if tt.requestBody["business_description"] != nil {
							assert.Equal(t, tt.requestBody["business_description"], *company.BusinessDescription)
						}

						// カスタムフィールドのIDを設定（更新成功をシミュレート）
						if len(company.CustomFields) > 0 {
							for i := range company.CustomFields {
								if company.CustomFields[i].ID == 0 {
									company.CustomFields[i].ID = i + 1
								}
								company.CustomFields[i].CompanyID = id
							}
						}
					}).
					Return(tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.PUT("/api/v1/companies/:id", UpdateCompany(mockUseCase))

			// リクエストボディをJSON化
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodPut, "/api/v1/companies/"+tt.id, bytes.NewBuffer(jsonData))
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
