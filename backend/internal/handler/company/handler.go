package company

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// CompanyHandler は企業情報に関する全てのハンドラーをまとめた構造体です
type CompanyHandler struct {
	createUsecase company.CreateCompanyUsecase
	getUsecase    company.GetCompanyUsecase
	listUsecase   company.GetCompaniesUsecase
	updateUsecase company.UpdateCompanyUsecase
	deleteUsecase company.DeleteCompanyUsecase
}

// NewCompanyHandler は新しいCompanyHandlerを作成します
func NewCompanyHandler(
	createUsecase company.CreateCompanyUsecase,
	getUsecase company.GetCompanyUsecase,
	listUsecase company.GetCompaniesUsecase,
	updateUsecase company.UpdateCompanyUsecase,
	deleteUsecase company.DeleteCompanyUsecase,
) *CompanyHandler {
	return &CompanyHandler{
		createUsecase: createUsecase,
		getUsecase:    getUsecase,
		listUsecase:   listUsecase,
		updateUsecase: updateUsecase,
		deleteUsecase: deleteUsecase,
	}
}

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

// UpdateCompanyRequest は企業情報更新のリクエストを表す構造体です
type UpdateCompanyRequest struct {
	Name                string               `json:"name" binding:"required"`
	BusinessDescription *string              `json:"business_description"`
	CustomFields        []CompanyCustomField `json:"custom_fields"`
}

// Validate はCreateCompanyRequestのバリデーションを行います
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

// Validate はUpdateCompanyRequestのバリデーションを行います
func (r *UpdateCompanyRequest) Validate() error {
	if r.Name == "" {
		return errors.New("company name is required")
	}
	return nil
}

// Create は新しい企業情報を作成します
func (h *CompanyHandler) Create(c *gin.Context) {
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

	if err := h.createUsecase.Execute(c.Request.Context(), company); err != nil {
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

// Get は指定されたIDの企業情報を取得します
func (h *CompanyHandler) Get(c *gin.Context) {
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

	company, err := h.getUsecase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "COMPANY_NOT_FOUND",
				"message": "指定されたIDの企業が見つかりません",
			},
		})
		return
	}

	c.JSON(http.StatusOK, company)
}

// List は企業情報の一覧を取得します
func (h *CompanyHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

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

	response, err := h.listUsecase.Execute(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update は既存の企業情報を更新します
func (h *CompanyHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "無効な企業IDです",
				"detail":  err.Error(),
				"code":    "INVALID_ID",
			},
		})
		return
	}

	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "リクエストの形式が正しくありません",
				"detail":  err.Error(),
				"code":    "INVALID_REQUEST",
			},
		})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"message": "バリデーションに失敗しました",
				"detail":  err.Error(),
				"code":    "VALIDATION_FAILED",
			},
		})
		return
	}

	company := &entity.Company{
		ID:                  id,
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

	if err := h.updateUsecase.Execute(c.Request.Context(), company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "企業情報の更新に失敗しました",
				"detail":  err.Error(),
				"code":    "UPDATE_FAILED",
			},
		})
		return
	}

	c.JSON(http.StatusOK, company)
}

// Delete は指定されたIDの企業情報を削除します
func (h *CompanyHandler) Delete(c *gin.Context) {
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

	if err := h.deleteUsecase.Execute(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// 以下は後方互換性のためのヘルパー関数です

func CreateCompany(createCompanyUC company.CreateCompanyUsecase) gin.HandlerFunc {
	handler := NewCompanyHandler(createCompanyUC, nil, nil, nil, nil)
	return handler.Create
}

func GetCompany(getCompanyUC company.GetCompanyUsecase) gin.HandlerFunc {
	handler := NewCompanyHandler(nil, getCompanyUC, nil, nil, nil)
	return handler.Get
}

func GetCompanies(getCompaniesUC company.GetCompaniesUsecase) gin.HandlerFunc {
	handler := NewCompanyHandler(nil, nil, getCompaniesUC, nil, nil)
	return handler.List
}

func UpdateCompany(updateCompanyUC company.UpdateCompanyUsecase) gin.HandlerFunc {
	handler := NewCompanyHandler(nil, nil, nil, updateCompanyUC, nil)
	return handler.Update
}

func DeleteCompany(deleteCompanyUC company.DeleteCompanyUsecase) gin.HandlerFunc {
	handler := NewCompanyHandler(nil, nil, nil, nil, deleteCompanyUC)
	return handler.Delete
}
