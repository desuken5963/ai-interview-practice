package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// UpdateJob は既存の求人情報を更新するハンドラーです
func UpdateJob(jobUseCase job.JobUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータから企業IDを取得
		companyIDStr := c.Param("id")
		companyID, err := strconv.Atoi(companyIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "企業IDは整数である必要があります",
				},
			})
			return
		}

		// パスパラメータから求人IDを取得
		jobIDStr := c.Param("job_id")
		jobID, err := strconv.Atoi(jobIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "求人IDは整数である必要があります",
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

		// IDを設定
		jobData.ID = jobID
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
		if err := jobUseCase.UpdateJob(c.Request.Context(), &jobData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			})
			return
		}

		// 成功レスポンスを返す
		c.JSON(http.StatusOK, jobData)
	}
}
