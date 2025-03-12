package routes

import (
	"github.com/gin-gonic/gin"
	companyHandler "github.com/takanoakira/ai-interview-practice/backend/internal/handler/company"
	companyUseCase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// RegisterCompanyRoutes は企業関連のルートを登録します
func RegisterCompanyRoutes(
	router *gin.Engine,
	getCompaniesUC companyUseCase.GetCompaniesUsecase,
	getCompanyUC companyUseCase.GetCompanyUsecase,
	createCompanyUC companyUseCase.CreateCompanyUsecase,
	updateCompanyUC companyUseCase.UpdateCompanyUsecase,
	deleteCompanyUC companyUseCase.DeleteCompanyUsecase,
) {
	api := router.Group("/api/v1")
	{
		companies := api.Group("/companies")
		{
			companies.GET("", companyHandler.GetCompanies(getCompaniesUC))
			companies.GET("/:id", companyHandler.GetCompany(getCompanyUC))
			companies.POST("", companyHandler.CreateCompany(createCompanyUC))
			companies.PUT("/:id", companyHandler.UpdateCompany(updateCompanyUC))
			companies.DELETE("/:id", companyHandler.DeleteCompany(deleteCompanyUC))
		}
	}
}
