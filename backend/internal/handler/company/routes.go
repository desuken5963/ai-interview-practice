package company

import (
	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// RegisterRoutes はルーターにハンドラーのルートを登録します
func RegisterRoutes(router *gin.Engine, companyUseCase company.CompanyUseCase) {
	api := router.Group("/api/v1")
	{
		companies := api.Group("/companies")
		{
			companies.GET("", GetCompanies(companyUseCase))
			companies.GET("/:id", GetCompany(companyUseCase))
			companies.POST("", CreateCompany(companyUseCase))
			companies.PUT("/:id", UpdateCompany(companyUseCase))
			companies.DELETE("/:id", DeleteCompany(companyUseCase))
		}
	}
}
