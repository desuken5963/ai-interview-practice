package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// UpdateCompanyHandler は既存の企業情報を更新するハンドラーです
type UpdateCompanyHandler struct {
	Usecase company.CompanyUseCase
}

// NewUpdateCompanyHandler は新しいUpdateCompanyHandlerインスタンスを作成します
func NewUpdateCompanyHandler(usecase company.CompanyUseCase) *UpdateCompanyHandler {
	return &UpdateCompanyHandler{Usecase: usecase}
}

// Handle は企業情報更新リクエストを処理します
func (h *UpdateCompanyHandler) Handle(c *gin.Context) {
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

	var companyData entity.Company

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&companyData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "リクエストの形式が正しくありません",
			},
		})
		return
	}

	// IDを設定
	companyData.ID = id

	// バリデーション
	if companyData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "バリデーションエラーが発生しました",
				"details": []gin.H{
					{
						"field":   "name",
						"message": "企業名は必須です",
					},
				},
			},
		})
		return
	}

	// カスタムフィールドのバリデーション
	for i, field := range companyData.CustomFields {
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
	if err := h.Usecase.UpdateCompany(c.Request.Context(), &companyData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, companyData)
}

// UpdateCompany は既存の企業情報を更新するハンドラー関数を返します
// 後方互換性のために残しています
func UpdateCompany(companyUseCase company.CompanyUseCase) gin.HandlerFunc {
	handler := NewUpdateCompanyHandler(companyUseCase)
	return func(c *gin.Context) {
		handler.Handle(c)
	}
}
