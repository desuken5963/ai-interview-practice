package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobAPI_E2E(t *testing.T) {
	// E2Eテストをスキップするかどうかの環境変数をチェック
	if os.Getenv("SKIP_E2E_TESTS") == "true" {
		t.Skip("E2Eテストをスキップします")
	}

	// テスト用DBのセットアップ
	db := setupTestDB(t)

	// テスト前にテーブルをクリーンアップ
	err := cleanupTables(db)
	require.NoError(t, err, "テーブルのクリーンアップに失敗しました")

	// テスト用のルーターを設定
	router := setupRouter(db)

	// テスト用の企業データ
	testCompany := map[string]interface{}{
		"name":                 "E2Eテスト企業",
		"business_description": "E2Eテスト企業の説明",
		"custom_fields": []map[string]interface{}{
			{
				"field_name": "業界",
				"content":    "IT",
			},
		},
	}

	var companyID int
	var jobPostingID int

	// 1. 企業の作成をテスト
	t.Run("CreateCompany", func(t *testing.T) {
		// リクエストボディを作成
		jsonData, err := json.Marshal(testCompany)
		require.NoError(t, err)

		// POSTリクエストを作成
		req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusCreated, rec.Code)

		// レスポンスボディをパース
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 企業IDを取得
		companyID = int(response["id"].(float64))
		assert.NotZero(t, companyID)
	})

	// 2. 求人の作成をテスト
	t.Run("CreateJobPosting", func(t *testing.T) {
		// テスト用の求人データ
		testJobPosting := map[string]interface{}{
			"company_id":           companyID,
			"title":                "E2Eテストエンジニア",
			"description":          "E2Eテスト求人の説明",
			"location":             "東京",
			"salary":               "年収500万円〜800万円",
			"employment_type":      "正社員",
			"application_deadline": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
			"status":               "公開中",
			"required_skills":      "Go, Docker, MySQL",
			"preferred_skills":     "Kubernetes, AWS",
			"experience":           3,
		}

		// リクエストボディを作成
		jsonData, err := json.Marshal(testJobPosting)
		require.NoError(t, err)

		// POSTリクエストを作成
		req := httptest.NewRequest(http.MethodPost, "/api/v1/job-postings", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusCreated, rec.Code)

		// レスポンスボディをパース
		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 求人IDを取得
		jobPostingID = int(response["id"].(float64))
		assert.NotZero(t, jobPostingID)
	})

	// 3. 求人の取得をテスト
	t.Run("GetJobPosting", func(t *testing.T) {
		// GETリクエストを作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/job-postings/%d", jobPostingID), nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスボディをパース
		var jobPosting map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &jobPosting)
		assert.NoError(t, err)

		// 求人データを検証
		assert.Equal(t, float64(jobPostingID), jobPosting["id"])
		assert.Equal(t, float64(companyID), jobPosting["company_id"])
		assert.Equal(t, "E2Eテストエンジニア", jobPosting["title"])
		assert.Equal(t, "E2Eテスト求人の説明", jobPosting["description"])
		assert.Equal(t, "東京", jobPosting["location"])
		assert.Equal(t, "年収500万円〜800万円", jobPosting["salary"])
		assert.Equal(t, "正社員", jobPosting["employment_type"])
		assert.Equal(t, "公開中", jobPosting["status"])
		assert.Equal(t, "Go, Docker, MySQL", jobPosting["required_skills"])
		assert.Equal(t, "Kubernetes, AWS", jobPosting["preferred_skills"])
		assert.Equal(t, float64(3), jobPosting["experience"])
	})

	// 4. 求人一覧の取得をテスト
	t.Run("GetJobPostings", func(t *testing.T) {
		// 追加の求人を作成
		for i := 0; i < 3; i++ {
			additionalJobPosting := map[string]interface{}{
				"company_id":           companyID,
				"title":                fmt.Sprintf("追加E2Eテスト求人 %d", i+1),
				"description":          fmt.Sprintf("追加E2Eテスト求人 %d の説明", i+1),
				"location":             "大阪",
				"salary":               "年収400万円〜600万円",
				"employment_type":      "契約社員",
				"application_deadline": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
				"status":               "公開中",
			}

			jsonData, err := json.Marshal(additionalJobPosting)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/job-postings", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// GETリクエストを作成
		req := httptest.NewRequest(http.MethodGet, "/api/v1/job-postings?page=1&limit=10", nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスボディをパース
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 求人一覧を検証
		jobPostings := response["jobPostings"].([]interface{})
		assert.Equal(t, 4, len(jobPostings)) // 元の1つ + 追加の3つ
		assert.Equal(t, float64(4), response["total"])
		assert.Equal(t, float64(1), response["page"])
		assert.Equal(t, float64(10), response["limit"])
	})

	// 5. 企業IDによる求人一覧の取得をテスト
	t.Run("GetJobPostingsByCompanyID", func(t *testing.T) {
		// GETリクエストを作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/companies/%d/job-postings?page=1&limit=10", companyID), nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスボディをパース
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 求人一覧を検証
		jobPostings := response["jobPostings"].([]interface{})
		assert.Equal(t, 4, len(jobPostings)) // 元の1つ + 追加の3つ
		assert.Equal(t, float64(4), response["total"])
		assert.Equal(t, float64(1), response["page"])
		assert.Equal(t, float64(10), response["limit"])

		// すべての求人が同じ企業IDであることを確認
		for _, jobPosting := range jobPostings {
			jobPostingMap := jobPosting.(map[string]interface{})
			assert.Equal(t, float64(companyID), jobPostingMap["company_id"])
		}
	})

	// 6. 求人の更新をテスト
	t.Run("UpdateJobPosting", func(t *testing.T) {
		// 更新用のデータ
		updateData := map[string]interface{}{
			"company_id":           companyID,
			"title":                "更新後のE2Eテストエンジニア",
			"description":          "更新後のE2Eテスト求人の説明",
			"location":             "名古屋",
			"salary":               "年収600万円〜900万円",
			"employment_type":      "正社員",
			"application_deadline": time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
			"status":               "公開中",
			"required_skills":      "Go, Docker, MySQL, Kubernetes",
			"preferred_skills":     "AWS, GCP",
			"experience":           5,
		}

		// リクエストボディを作成
		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// PUTリクエストを作成
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/job-postings/%d", jobPostingID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// 更新後のデータを取得して検証
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/job-postings/%d", jobPostingID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var updatedJobPosting map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &updatedJobPosting)
		assert.NoError(t, err)

		assert.Equal(t, updateData["title"], updatedJobPosting["title"])
		assert.Equal(t, updateData["description"], updatedJobPosting["description"])
		assert.Equal(t, updateData["location"], updatedJobPosting["location"])
		assert.Equal(t, updateData["salary"], updatedJobPosting["salary"])
		assert.Equal(t, updateData["required_skills"], updatedJobPosting["required_skills"])
		assert.Equal(t, updateData["preferred_skills"], updatedJobPosting["preferred_skills"])
		assert.Equal(t, updateData["experience"], updatedJobPosting["experience"])
	})

	// 7. 求人の削除をテスト
	t.Run("DeleteJobPosting", func(t *testing.T) {
		// DELETEリクエストを作成
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/job-postings/%d", jobPostingID), nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// 削除後に取得を試みる
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/job-postings/%d", jobPostingID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// 404エラーが返されることを確認
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
