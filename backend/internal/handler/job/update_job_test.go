package job

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

func TestUpdateJob(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		companyID      string
		jobID          string
		requestBody    map[string]interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "正常に求人を更新できる",
			companyID: "1",
			jobID:     "1",
			requestBody: map[string]interface{}{
				"title":       "更新後のテストエンジニア",
				"description": "更新後のテスト求人の説明",
				"custom_fields": []map[string]interface{}{
					{
						"id":         1,
						"job_id":     1,
						"field_name": "勤務地",
						"content":    "更新後の東京",
					},
					{
						"id":         2,
						"job_id":     1,
						"field_name": "給与",
						"content":    "更新後の年収600万円〜900万円",
					},
					{
						"field_name": "新しいフィールド",
						"content":    "新しい値",
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"company_id":  float64(1),
				"title":       "更新後のテストエンジニア",
				"description": "更新後のテスト求人の説明",
				"custom_fields": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"job_id":     float64(1),
						"field_name": "勤務地",
						"content":    "更新後の東京",
					},
					map[string]interface{}{
						"id":         float64(2),
						"job_id":     float64(1),
						"field_name": "給与",
						"content":    "更新後の年収600万円〜900万円",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "新しいフィールド",
						"content":    "新しい値",
					},
				},
			},
		},
		{
			name:      "不正な企業IDパラメータの場合はエラーを返す",
			companyID: "invalid",
			jobID:     "1",
			requestBody: map[string]interface{}{
				"title": "更新後のテストエンジニア",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_ID",
					"message": "企業IDは整数である必要があります",
				},
			},
		},
		{
			name:      "不正な求人IDパラメータの場合はエラーを返す",
			companyID: "1",
			jobID:     "invalid",
			requestBody: map[string]interface{}{
				"title": "更新後のテストエンジニア",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_ID",
					"message": "求人IDは整数である必要があります",
				},
			},
		},
		{
			name:      "必須フィールドがない場合はエラーを返す",
			companyID: "1",
			jobID:     "1",
			requestBody: map[string]interface{}{
				"description": "タイトルのない求人の説明",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "VALIDATION_ERROR",
					"message": "バリデーションエラーが発生しました",
					"details": []interface{}{
						map[string]interface{}{
							"field":   "title",
							"message": "求人タイトルは必須です",
						},
					},
				},
			},
		},
		{
			name:      "サーバーエラーの場合は500エラーを返す",
			companyID: "1",
			jobID:     "1",
			requestBody: map[string]interface{}{
				"title":       "エラー求人",
				"description": "エラー求人の説明",
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
			mockUseCase := new(MockJobUseCase)

			// 正常なリクエストの場合のみモックの振る舞いを設定
			if tt.companyID != "invalid" && tt.jobID != "invalid" && tt.requestBody["title"] != nil && tt.requestBody["title"] != "" {
				mockUseCase.On("UpdateJob", mock.Anything, mock.AnythingOfType("*entity.JobPosting")).
					Run(func(args mock.Arguments) {
						// モックに渡された求人オブジェクトを検証
						job := args.Get(1).(*entity.JobPosting)
						assert.Equal(t, tt.requestBody["title"], job.Title)
						if tt.requestBody["description"] != nil {
							assert.Equal(t, tt.requestBody["description"], *job.Description)
						}
					}).
					Return(tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.PUT("/api/v1/companies/:id/jobs/:job_id", UpdateJob(mockUseCase))

			// リクエストボディをJSON化
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// テストリクエストを作成
			url := fmt.Sprintf("/api/v1/companies/%s/jobs/%s", tt.companyID, tt.jobID)
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
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
			if _, ok := response["created_at"]; ok {
				delete(response, "created_at")
				delete(response, "updated_at")
			}
			if customFields, ok := response["custom_fields"].([]interface{}); ok {
				for _, field := range customFields {
					if cf, ok := field.(map[string]interface{}); ok {
						delete(cf, "created_at")
						delete(cf, "updated_at")
					}
				}
			}

			// 期待されるレスポンスボディを検証
			if tt.expectedBody != nil {
				// 日付フィールドは動的に生成されるため、期待値からも削除
				if _, ok := tt.expectedBody["created_at"]; ok {
					delete(tt.expectedBody, "created_at")
					delete(tt.expectedBody, "updated_at")
				}
				if customFields, ok := tt.expectedBody["custom_fields"].([]interface{}); ok {
					for _, field := range customFields {
						if cf, ok := field.(map[string]interface{}); ok {
							delete(cf, "created_at")
							delete(cf, "updated_at")
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
