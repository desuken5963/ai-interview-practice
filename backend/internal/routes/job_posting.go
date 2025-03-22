package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/takanoakira/ai-interview-practice/backend/internal/handler/job_posting"
)

func SetupJobPostingRoutes(r *gin.Engine, h job_posting.Handler) {
	jobPostings := r.Group("/api/v1/job-postings")
	{
		jobPostings.POST("", h.CreateJobPosting)
		jobPostings.PUT("/:id", h.UpdateJobPosting)
		jobPostings.DELETE("/:id", h.DeleteJobPosting)
	}
}
