package routes

import (
	"github.com/gin-gonic/gin"
	jobHandler "github.com/takanoakira/ai-interview-practice/backend/internal/handler/job_posting"
	jobUseCase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job_posting"
)

// RegisterJobRoutes は求人関連のルートを登録します
func RegisterJobRoutes(
	router *gin.Engine,
	jobPostingUC jobUseCase.JobPostingUsecase,
) {
	handler := jobHandler.NewJobPostingHandler(jobPostingUC)

	api := router.Group("/api/v1")
	{
		// 求人情報のエンドポイント
		jobs := api.Group("/job-postings")
		{
			jobs.GET("", handler.GetJobPostings)
			jobs.GET("/:id", handler.GetJobPosting)
			jobs.POST("", handler.CreateJobPosting)
			jobs.PUT("/:id", handler.UpdateJobPosting)
			jobs.DELETE("/:id", handler.DeleteJobPosting)
		}

		// 企業に紐づく求人情報のエンドポイント
		companies := api.Group("/companies")
		{
			companies.GET("/:id/job-postings", handler.GetJobPostingsByCompanyID)
			companies.GET("/:id/with-job-postings", handler.GetCompanyWithJobPostings)
			companies.POST("/:id/job-postings", handler.CreateJobPosting)
			companies.PUT("/:id/job-postings/:job_id", handler.UpdateJobPosting)
			companies.DELETE("/:id/job-postings/:job_id", handler.DeleteJobPosting)
		}
	}
}
