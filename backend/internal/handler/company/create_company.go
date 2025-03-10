package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// CreateCompanyHandler は新しい企業情報を作成するハンドラーです
type CreateCompanyHandler struct {
	Usecase company.CompanyUseCase
}

// NewCreateCompanyHandler は新しいCreateCompanyHandlerインスタンスを作成します
func NewCreateCompanyHandler(usecase company.CompanyUseCase) *CreateCompanyHandler {
	return &CreateCompanyHandler{Usecase: usecase}
}

// Handle は企業情報作成リクエストを処理します
func (h *CreateCompanyHandler) Handle(c *gin.Context) {
	var company entity.Company

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "リクエストの形式が正しくありません",
			},
		})
		return
	}

	// バリデーション
	if company.Name == "" {
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
	for i, field := range company.CustomFields {
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
	if err := h.Usecase.CreateCompany(c.Request.Context(), &company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusCreated, company)
}

// CreateCompany は新しい企業情報を作成するハンドラー関数を返します
// 後方互換性のために残しています
func CreateCompany(companyUseCase company.CompanyUseCase) gin.HandlerFunc {
	handler := NewCreateCompanyHandler(companyUseCase)
	return func(c *gin.Context) {
		handler.Handle(c)
	}
}
