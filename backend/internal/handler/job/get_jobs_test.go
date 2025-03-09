package job

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
)

// GetJobsのテスト用のモック
type GetJobsMockUseCase struct {
	mock.Mock
}

func (m *GetJobsMockUseCase) GetJobs(ctx context.Context, page, limit int) (*entity.JobResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobResponse), args.Error(1)
}

func (m *GetJobsMockUseCase) GetJobsByCompanyID(ctx context.Context, companyID, page, limit int) (*entity.JobResponse, error) {
	args := m.Called(ctx, companyID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobResponse), args.Error(1)
}

func (m *GetJobsMockUseCase) GetJob(ctx context.Context, id int) (*entity.JobPosting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobPosting), args.Error(1)
}

func (m *GetJobsMockUseCase) CreateJob(ctx context.Context, job *entity.JobPosting) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *GetJobsMockUseCase) UpdateJob(ctx context.Context, job *entity.JobPosting) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *GetJobsMockUseCase) DeleteJob(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *GetJobsMockUseCase) GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.Company), args.Get(1).([]entity.JobPosting), args.Error(2)
}

// GetJobsのテスト用のヘルパー関数
func getJobsStringPtr(s string) *string {
	return &s
}

func TestGetJobs(t *testing.T) {
	// テスト用の求人データ
	now := time.Now()
	mockJobs := []entity.JobPosting{
		{
			ID:          1,
			CompanyID:   1,
			Title:       "テストエンジニア1",
			Description: getJobsStringPtr("テスト求人の説明1"),
			CustomFields: []entity.JobCustomField{
				{
					ID:        1,
					JobID:     1,
					FieldName: "勤務地",
					Content:   "東京",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:          2,
			CompanyID:   1,
			Title:       "テストエンジニア2",
			Description: getJobsStringPtr("テスト求人の説明2"),
			CustomFields: []entity.JobCustomField{
				{
					ID:        2,
					JobID:     2,
					FieldName: "勤務地",
					Content:   "大阪",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	mockResponse := &entity.JobResponse{
		Jobs:  mockJobs,
		Total: 2,
		Page:  1,
		Limit: 10,
	}

	// テストケース
	tests := []struct {
		name           string
		query          string
		mockResponse   *entity.JobResponse
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "正常に求人一覧を取得できる",
			query:          "?page=1&limit=10",
			mockResponse:   mockResponse,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"jobs": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"company_id":  float64(1),
						"title":       "テストエンジニア1",
						"description": "テスト求人の説明1",
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
						"description": "テスト求人の説明2",
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
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "不正なページパラメータの場合はデフォルト値を使用する",
			query:          "?page=invalid&limit=10",
			mockResponse:   mockResponse,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"jobs": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"company_id":  float64(1),
						"title":       "テストエンジニア1",
						"description": "テスト求人の説明1",
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
						"description": "テスト求人の説明2",
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
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "不正なリミットパラメータの場合はデフォルト値を使用する",
			query:          "?page=1&limit=invalid",
			mockResponse:   mockResponse,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"jobs": []interface{}{
					map[string]interface{}{
						"id":          float64(1),
						"company_id":  float64(1),
						"title":       "テストエンジニア1",
						"description": "テスト求人の説明1",
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
						"description": "テスト求人の説明2",
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
				"total": float64(2),
				"page":  float64(1),
				"limit": float64(10),
			},
		},
		{
			name:           "サーバーエラーの場合は500エラーを返す",
			query:          "?page=1&limit=10",
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
			mockUseCase := new(GetJobsMockUseCase)

			// モックの振る舞いを設定
			page := 1
			limit := 10
			mockUseCase.On("GetJobs", mock.Anything, page, limit).
				Return(tt.mockResponse, tt.mockError)

			// テスト用のルーターを作成
			router := gin.New()
			handler := func(c *gin.Context) {
				// クエリパラメータを取得
				pageStr := c.DefaultQuery("page", "1")
				limitStr := c.DefaultQuery("limit", "10")

				// 文字列を整数に変換
				page, err := strconv.Atoi(pageStr)
				if err != nil || page <= 0 {
					page = 1
				}

				limit, err := strconv.Atoi(limitStr)
				if err != nil || limit <= 0 || limit > 100 {
					limit = 10
				}

				// ユースケースを呼び出し
				response, err := mockUseCase.GetJobs(c.Request.Context(), page, limit)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code":    "SERVER_ERROR",
							"message": "サーバーエラーが発生しました",
						},
					})
					return
				}

				// 成功レスポンスを返す
				c.JSON(http.StatusOK, response)
			}
			router.GET("/api/v1/jobs", handler)

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs"+tt.query, nil)
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
				assert.Equal(t, tt.expectedBody, response)
			}

			// モックが期待通り呼ばれたことを検証
			mockUseCase.AssertExpectations(t)
		})
	}
}
