package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// GetCompanyHandler は指定されたIDの企業情報を取得するハンドラーです
type GetCompanyHandler struct {
	Usecase company.GetCompanyUsecase
}

// NewGetCompanyHandler は新しいGetCompanyHandlerインスタンスを作成します
func NewGetCompanyHandler(usecase company.GetCompanyUsecase) *GetCompanyHandler {
	return &GetCompanyHandler{Usecase: usecase}
}

// Handle は企業情報取得リクエストを処理します
func (h *GetCompanyHandler) Handle(c *gin.Context) {
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
	company, err := h.Usecase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	// 企業が見つからない場合
	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "COMPANY_NOT_FOUND",
				"message": "指定されたIDの企業が見つかりません",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, company)
}

// GetCompany は指定されたIDの企業情報を取得するハンドラー関数を返します
// 後方互換性のために残しています
func GetCompany(companyUseCase company.CompanyUseCase) gin.HandlerFunc {
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
		company, err := companyUseCase.GetCompanyByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			})
			return
		}

		// 企業が見つからない場合
		if company == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": gin.H{
					"code":    "COMPANY_NOT_FOUND",
					"message": "指定されたIDの企業が見つかりません",
				},
			})
			return
		}

		// 成功レスポンスを返す
		c.JSON(http.StatusOK, company)
	}
}
