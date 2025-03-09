package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

func TestGetJobsByCompanyID(t *testing.T) {
	// テスト用の求人データ
	now := time.Now()
	companyID := 1
	mockJobs := []entity.JobPosting{
		{
			ID:          1,
			CompanyID:   companyID,
			Title:       "テストエンジニア1",
			Description: stringPtr("テスト求人1の説明"),
			CustomFields: []entity.JobCustomField{
				{
					ID:        1,
					JobID:     1,
					FieldName: "勤務地",
					Content:   "東京",
				},
				{
					ID:        2,
					JobID:     1,
					FieldName: "給与",
					Content:   "年収500万円〜800万円",
				},
				{
					ID:        3,
					JobID:     1,
					FieldName: "雇用形態",
					Content:   "正社員",
				},
				{
					ID:        4,
					JobID:     1,
					FieldName: "応募締切",
					Content:   now.AddDate(0, 1, 0).Format("2006-01-02"),
				},
				{
					ID:        5,
					JobID:     1,
					FieldName: "ステータス",
					Content:   "公開中",
				},
				{
					ID:        6,
					JobID:     1,
					FieldName: "必須スキル",
					Content:   "Go, Docker, MySQL",
				},
				{
					ID:        7,
					JobID:     1,
					FieldName: "歓迎スキル",
					Content:   "Kubernetes, AWS",
				},
				{
					ID:        8,
					JobID:     1,
					FieldName: "経験年数",
					Content:   "3年以上",
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          2,
			CompanyID:   companyID,
			Title:       "テストエンジニア2",
			Description: stringPtr("テスト求人2の説明"),
			CustomFields: []entity.JobCustomField{
				{
					ID:        9,
					JobID:     2,
					FieldName: "勤務地",
					Content:   "大阪",
				},
				{
					ID:        10,
					JobID:     2,
					FieldName: "給与",
					Content:   "年収400万円〜700万円",
				},
				{
					ID:        11,
					JobID:     2,
					FieldName: "雇用形態",
					Content:   "契約社員",
				},
				{
					ID:        12,
					JobID:     2,
					FieldName: "応募締切",
					Content:   now.AddDate(0, 1, 0).Format("2006-01-02"),
				},
				{
					ID:        13,
					JobID:     2,
					FieldName: "ステータス",
					Content:   "公開中",
				},
				{
					ID:        14,
					JobID:     2,
					FieldName: "必須スキル",
					Content:   "Java, Spring, PostgreSQL",
				},
				{
					ID:        15,
					JobID:     2,
					FieldName: "歓迎スキル",
					Content:   "Docker, Kubernetes",
				},
				{
					ID:        16,
					JobID:     2,
					FieldName: "経験年数",
					Content:   "2年以上",
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
		queryParams    string
		mockResponse   *entity.JobResponse
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "正常に企業の求人一覧を取得できる",
			companyID:   "1",
			queryParams: "?page=1&limit=10",
			mockResponse: &entity.JobResponse{
				Jobs:  mockJobs,
				Total: 2,
				Page:  1,
				Limit: 10,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
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
							map[string]interface{}{
								"id":         float64(2),
								"job_id":     float64(1),
								"field_name": "給与",
								"content":    "年収500万円〜800万円",
							},
							map[string]interface{}{
								"id":         float64(3),
								"job_id":     float64(1),
								"field_name": "雇用形態",
								"content":    "正社員",
							},
							map[string]interface{}{
								"id":         float64(4),
								"job_id":     float64(1),
								"field_name": "応募締切",
								"content":    "2025-04-09",
							},
							map[string]interface{}{
								"id":         float64(5),
								"job_id":     float64(1),
								"field_name": "ステータス",
								"content":    "公開中",
							},
							map[string]interface{}{
								"id":         float64(6),
								"job_id":     float64(1),
								"field_name": "必須スキル",
								"content":    "Go, Docker, MySQL",
							},
							map[string]interface{}{
								"id":         float64(7),
								"job_id":     float64(1),
								"field_name": "歓迎スキル",
								"content":    "Kubernetes, AWS",
							},
							map[string]interface{}{
								"id":         float64(8),
								"job_id":     float64(1),
								"field_name": "経験年数",
								"content":    "3年以上",
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
								"id":         float64(9),
								"job_id":     float64(2),
								"field_name": "勤務地",
								"content":    "大阪",
							},
							map[string]interface{}{
								"id":         float64(10),
								"job_id":     float64(2),
								"field_name": "給与",
								"content":    "年収400万円〜700万円",
							},
							map[string]interface{}{
								"id":         float64(11),
								"job_id":     float64(2),
								"field_name": "雇用形態",
								"content":    "契約社員",
							},
							map[string]interface{}{
								"id":         float64(12),
								"job_id":     float64(2),
								"field_name": "応募締切",
								"content":    "2025-04-09",
							},
							map[string]interface{}{
								"id":         float64(13),
								"job_id":     float64(2),
								"field_name": "ステータス",
								"content":    "公開中",
							},
							map[string]interface{}{
								"id":         float64(14),
								"job_id":     float64(2),
								"field_name": "必須スキル",
								"content":    "Java, Spring, PostgreSQL",
							},
							map[string]interface{}{
								"id":         float64(15),
								"job_id":     float64(2),
								"field_name": "歓迎スキル",
								"content":    "Docker, Kubernetes",
							},
							map[string]interface{}{
								"id":         float64(16),
								"job_id":     float64(2),
								"field_name": "経験年数",
								"content":    "2年以上",
							},
						},
					},
				},
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "不正な企業IDパラメータの場合はエラーを返す",
			companyID:      "invalid",
			queryParams:    "",
			mockResponse:   nil,
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
			name:           "不正なページパラメータの場合はエラーを返す",
			companyID:      "1",
			queryParams:    "?page=invalid&limit=10",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_PAGE",
					"message": "ページは1以上の整数である必要があります",
				},
			},
		},
		{
			name:           "不正なリミットパラメータの場合はエラーを返す",
			companyID:      "1",
			queryParams:    "?page=1&limit=invalid",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_LIMIT",
					"message": "リミットは1から100の間の整数である必要があります",
				},
			},
		},
		{
			name:        "存在しない企業IDの場合は空の結果を返す",
			companyID:   "999",
			queryParams: "?page=1&limit=10",
			mockResponse: &entity.JobResponse{
				Jobs:  []entity.JobPosting{},
				Total: 0,
				Page:  1,
				Limit: 10,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"jobs":  []interface{}{},
				"total": float64(0),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "サーバーエラーの場合は500エラーを返す",
			companyID:      "1",
			queryParams:    "?page=1&limit=10",
			mockResponse:   nil,
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
			if tt.companyID != "invalid" && !containsInvalidParam(tt.queryParams) {
				companyID := 0
				fmt.Sscanf(tt.companyID, "%d", &companyID)
				mockUseCase.On("GetJobsByCompanyID", mock.Anything, companyID, 1, 10).
					Return(tt.mockResponse, tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.GET("/api/v1/companies/:id/jobs", GetJobsByCompanyID(mockUseCase))

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/companies/"+tt.companyID+"/jobs"+tt.queryParams, nil)
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

// containsInvalidParam はクエリパラメータに不正な値が含まれているかをチェックするヘルパー関数
func containsInvalidParam(queryParams string) bool {
	return queryParams == "?page=invalid&limit=10" || queryParams == "?page=1&limit=invalid"
}
