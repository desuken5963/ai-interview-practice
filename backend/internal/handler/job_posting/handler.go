package job_posting

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job_posting"
)

// JobPostingHandler は求人情報に関するHTTPハンドラーを提供します
type JobPostingHandler struct {
	usecase job_posting.JobPostingUsecase
}

// NewJobPostingHandler は新しいJobPostingHandlerインスタンスを作成します
func NewJobPostingHandler(usecase job_posting.JobPostingUsecase) *JobPostingHandler {
	return &JobPostingHandler{usecase: usecase}
}

// CreateJobPosting は新しい求人情報を作成します
func (h *JobPostingHandler) CreateJobPosting(c *gin.Context) {
	var job entity.JobPosting
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.usecase.Create(c.Request.Context(), &job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, job)
}

// GetJobPostings は求人情報の一覧を取得します
func (h *JobPostingHandler) GetJobPostings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	jobs, err := h.usecase.List(c.Request.Context(), (page-1)*limit, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// GetJobPosting は指定されたIDの求人情報を取得します
func (h *JobPostingHandler) GetJobPosting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	job, err := h.usecase.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job posting not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// UpdateJobPosting は既存の求人情報を更新します
func (h *JobPostingHandler) UpdateJobPosting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var job entity.JobPosting
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	job.ID = id
	if err := h.usecase.Update(c.Request.Context(), &job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

// DeleteJobPosting は指定されたIDの求人情報を削除します
func (h *JobPostingHandler) DeleteJobPosting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetJobPostingsByCompanyID は指定された企業IDの求人情報一覧を取得します
func (h *JobPostingHandler) GetJobPostingsByCompanyID(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	response, err := h.usecase.ListByCompanyID(c.Request.Context(), companyID, (page-1)*limit, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetCompanyWithJobPostings は企業情報と関連する求人情報を取得します
func (h *JobPostingHandler) GetCompanyWithJobPostings(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.usecase.GetCompanyWithJobs(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}
