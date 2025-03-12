package company

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// CreateCompanyRequest は企業情報作成のリクエストを表す構造体です
type CreateCompanyRequest struct {
	Name                string               `json:"name" binding:"required"`
	BusinessDescription *string              `json:"business_description"`
	CustomFields        []CompanyCustomField `json:"custom_fields"`
}

// CompanyCustomField は企業の追加情報を表す構造体です
type CompanyCustomField struct {
	FieldName string `json:"field_name" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// Validate はリクエストのバリデーションを行います
func (r *CreateCompanyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("company name is required")
	}
	for _, field := range r.CustomFields {
		if field.FieldName == "" {
			return errors.New("field name is required for custom fields")
		}
	}
	return nil
}

// CreateCompanyHandler は企業情報作成のハンドラーです
type CreateCompanyHandler struct {
	usecase company.CreateCompanyUsecase
}

// NewCreateCompanyHandler は新しいCreateCompanyHandlerを作成します
func NewCreateCompanyHandler(usecase company.CreateCompanyUsecase) *CreateCompanyHandler {
	return &CreateCompanyHandler{
		usecase: usecase,
	}
}

// Handle は企業情報作成のリクエストを処理します
func (h *CreateCompanyHandler) Handle(c *gin.Context) {
	var req CreateCompanyRequest
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
		Name:                req.Name,
		BusinessDescription: req.BusinessDescription,
		CustomFields:        make([]entity.CompanyCustomField, len(req.CustomFields)),
	}

	for i, field := range req.CustomFields {
		company.CustomFields[i] = entity.CompanyCustomField{
			FieldName: field.FieldName,
			Content:   field.Content,
		}
	}

	if err := h.usecase.Execute(c.Request.Context(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Failed to create company",
				"detail":  err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Company created successfully",
		"data": gin.H{
			"company": company,
		},
	})
}

// CreateCompany は新しい企業情報を作成するハンドラー関数を返します
// 後方互換性のために残しています
func CreateCompany(createCompanyUC company.CreateCompanyUsecase) gin.HandlerFunc {
	handler := NewCreateCompanyHandler(createCompanyUC)
	return func(c *gin.Context) {
		handler.Handle(c)
	}
}
