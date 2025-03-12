package job

import (
	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// RegisterRoutes はルーターにハンドラーのルートを登録します
func RegisterRoutes(
	router *gin.Engine,
	createJobUC job.CreateJobUsecase,
	getJobUC job.GetJobUsecase,
	updateJobUC job.UpdateJobUsecase,
	deleteJobUC job.DeleteJobUsecase,
	getJobsUC job.GetJobsUsecase,
	getJobsByCompanyUC job.GetJobsByCompanyIDUsecase,
	getCompanyWithJobsUC job.GetCompanyWithJobsUsecase,
) {
	api := router.Group("/api/v1")
	{
		// 求人情報のエンドポイント
		jobs := api.Group("/jobs")
		{
			jobs.GET("", NewGetJobsHandler(getJobsUC).Handle)
			jobs.GET("/:id", NewGetJobHandler(getJobUC).Handle)
		}

		// 企業に紐づく求人情報のエンドポイント
		companies := api.Group("/companies")
		{
			companies.GET("/:id/jobs", NewGetJobsByCompanyIDHandler(getJobsByCompanyUC).Handle)
			companies.GET("/:id/with-jobs", NewGetCompanyWithJobsHandler(getCompanyWithJobsUC).Handle)
			companies.POST("/:id/jobs", NewCreateJobHandler(createJobUC).Handle)
			companies.PUT("/:id/jobs/:job_id", NewUpdateJobHandler(updateJobUC).Handle)
			companies.DELETE("/:id/jobs/:job_id", NewDeleteJobHandler(deleteJobUC).Handle)
		}
	}
}
