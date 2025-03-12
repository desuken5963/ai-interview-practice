package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// DeleteCompanyHandler は企業情報を削除するハンドラーです
type DeleteCompanyHandler struct {
	Usecase company.DeleteCompanyUsecase
}

// NewDeleteCompanyHandler は新しいDeleteCompanyHandlerインスタンスを作成します
func NewDeleteCompanyHandler(usecase company.DeleteCompanyUsecase) *DeleteCompanyHandler {
	return &DeleteCompanyHandler{Usecase: usecase}
}

// Handle は企業情報削除リクエストを処理します
func (h *DeleteCompanyHandler) Handle(c *gin.Context) {
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
	if err := h.Usecase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.Status(http.StatusNoContent)
}

// DeleteCompany は企業情報を削除するハンドラー関数を返します
// 後方互換性のために残しています
func DeleteCompany(deleteCompanyUC company.DeleteCompanyUsecase) gin.HandlerFunc {
	handler := NewDeleteCompanyHandler(deleteCompanyUC)
	return func(c *gin.Context) {
		handler.Handle(c)
	}
}
