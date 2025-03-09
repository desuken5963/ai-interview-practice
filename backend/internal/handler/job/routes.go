package job

import (
	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// RegisterRoutes はルーターにハンドラーのルートを登録します
func RegisterRoutes(router *gin.Engine, jobUseCase job.JobUseCase) {
	api := router.Group("/api/v1")
	{
		// 求人情報のエンドポイント
		jobs := api.Group("/jobs")
		{
			jobs.GET("/:id", GetJob(jobUseCase))
		}

		// 企業に紐づく求人情報のエンドポイント
		companies := api.Group("/companies")
		{
			companies.GET("/:id/jobs", GetJobsByCompanyID(jobUseCase))
			companies.GET("/:id/with-jobs", GetCompanyWithJobs(jobUseCase))
			companies.POST("/:id/jobs", CreateJob(jobUseCase))
			companies.PUT("/:id/jobs/:job_id", UpdateJob(jobUseCase))
			companies.DELETE("/:id/jobs/:job_id", DeleteJob(jobUseCase))
		}
	}
}
