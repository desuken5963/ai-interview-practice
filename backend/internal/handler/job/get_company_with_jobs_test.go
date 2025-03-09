package job

import (
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

func TestGetCompanyWithJobs(t *testing.T) {
	now := time.Now()
	mockCompany := &entity.Company{
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
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockJobs := []entity.JobPosting{
		{
			ID:          1,
			CompanyID:   1,
			Title:       "テストエンジニア1",
			Description: stringPtr("テスト求人1の説明"),
			CustomFields: []entity.JobCustomField{
				{
					ID:        1,
					JobID:     1,
					FieldName: "勤務地",
					Content:   "東京",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          2,
			CompanyID:   1,
			Title:       "テストエンジニア2",
			Description: stringPtr("テスト求人2の説明"),
			CustomFields: []entity.JobCustomField{
				{
					ID:        2,
					JobID:     2,
					FieldName: "勤務地",
					Content:   "大阪",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// テストケース
	tests := []struct {
		name           string
		companyID      string
		mockCompany    *entity.Company
		mockJobs       []entity.JobPosting
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "正常に企業と求人情報を取得できる",
			companyID:      "1",
			mockCompany:    mockCompany,
			mockJobs:       mockJobs,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":                   float64(1),
				"name":                 "テスト企業",
				"business_description": "テスト企業の説明",
				"job_count":            float64(2),
				"custom_fields": []interface{}{
					map[string]interface{}{
						"id":         float64(1),
						"company_id": float64(1),
						"field_name": "業界",
						"content":    "IT",
					},
				},
				"jobs": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"company_id":  float64(1),
						"title":       "テストエンジニア1",
						"description": "テスト求人1の説明",
						"custom_fields": []interface{}{
							map[string]interface{}{
								"id":         float64(1),
								"job_id":     float64(1),
								"field_name": "勤務地",
								"content":    "東京",
							},
						},
					},
					map[string]interface{}{
						"id":          float64(2),
						"company_id":  float64(1),
						"title":       "テストエンジニア2",
						"description": "テスト求人2の説明",
						"custom_fields": []interface{}{
							map[string]interface{}{
								"id":         float64(2),
								"job_id":     float64(2),
								"field_name": "勤務地",
								"content":    "大阪",
							},
						},
					},
				},
			},
		},
		{
			name:           "不正なIDパラメータの場合はエラーを返す",
			companyID:      "invalid",
			mockCompany:    nil,
			mockJobs:       nil,
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
			name:           "サーバーエラーの場合は500エラーを返す",
			companyID:      "1",
			mockCompany:    nil,
			mockJobs:       nil,
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

			// 正常なパラメータの場合のみモックの振る舞いを設定
			if tt.companyID != "invalid" {
				companyID := 1
				mockUseCase.On("GetCompanyWithJobs", mock.Anything, companyID).
					Return(tt.mockCompany, tt.mockJobs, tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.GET("/api/v1/companies/:id/with-jobs", GetCompanyWithJobs(mockUseCase))

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+tt.companyID+"/with-jobs", nil)
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
			if jobs, ok := response["jobs"].([]interface{}); ok {
				for _, job := range jobs {
					if j, ok := job.(map[string]interface{}); ok {
						delete(j, "created_at")
						delete(j, "updated_at")
						if customFields, ok := j["custom_fields"].([]interface{}); ok {
							for _, field := range customFields {
								if cf, ok := field.(map[string]interface{}); ok {
									delete(cf, "created_at")
									delete(cf, "updated_at")
								}
							}
						}
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
				if jobs, ok := tt.expectedBody["jobs"].([]interface{}); ok {
					for _, job := range jobs {
						if j, ok := job.(map[string]interface{}); ok {
							delete(j, "created_at")
							delete(j, "updated_at")
							if customFields, ok := j["custom_fields"].([]interface{}); ok {
								for _, field := range customFields {
									if cf, ok := field.(map[string]interface{}); ok {
										delete(cf, "created_at")
										delete(cf, "updated_at")
									}
								}
							}
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
