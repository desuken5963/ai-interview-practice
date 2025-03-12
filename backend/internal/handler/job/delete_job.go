package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// DeleteJobHandler は求人情報を削除するハンドラーです
type DeleteJobHandler struct {
	Usecase job.DeleteJobUsecase
}

// NewDeleteJobHandler は新しいDeleteJobHandlerインスタンスを作成します
func NewDeleteJobHandler(usecase job.DeleteJobUsecase) *DeleteJobHandler {
	return &DeleteJobHandler{Usecase: usecase}
}

// Handle は求人情報削除リクエストを処理します
func (h *DeleteJobHandler) Handle(c *gin.Context) {
	// パスパラメータからIDを取得
	companyIDStr := c.Param("id")
	_, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "企業IDは整数である必要があります",
			},
		})
		return
	}

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

	// ユースケースを呼び出し
	if err := h.Usecase.Execute(c.Request.Context(), jobID); err != nil {
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
		"message": "求人情報が正常に削除されました",
	})
}

// DeleteJob は求人情報を削除するハンドラー関数を返します
// 後方互換性のために残しています
func DeleteJob(usecase job.DeleteJobUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		companyIDStr := c.Param("id")
		_, err := strconv.Atoi(companyIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "企業IDは整数である必要があります",
				},
			})
			return
		}

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

		// ユースケースを呼び出し
		if err := usecase.Execute(c.Request.Context(), jobID); err != nil {
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
			"message": "求人情報が正常に削除されました",
		})
	}
}
