package job

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

func TestDeleteJob(t *testing.T) {
	// テストケース
	tests := []struct {
		name            string
		companyID       string
		jobID           string
		mockJob         *entity.JobPosting
		mockGetError    error
		mockDeleteError error
		expectedStatus  int
		expectedBody    map[string]interface{}
	}{
		{
			name:            "正常に求人を削除できる",
			companyID:       "1",
			jobID:           "1",
			mockJob:         &entity.JobPosting{ID: 1, CompanyID: 1, Title: "テスト求人"},
			mockGetError:    nil,
			mockDeleteError: nil,
			expectedStatus:  http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "求人情報を削除しました",
			},
		},
		{
			name:            "不正な企業IDパラメータの場合はエラーを返す",
			companyID:       "invalid",
			jobID:           "1",
			mockJob:         nil,
			mockGetError:    nil,
			mockDeleteError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_ID",
					"message": "企業IDは整数である必要があります",
				},
			},
		},
		{
			name:            "不正な求人IDパラメータの場合はエラーを返す",
			companyID:       "1",
			jobID:           "invalid",
			mockJob:         nil,
			mockGetError:    nil,
			mockDeleteError: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "INVALID_ID",
					"message": "求人IDは整数である必要があります",
				},
			},
		},
		{
			name:            "求人が見つからない場合は404エラーを返す",
			companyID:       "1",
			jobID:           "999",
			mockJob:         nil,
			mockGetError:    errors.New("求人が見つかりません"),
			mockDeleteError: nil,
			expectedStatus:  http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "NOT_FOUND",
					"message": "求人情報が見つかりません",
				},
			},
		},
		{
			name:            "企業IDと求人の企業IDが一致しない場合は403エラーを返す",
			companyID:       "2",
			jobID:           "1",
			mockJob:         &entity.JobPosting{ID: 1, CompanyID: 1, Title: "テスト求人"},
			mockGetError:    nil,
			mockDeleteError: nil,
			expectedStatus:  http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    "FORBIDDEN",
					"message": "この企業の求人ではありません",
				},
			},
		},
		{
			name:            "サーバーエラーの場合は500エラーを返す",
			companyID:       "1",
			jobID:           "1",
			mockJob:         &entity.JobPosting{ID: 1, CompanyID: 1, Title: "テスト求人"},
			mockGetError:    nil,
			mockDeleteError: errors.New("database error"),
			expectedStatus:  http.StatusInternalServerError,
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

			// GetJobのモック設定
			if tt.companyID != "invalid" && tt.jobID != "invalid" {
				jobID := 0
				fmt.Sscanf(tt.jobID, "%d", &jobID)
				mockUseCase.On("GetJob", mock.Anything, jobID).
					Return(tt.mockJob, tt.mockGetError)
			}

			// DeleteJobのモック設定
			if tt.companyID != "invalid" && tt.jobID != "invalid" && tt.mockJob != nil && tt.mockJob.CompanyID == 1 && tt.companyID == "1" {
				jobID := 0
				fmt.Sscanf(tt.jobID, "%d", &jobID)
				mockUseCase.On("DeleteJob", mock.Anything, jobID).
					Return(tt.mockDeleteError)
			}

			// テスト用のルーターを作成
			router := gin.New()
			router.DELETE("/api/v1/companies/:id/jobs/:job_id", DeleteJob(mockUseCase))

			// テストリクエストを作成
			url := fmt.Sprintf("/api/v1/companies/%s/jobs/%s", tt.companyID, tt.jobID)
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			rec := httptest.NewRecorder()

			// リクエストを実行
			router.ServeHTTP(rec, req)

			// レスポンスを検証
			assert.Equal(t, tt.expectedStatus, rec.Code)

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

			// モックが期待通り呼ばれたことを検証
			mockUseCase.AssertExpectations(t)
		})
	}
}
