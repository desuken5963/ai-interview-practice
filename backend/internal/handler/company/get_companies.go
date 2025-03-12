package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// GetCompaniesHandler は企業情報の一覧を取得するハンドラーです
type GetCompaniesHandler struct {
	Usecase company.GetCompaniesUsecase
}

// NewGetCompaniesHandler は新しいGetCompaniesHandlerインスタンスを作成します
func NewGetCompaniesHandler(usecase company.GetCompaniesUsecase) *GetCompaniesHandler {
	return &GetCompaniesHandler{Usecase: usecase}
}

// Handle は企業情報一覧取得リクエストを処理します
func (h *GetCompaniesHandler) Handle(c *gin.Context) {
	// クエリパラメータを取得
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// 文字列を整数に変換
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PAGE",
				"message": "ページは1以上の整数である必要があります",
			},
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_LIMIT",
				"message": "リミットは1から100の間の整数である必要があります",
			},
		})
		return
	}

	// ユースケースを呼び出し
	response, err := h.Usecase.Execute(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, response)
}

// GetCompanies は企業情報の一覧を取得するハンドラー関数を返します
// 後方互換性のために残しています
func GetCompanies(getCompaniesUC company.GetCompaniesUsecase) gin.HandlerFunc {
	handler := NewGetCompaniesHandler(getCompaniesUC)
	return func(c *gin.Context) {
		handler.Handle(c)
	}
}
