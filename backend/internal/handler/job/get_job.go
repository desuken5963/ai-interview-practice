package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// GetJobHandler は求人情報を取得するハンドラーです
type GetJobHandler struct {
	Usecase job.GetJobUsecase
}

// NewGetJobHandler は新しいGetJobHandlerインスタンスを作成します
func NewGetJobHandler(usecase job.GetJobUsecase) *GetJobHandler {
	return &GetJobHandler{Usecase: usecase}
}

// Handle は求人情報取得リクエストを処理します
func (h *GetJobHandler) Handle(c *gin.Context) {
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
	job, err := h.Usecase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "求人情報が見つかりませんでした",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, job)
}
