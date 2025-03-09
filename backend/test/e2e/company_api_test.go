package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	companyHandler "github.com/takanoakira/ai-interview-practice/backend/internal/handler/company"
	companyRepository "github.com/takanoakira/ai-interview-practice/backend/internal/repository/company"
	companyUseCase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// テスト用のデータベース接続を設定
func setupTestDB(t *testing.T) *gorm.DB {
	// 環境変数からテスト用DBの接続情報を取得
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		// デフォルトの接続情報
		dsn = "user:password@tcp(db:3306)/ai_interview_test?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// データベース接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "データベース接続に失敗しました")

	// テスト用のテーブルを作成
	err = db.AutoMigrate(&entity.Company{}, &entity.CompanyCustomField{})
	require.NoError(t, err, "テーブルのマイグレーションに失敗しました")

	return db
}

// テスト前にテーブルをクリーンアップ
func cleanupTables(db *gorm.DB) error {
	if err := db.Exec("DELETE FROM company_custom_fields").Error; err != nil {
		return err
	}
	if err := db.Exec("DELETE FROM companies").Error; err != nil {
		return err
	}
	return nil
}

// テスト用のルーターを設定
func setupRouter(db *gorm.DB) *gin.Engine {
	// Ginのテストモードを設定
	gin.SetMode(gin.TestMode)

	// リポジトリの初期化
	companyRepo := companyRepository.NewCompanyRepository(db)

	// ユースケースの初期化
	companyUC := companyUseCase.NewCompanyUseCase(companyRepo)

	// ルーターの初期化
	router := gin.New()
	router.Use(gin.Recovery())

	// ルートの登録
	companyHandler.RegisterRoutes(router, companyUC)

	return router
}

func TestCompanyAPI_E2E(t *testing.T) {
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
			{
				"field_name": "従業員数",
				"content":    "100人",
			},
		},
	}

	var companyID int

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

	// 2. 企業の取得をテスト
	t.Run("GetCompany", func(t *testing.T) {
		// GETリクエストを作成
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/companies/%d", companyID), nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスボディをパース
		var company map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &company)
		assert.NoError(t, err)

		// 企業データを検証
		assert.Equal(t, float64(companyID), company["id"])
		assert.Equal(t, testCompany["name"], company["name"])
		assert.Equal(t, testCompany["business_description"], company["business_description"])

		// カスタムフィールドを検証
		customFields := company["custom_fields"].([]interface{})
		assert.Equal(t, 2, len(customFields))
	})

	// 3. 企業一覧の取得をテスト
	t.Run("GetCompanies", func(t *testing.T) {
		// 追加の企業を作成
		for i := 0; i < 3; i++ {
			additionalCompany := map[string]interface{}{
				"name":                 fmt.Sprintf("追加E2Eテスト企業 %d", i+1),
				"business_description": fmt.Sprintf("追加E2Eテスト企業 %d の説明", i+1),
			}

			jsonData, err := json.Marshal(additionalCompany)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/companies", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)
		}

		// GETリクエストを作成
		req := httptest.NewRequest(http.MethodGet, "/api/v1/companies?page=1&limit=10", nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// レスポンスボディをパース
		var response map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)

		// 企業一覧を検証
		companies := response["companies"].([]interface{})
		assert.Equal(t, 4, len(companies)) // 元の1つ + 追加の3つ
		assert.Equal(t, float64(4), response["total"])
		assert.Equal(t, float64(1), response["page"])
		assert.Equal(t, float64(10), response["limit"])
	})

	// 4. 企業の更新をテスト
	t.Run("UpdateCompany", func(t *testing.T) {
		// 更新用のデータ
		updateData := map[string]interface{}{
			"name":                 "更新後のE2Eテスト企業",
			"business_description": "更新後のE2Eテスト企業の説明",
			"custom_fields": []map[string]interface{}{
				{
					"field_name": "業界",
					"content":    "更新後のIT業界",
				},
				{
					"field_name": "従業員数",
					"content":    "200人",
				},
			},
		}

		// リクエストボディを作成
		jsonData, err := json.Marshal(updateData)
		require.NoError(t, err)

		// PUTリクエストを作成
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/companies/%d", companyID), bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusOK, rec.Code)

		// 更新後のデータを取得して検証
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/companies/%d", companyID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var updatedCompany map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &updatedCompany)
		assert.NoError(t, err)

		assert.Equal(t, updateData["name"], updatedCompany["name"])
		assert.Equal(t, updateData["business_description"], updatedCompany["business_description"])
	})

	// 5. 企業の削除をテスト
	t.Run("DeleteCompany", func(t *testing.T) {
		// DELETEリクエストを作成
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/companies/%d", companyID), nil)
		rec := httptest.NewRecorder()

		// リクエストを実行
		router.ServeHTTP(rec, req)

		// レスポンスを検証
		assert.Equal(t, http.StatusNoContent, rec.Code)

		// 削除後に取得を試みる
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/companies/%d", companyID), nil)
		rec = httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		// 404エラーが返されることを確認
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
