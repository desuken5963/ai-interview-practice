package company

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// UpdateCompanyRequest は企業情報更新のリクエストを表す構造体です
type UpdateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// Validate はリクエストのバリデーションを行います
func (r *UpdateCompanyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("company name is required")
	}
	return nil
}

// UpdateCompanyHandler は企業情報更新のハンドラーです
type UpdateCompanyHandler struct {
	usecase company.UpdateCompanyUsecase
}

// NewUpdateCompanyHandler は新しいUpdateCompanyHandlerを作成します
func NewUpdateCompanyHandler(usecase company.UpdateCompanyUsecase) *UpdateCompanyHandler {
	return &UpdateCompanyHandler{
		usecase: usecase,
	}
}

// Handle は企業情報更新のリクエストを処理します
func (h *UpdateCompanyHandler) Handle(c *gin.Context) {
	// パスパラメータからIDを取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid company ID",
				"detail":  err.Error(),
			},
		})
		return
	}

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Invalid request format",
				"detail":  err.Error(),
			},
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "Validation failed",
				"detail":  err.Error(),
			},
		})
		return
	}

	company := &entity.Company{
		ID:                  id,
		Name:                req.Name,
		BusinessDescription: &req.Description,
	}

	if err := h.usecase.Execute(c.Request.Context(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to update company",
				"detail":  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Company updated successfully",
		"data": gin.H{
			"company": company,
		},
	})
}

// UpdateCompany は既存の企業情報を更新するハンドラー関数を返します
// 後方互換性のために残しています
func UpdateCompany(updateCompanyUC company.UpdateCompanyUsecase) gin.HandlerFunc {
	handler := NewUpdateCompanyHandler(updateCompanyUC)
	return func(c *gin.Context) {
		handler.Handle(c)
	}
}
