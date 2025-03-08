package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// GetCompanies は企業情報の一覧を取得するハンドラーです
func GetCompanies(companyUseCase company.CompanyUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		response, err := companyUseCase.GetCompanies(c.Request.Context(), page, limit)
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
}
