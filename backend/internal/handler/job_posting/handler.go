package job_posting

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job_posting"
)

type Handler interface {
	CreateJobPosting(c *gin.Context)
	UpdateJobPosting(c *gin.Context)
	DeleteJobPosting(c *gin.Context)
}

type handler struct {
	usecase job_posting.UseCase
}

func NewHandler(usecase job_posting.UseCase) Handler {
	return &handler{usecase: usecase}
}

type CreateJobPostingRequest struct {
	CompanyID    int                           `json:"company_id" binding:"required"`
	Title        string                        `json:"title" binding:"required,max=100"`
	Description  *string                       `json:"description,omitempty" binding:"omitempty,max=1000"`
	CustomFields []CreateJobCustomFieldRequest `json:"custom_fields,omitempty"`
}

type CreateJobCustomFieldRequest struct {
	FieldName string `json:"field_name" binding:"required,max=50"`
	Content   string `json:"content" binding:"required,max=500"`
}

func (h *handler) CreateJobPosting(c *gin.Context) {
	var req CreateJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobPosting := &entity.JobPosting{
		CompanyID:   req.CompanyID,
		Title:       req.Title,
		Description: req.Description,
	}

	for _, field := range req.CustomFields {
		jobPosting.CustomFields = append(jobPosting.CustomFields, entity.JobCustomField{
			FieldName: field.FieldName,
			Content:   field.Content,
		})
	}

	result, err := h.usecase.CreateJobPosting(c.Request.Context(), jobPosting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *handler) UpdateJobPosting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	var req CreateJobPostingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobPosting := &entity.JobPosting{
		ID:          id,
		CompanyID:   req.CompanyID,
		Title:       req.Title,
		Description: req.Description,
	}

	for _, field := range req.CustomFields {
		jobPosting.CustomFields = append(jobPosting.CustomFields, entity.JobCustomField{
			FieldName: field.FieldName,
			Content:   field.Content,
		})
	}

	result, err := h.usecase.UpdateJobPosting(c.Request.Context(), jobPosting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *handler) DeleteJobPosting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	if err := h.usecase.DeleteJobPosting(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
