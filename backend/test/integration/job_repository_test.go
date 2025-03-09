package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/repository/job"
)

func TestJobRepository_Integration(t *testing.T) {
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
	repo := job.NewJobRepository(db)

	// テスト用のコンテキスト
	ctx := context.Background()

	// テスト用の企業データを作成
	testCompany := &entity.Company{
		Name:                "テスト企業",
		BusinessDescription: stringPtr("テスト企業の説明"),
	}

	err = db.Create(testCompany).Error
	require.NoError(t, err)
	require.NotZero(t, testCompany.ID)

	// テスト用の求人データ
	testJob := &entity.JobPosting{
		CompanyID:           testCompany.ID,
		Title:               "テストエンジニア",
		Description:         "テスト求人の説明",
		Location:            "東京",
		Salary:              "年収500万円〜800万円",
		EmploymentType:      "正社員",
		ApplicationDeadline: time.Now().AddDate(0, 1, 0),
		Status:              "公開中",
		RequiredSkills:      "Go, Docker, MySQL",
		PreferredSkills:     stringPtr("Kubernetes, AWS"),
		Experience:          intPtr(3),
	}

	// 1. Create - 求人の作成をテスト
	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, testJob)
		assert.NoError(t, err)
		assert.NotZero(t, testJob.ID)
		assert.NotZero(t, testJob.CreatedAt)
		assert.NotZero(t, testJob.UpdatedAt)
	})

	// 2. FindByID - 求人の取得をテスト
	t.Run("FindByID", func(t *testing.T) {
		job, err := repo.FindByID(ctx, testJob.ID)
		assert.NoError(t, err)
		assert.NotNil(t, job)
		assert.Equal(t, testJob.ID, job.ID)
		assert.Equal(t, testJob.CompanyID, job.CompanyID)
		assert.Equal(t, testJob.Title, job.Title)
		assert.Equal(t, testJob.Description, job.Description)
		assert.Equal(t, testJob.Location, job.Location)
		assert.Equal(t, testJob.Salary, job.Salary)
		assert.Equal(t, testJob.EmploymentType, job.EmploymentType)
		assert.Equal(t, testJob.Status, job.Status)
		assert.Equal(t, testJob.RequiredSkills, job.RequiredSkills)
		assert.Equal(t, testJob.PreferredSkills, job.PreferredSkills)
		assert.Equal(t, testJob.Experience, job.Experience)
	})

	// 3. FindAll - 求人一覧の取得をテスト
	t.Run("FindAll", func(t *testing.T) {
		// 追加の求人を作成
		for i := 0; i < 5; i++ {
			additionalJob := &entity.JobPosting{
				CompanyID:           testCompany.ID,
				Title:               fmt.Sprintf("追加テスト求人 %d", i+1),
				Description:         fmt.Sprintf("追加テスト求人 %d の説明", i+1),
				Location:            "大阪",
				Salary:              "年収400万円〜600万円",
				EmploymentType:      "契約社員",
				ApplicationDeadline: time.Now().AddDate(0, 1, 0),
				Status:              "公開中",
			}
			err := repo.Create(ctx, additionalJob)
			assert.NoError(t, err)
		}

		// 求人一覧を取得
		jobs, total, err := repo.FindAll(ctx, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total) // 元の1つ + 追加の5つ
		assert.Equal(t, 6, len(jobs))
	})

	// 4. FindByCompanyID - 企業IDによる求人一覧の取得をテスト
	t.Run("FindByCompanyID", func(t *testing.T) {
		// 別の企業を作成
		anotherCompany := &entity.Company{
			Name:                "別のテスト企業",
			BusinessDescription: stringPtr("別のテスト企業の説明"),
		}

		err := db.Create(anotherCompany).Error
		require.NoError(t, err)

		// 別の企業の求人を作成
		anotherJob := &entity.JobPosting{
			CompanyID:           anotherCompany.ID,
			Title:               "別の企業の求人",
			Description:         "別の企業の求人の説明",
			Location:            "福岡",
			Salary:              "年収300万円〜500万円",
			EmploymentType:      "正社員",
			ApplicationDeadline: time.Now().AddDate(0, 1, 0),
			Status:              "公開中",
		}

		err = repo.Create(ctx, anotherJob)
		assert.NoError(t, err)

		// 元の企業の求人一覧を取得
		jobs, total, err := repo.FindByCompanyID(ctx, testCompany.ID, 1, 10)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), total) // 元の企業の求人のみ
		assert.Equal(t, 6, len(jobs))

		// すべての求人が元の企業のものであることを確認
		for _, job := range jobs {
			assert.Equal(t, testCompany.ID, job.CompanyID)
		}
	})

	// 5. Update - 求人の更新をテスト
	t.Run("Update", func(t *testing.T) {
		// 更新用のデータ
		testJob.Title = "更新後のタイトル"
		testJob.Description = "更新後の説明"
		testJob.Location = "名古屋"
		testJob.Salary = "年収600万円〜900万円"
		testJob.RequiredSkills = "Go, Docker, MySQL, Kubernetes"

		// 更新を実行
		err := repo.Update(ctx, testJob)
		assert.NoError(t, err)

		// 更新後のデータを取得して検証
		updatedJob, err := repo.FindByID(ctx, testJob.ID)
		assert.NoError(t, err)
		assert.Equal(t, "更新後のタイトル", updatedJob.Title)
		assert.Equal(t, "更新後の説明", updatedJob.Description)
		assert.Equal(t, "名古屋", updatedJob.Location)
		assert.Equal(t, "年収600万円〜900万円", updatedJob.Salary)
		assert.Equal(t, "Go, Docker, MySQL, Kubernetes", updatedJob.RequiredSkills)
	})

	// 6. Delete - 求人の削除をテスト
	t.Run("Delete", func(t *testing.T) {
		// 削除を実行
		err := repo.Delete(ctx, testJob.ID)
		assert.NoError(t, err)

		// 削除後に取得を試みる
		_, err = repo.FindByID(ctx, testJob.ID)
		assert.Error(t, err) // レコードが見つからないエラーが発生するはず
	})
}
