package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/takanoakira/ai-interview-practice/backend/internal/handler/company"
)

func SetupCompanyRoutes(r *gin.Engine, h company.Handler) {
	companies := r.Group("/api/v1/companies")
	{
		companies.GET("", h.GetCompanies)
		companies.GET("/:id", h.GetCompanyByID)
		companies.POST("", h.CreateCompany)
		companies.PUT("/:id", h.UpdateCompany)
		companies.DELETE("/:id", h.DeleteCompany)
	}
}
