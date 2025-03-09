package job

import (
	"bytes"
	"context"
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

// MockJobUseCase はテスト用のモックです
type MockJobUseCase struct {
	mock.Mock
}

// GetJobs は GetJobs メソッドのモック実装です
func (m *MockJobUseCase) GetJobs(ctx context.Context, page, limit int) (*entity.JobResponse, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobResponse), args.Error(1)
}

// GetJobsByCompanyID は GetJobsByCompanyID メソッドのモック実装です
func (m *MockJobUseCase) GetJobsByCompanyID(ctx context.Context, companyID, page, limit int) (*entity.JobResponse, error) {
	args := m.Called(ctx, companyID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobResponse), args.Error(1)
}

// GetJob は GetJob メソッドのモック実装です
func (m *MockJobUseCase) GetJob(ctx context.Context, id int) (*entity.JobPosting, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.JobPosting), args.Error(1)
}

// CreateJob は CreateJob メソッドのモック実装です
func (m *MockJobUseCase) CreateJob(ctx context.Context, job *entity.JobPosting) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

// UpdateJob は UpdateJob メソッドのモック実装です
func (m *MockJobUseCase) UpdateJob(ctx context.Context, job *entity.JobPosting) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

// DeleteJob は DeleteJob メソッドのモック実装です
func (m *MockJobUseCase) DeleteJob(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetCompanyWithJobs は GetCompanyWithJobs メソッドのモック実装です
func (m *MockJobUseCase) GetCompanyWithJobs(ctx context.Context, companyID int) (*entity.Company, []entity.JobPosting, error) {
	args := m.Called(ctx, companyID)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*entity.Company), args.Get(1).([]entity.JobPosting), args.Error(2)
}

// ヘルパー関数
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestCreateJob(t *testing.T) {
	// テストケース
	tests := []struct {
		name           string
		companyID      string
		requestBody    map[string]interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "正常に求人を作成できる",
			companyID: "1",
			requestBody: map[string]interface{}{
				"title":       "テストエンジニア",
				"description": "テスト求人の説明",
				"custom_fields": []map[string]interface{}{
					{
						"field_name": "勤務地",
						"content":    "東京",
					},
					{
						"field_name": "給与",
						"content":    "年収500万円〜800万円",
					},
					{
						"field_name": "雇用形態",
						"content":    "正社員",
					},
					{
						"field_name": "応募締切",
						"content":    "2025-04-09",
					},
					{
						"field_name": "ステータス",
						"content":    "公開中",
					},
					{
						"field_name": "必須スキル",
						"content":    "Go, Docker, MySQL",
					},
					{
						"field_name": "歓迎スキル",
						"content":    "Kubernetes, AWS",
					},
					{
						"field_name": "経験年数",
						"content":    "3年以上",
					},
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          float64(1),
				"company_id":  float64(1),
				"title":       "テストエンジニア",
				"description": "テスト求人の説明",
				"custom_fields": []interface{}{
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "勤務地",
						"content":    "東京",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "給与",
						"content":    "年収500万円〜800万円",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "雇用形態",
						"content":    "正社員",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "応募締切",
						"content":    "2025-04-09",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "ステータス",
						"content":    "公開中",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "必須スキル",
						"content":    "Go, Docker, MySQL",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "歓迎スキル",
						"content":    "Kubernetes, AWS",
					},
					map[string]interface{}{
						"id":         float64(0),
						"job_id":     float64(0),
						"field_name": "経験年数",
						"content":    "3年以上",
					},
				},
			},
		},
		{
			name:      "不正な企業IDパラメータの場合はエラーを返す",
			companyID: "invalid",
			requestBody: map[string]interface{}{
				"title": "テストエンジニア",
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
			name:      "必須フィールドがない場合はエラーを返す",
			companyID: "1",
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
			if tt.companyID != "invalid" && tt.requestBody["title"] != nil && tt.requestBody["title"] != "" {
				mockUseCase.On("CreateJob", mock.Anything, mock.AnythingOfType("*entity.JobPosting")).
					Run(func(args mock.Arguments) {
						// モックに渡された求人オブジェクトを検証
						job := args.Get(1).(*entity.JobPosting)
						assert.Equal(t, tt.requestBody["title"], job.Title)
						if tt.requestBody["description"] != nil {
							assert.Equal(t, tt.requestBody["description"], *job.Description)
						}

						// IDを設定（作成成功をシミュレート）
						job.ID = 1
					}).
					Return(tt.mockError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.POST("/api/v1/companies/:id/jobs", CreateJob(mockUseCase))

			// リクエストボディをJSON化
			jsonData, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			// テストリクエストを作成
			req := httptest.NewRequest(http.MethodPost, "/api/v1/companies/"+tt.companyID+"/jobs", bytes.NewBuffer(jsonData))
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
