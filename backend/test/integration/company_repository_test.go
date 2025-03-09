package integration

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/repository/company"
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

func TestCompanyRepository_Integration(t *testing.T) {
	// 統合テストをスキップするかどうかの環境変数をチェック
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("統合テストをスキップします")
	}

	// テスト用DBのセットアップ
	db := setupTestDB(t)

	// テスト前にテーブルをクリーンアップ
	err := cleanupTables(db)
	require.NoError(t, err, "テーブルのクリーンアップに失敗しました")

	// リポジトリの初期化
	repo := company.NewCompanyRepository(db)

	// テスト用のコンテキスト
	ctx := context.Background()

	// テスト企業データ
	testCompany := &entity.Company{
		Name:                "統合テスト企業",
		BusinessDescription: stringPtr("統合テスト企業の説明"),
		CustomFields: []entity.CompanyCustomField{
			{
				FieldName: "業界",
				Content:   "IT",
			},
			{
				FieldName: "従業員数",
				Content:   "100人",
			},
		},
	}

	// 1. Create - 企業の作成をテスト
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, testCompany)
		assert.NoError(t, err)
		assert.NotZero(t, testCompany.ID)
		assert.NotZero(t, testCompany.CreatedAt)
		assert.NotZero(t, testCompany.UpdatedAt)
		assert.Equal(t, 2, len(testCompany.CustomFields))
		for _, cf := range testCompany.CustomFields {
			assert.NotZero(t, cf.ID)
		}
	})

	// 2. FindByID - 企業の取得をテスト
	t.Run("FindByID", func(t *testing.T) {
		company, err := repo.FindByID(ctx, testCompany.ID)
		assert.NoError(t, err)
		assert.NotNil(t, company)
		assert.Equal(t, testCompany.ID, company.ID)
		assert.Equal(t, testCompany.Name, company.Name)
		assert.Equal(t, testCompany.BusinessDescription, company.BusinessDescription)
		assert.Equal(t, 2, len(company.CustomFields))
	})

	// 3. FindAll - 企業一覧の取得をテスト
	t.Run("FindAll", func(t *testing.T) {
		// 追加の企業を作成
		for i := 0; i < 5; i++ {
			additionalCompany := &entity.Company{
				Name:                fmt.Sprintf("追加テスト企業 %d", i+1),
				BusinessDescription: stringPtr(fmt.Sprintf("追加テスト企業 %d の説明", i+1)),
			}
			err := repo.Create(ctx, additionalCompany)
			assert.NoError(t, err)
		}

		// 企業一覧を取得
		companies, total, err := repo.FindAll(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total) // 元の1つ + 追加の5つ
		assert.Equal(t, 6, len(companies))
	})

	// 4. Update - 企業の更新をテスト
	t.Run("Update", func(t *testing.T) {
		// 更新用のデータ
		testCompany.Name = "更新後の企業名"
		testCompany.BusinessDescription = stringPtr("更新後の説明")
		testCompany.CustomFields[0].Content = "更新後のIT業界"

		// 更新を実行
		err := repo.Update(ctx, testCompany)
		assert.NoError(t, err)

		// 更新後のデータを取得して検証
		updatedCompany, err := repo.FindByID(ctx, testCompany.ID)
		assert.NoError(t, err)
		assert.Equal(t, "更新後の企業名", updatedCompany.Name)
		assert.Equal(t, "更新後の説明", *updatedCompany.BusinessDescription)

		// カスタムフィールドの更新を検証
		found := false
		for _, cf := range updatedCompany.CustomFields {
			if cf.FieldName == "業界" {
				assert.Equal(t, "更新後のIT業界", cf.Content)
				found = true
				break
			}
		}
		assert.True(t, found, "更新されたカスタムフィールドが見つかりませんでした")
	})

	// 5. Delete - 企業の削除をテスト
	t.Run("Delete", func(t *testing.T) {
		// 削除を実行
		err := repo.Delete(ctx, testCompany.ID)
		assert.NoError(t, err)

		// 削除後に取得を試みる
		_, err = repo.FindByID(ctx, testCompany.ID)
		assert.Error(t, err) // レコードが見つからないエラーが発生するはず
	})
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}
