package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// GetJob は指定されたIDの求人情報を取得するハンドラーです
func GetJob(jobUseCase job.JobUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		job, err := jobUseCase.GetJobByID(c.Request.Context(), id)
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
}
