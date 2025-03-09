package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

// GetJobのテスト用のモック
type GetJobMockUseCase struct {
	mock.Mock
}

func (m *GetJobMockUseCase) GetJob(ctx context.Context, id int) (*entity.JobPosting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobPosting), args.Error(1)
}

func (m *GetJobMockUseCase) GetJobs(ctx context.Context, page, limit int) (*entity.JobResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobResponse), args.Error(1)
}

func (m *GetJobMockUseCase) GetJobsByCompanyID(ctx context.Context, companyID, page, limit int) (*entity.JobResponse, error) {
	args := m.Called(ctx, companyID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobResponse), args.Error(1)
}

func (m *GetJobMockUseCase) CreateJob(ctx context.Context, job *entity.JobPosting) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *GetJobMockUseCase) UpdateJob(ctx context.Context, job *entity.JobPosting) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *GetJobMockUseCase) DeleteJob(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *GetJobMockUseCase) GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.Company), args.Get(1).([]entity.JobPosting), args.Error(2)
}

// GetJobのテスト用のヘルパー関数
func getJobStringPtr(s string) *string {
	return &s
}

func TestGetJob(t *testing.T) {
	// テスト用の求人データ
	now := time.Now()
	mockJob := &entity.JobPosting{
		ID:          1,
		CompanyID:   1,
		Title:       "テストエンジニア",
		Description: getJobStringPtr("テスト求人の説明"),
		CustomFields: []entity.JobCustomField{
			{
				ID:        1,
				JobID:     1,
				FieldName: "勤務地",
				Content:   "東京",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				JobID:     1,
				FieldName: "給与",
				Content:   "年収500万円〜800万円",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        3,
				JobID:     1,
				FieldName: "雇用形態",
				Content:   "正社員",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        4,
				JobID:     1,
				FieldName: "応募締切",
				Content:   now.AddDate(0, 1, 0).Format("2006-01-02"),
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        5,
				JobID:     1,
				FieldName: "ステータス",
				Content:   "公開中",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        6,
				JobID:     1,
				FieldName: "必須スキル",
				Content:   "Go, Docker, MySQL",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        7,
				JobID:     1,
				FieldName: "歓迎スキル",
				Content:   "Kubernetes, AWS",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        8,
				JobID:     1,
				FieldName: "経験年数",
				Content:   "3年以上",
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	// テストケース
	tests := []struct {
		name           string
		id             string
		mockJob        *entity.JobPosting
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "正常に求人を取得できる",
			id:             "1",
			mockJob:        mockJob,
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"company_id":  float64(1),
				"title":       "テストエンジニア",
				"description": "テスト求人の説明",
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
						"content":    now.AddDate(0, 1, 0).Format("2006-01-02"),
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
		},
		{
			name:           "不正なIDパラメータの場合はエラーを返す",
			id:             "invalid",
			mockJob:        nil,
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
			name:           "存在しない求人IDの場合は404エラーを返す",
			id:             "999",
			mockJob:        nil,
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "JOB_NOT_FOUND",
					"message": "指定されたIDの求人が見つかりません",
				},
			},
		},
		{
			name:           "サーバーエラーの場合は500エラーを返す",
			id:             "1",
			mockJob:        nil,
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
			mockUseCase := new(GetJobMockUseCase)

			// 正常なパラメータの場合のみモックの振る舞いを設定
			if tt.id != "invalid" {
				id := 0
				fmt.Sscanf(tt.id, "%d", &id)
				mockUseCase.On("GetJob", mock.Anything, id).
					Return(tt.mockJob, tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			handler := func(c *gin.Context) {
				// パスパラメータからIDを取得
				idStr := c.Param("id")
				id, err := strconv.Atoi(idStr)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": gin.H{
							"code":    "INVALID_ID",
							"message": "IDは整数である必要があります",
						},
					})
					return
				}

				// ユースケースを呼び出し
				job, err := mockUseCase.GetJob(c.Request.Context(), id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": gin.H{
							"code":    "SERVER_ERROR",
							"message": "サーバーエラーが発生しました",
						},
					})
					return
				}

				// 求人が見つからない場合
				if job == nil {
					c.JSON(http.StatusNotFound, gin.H{
						"error": gin.H{
							"code":    "JOB_NOT_FOUND",
							"message": "指定されたIDの求人が見つかりません",
						},
					})
					return
				}

				// 成功レスポンスを返す
				c.JSON(http.StatusOK, job)
			}
			router.GET("/api/v1/jobs/:id", handler)

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/"+tt.id, nil)
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
			if response["created_at"] != nil {
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
				assert.Equal(t, tt.expectedBody, response)
			}

			// モックが期待通り呼ばれたことを検証
			mockUseCase.AssertExpectations(t)
		})
	}
}
