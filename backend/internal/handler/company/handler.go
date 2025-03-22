package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

type Handler interface {
	GetCompanies(c *gin.Context)
	CreateCompany(c *gin.Context)
	UpdateCompany(c *gin.Context)
	DeleteCompany(c *gin.Context)
}

type handler struct {
	usecase company.UseCase
}

func NewHandler(usecase company.UseCase) Handler {
	return &handler{usecase: usecase}
}

func (h *handler) GetCompanies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))

	if limit > 100 {
		limit = 100
	}

	companies, err := h.usecase.GetCompanies(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companies)
}

func (h *handler) CreateCompany(c *gin.Context) {
	var company entity.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateCompany(c.Request.Context(), &company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, company)
}

func (h *handler) UpdateCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company entity.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company.ID = id
	if err := h.usecase.UpdateCompany(c.Request.Context(), &company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新後の企業情報を取得
	companies, err := h.usecase.GetCompanies(c.Request.Context(), 1, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 更新した企業を探す
	var updatedCompany *entity.Company
	for _, comp := range companies.Companies {
		if comp.ID == id {
			updatedCompany = &comp
			break
		}
	}

	if updatedCompany == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Updated company not found"})
		return
	}

	c.JSON(http.StatusOK, updatedCompany)
}

func (h *handler) DeleteCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	if err := h.usecase.DeleteCompany(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
