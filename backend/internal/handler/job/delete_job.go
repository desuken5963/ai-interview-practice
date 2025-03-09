package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// DeleteJob は求人情報を削除するハンドラーです
func DeleteJob(jobUseCase job.JobUseCase) gin.HandlerFunc {
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

		// 求人が指定された企業に属しているか確認
		job, err := jobUseCase.GetJob(c.Request.Context(), jobID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "NOT_FOUND",
					"message": "求人情報が見つかりません",
				},
			})
			return
		}

		if job.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "この企業の求人ではありません",
				},
			})
			return
		}

		// ユースケースを呼び出し
		if err := jobUseCase.DeleteJob(c.Request.Context(), jobID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			})
			return
		}

		// 成功レスポンスを返す
		c.JSON(http.StatusOK, gin.H{
			"message": "求人情報を削除しました",
		})
	}
}
