package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// CreateJob は新しい求人情報を作成するハンドラーです
func CreateJob(jobUseCase job.JobUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		idStr := c.Param("id")
		companyID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "IDは整数である必要があります",
				},
			})
			return
		}

		var jobData entity.JobPosting

		// リクエストボディをバインド
		if err := c.ShouldBindJSON(&jobData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "リクエストの形式が正しくありません",
				},
			})
			return
		}

		// 企業IDを設定
		jobData.CompanyID = companyID

		// バリデーション
		if jobData.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "バリデーションエラーが発生しました",
					"details": []gin.H{
						{
							"field":   "title",
							"message": "求人タイトルは必須です",
						},
					},
				},
			})
			return
		}

		// カスタムフィールドのバリデーション
		for i, field := range jobData.CustomFields {
			if field.FieldName == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"code":    "VALIDATION_ERROR",
						"message": "バリデーションエラーが発生しました",
						"details": []gin.H{
							{
								"field":   "custom_fields[" + strconv.Itoa(i) + "].field_name",
								"message": "項目名は必須です",
							},
						},
					},
				})
				return
			}
		}

		// ユースケースを呼び出し
		if err := jobUseCase.CreateJob(c.Request.Context(), &jobData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			})
			return
		}

		// 成功レスポンスを返す
		c.JSON(http.StatusCreated, jobData)
	}
}
