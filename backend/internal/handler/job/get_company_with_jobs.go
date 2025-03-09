package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// GetCompanyWithJobs は企業情報と紐づく求人情報を取得するハンドラーです
func GetCompanyWithJobs(jobUseCase job.JobUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		idStr := c.Param("id")
		companyID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "IDは整数である必要があります",
				},
			})
			return
		}

		// ユースケースを呼び出し
		company, jobs, err := jobUseCase.GetCompanyWithJobs(c.Request.Context(), companyID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			})
			return
		}

		// 企業情報と求人情報を含むレスポンスを返す
		c.JSON(http.StatusOK, gin.H{
			"id":                   company.ID,
			"name":                 company.Name,
			"business_description": company.BusinessDescription,
			"custom_fields":        company.CustomFields,
			"job_count":            len(jobs),
			"created_at":           company.CreatedAt,
			"updated_at":           company.UpdatedAt,
			"jobs":                 jobs,
		})
	}
}
