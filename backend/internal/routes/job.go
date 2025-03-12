package routes

import (
	"github.com/gin-gonic/gin"
	jobHandler "github.com/takanoakira/ai-interview-practice/backend/internal/handler/job"
	jobUseCase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// RegisterJobRoutes は求人関連のルートを登録します
func RegisterJobRoutes(
	router *gin.Engine,
	createJobUC jobUseCase.CreateJobUsecase,
	getJobUC jobUseCase.GetJobUsecase,
	updateJobUC jobUseCase.UpdateJobUsecase,
	deleteJobUC jobUseCase.DeleteJobUsecase,
	getJobsUC jobUseCase.GetJobsUsecase,
	getJobsByCompanyUC jobUseCase.GetJobsByCompanyIDUsecase,
	getCompanyWithJobsUC jobUseCase.GetCompanyWithJobsUsecase,
) {
	api := router.Group("/api/v1")
	{
		// 求人情報のエンドポイント
		jobs := api.Group("/jobs")
		{
			jobs.GET("", jobHandler.NewGetJobsHandler(getJobsUC).Handle)
			jobs.GET("/:id", jobHandler.NewGetJobHandler(getJobUC).Handle)
		}

		// 企業に紐づく求人情報のエンドポイント
		companies := api.Group("/companies")
		{
			companies.GET("/:id/jobs", jobHandler.NewGetJobsByCompanyIDHandler(getJobsByCompanyUC).Handle)
			companies.GET("/:id/with-jobs", jobHandler.NewGetCompanyWithJobsHandler(getCompanyWithJobsUC).Handle)
			companies.POST("/:id/jobs", jobHandler.NewCreateJobHandler(createJobUC).Handle)
			companies.PUT("/:id/jobs/:job_id", jobHandler.NewUpdateJobHandler(updateJobUC).Handle)
			companies.DELETE("/:id/jobs/:job_id", jobHandler.NewDeleteJobHandler(deleteJobUC).Handle)
		}
	}
}
