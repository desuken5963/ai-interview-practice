package job

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
)

// UpdateJobHandler は求人情報を更新するハンドラーです
type UpdateJobHandler struct {
	Usecase job.UpdateJobUsecase
}

// NewUpdateJobHandler は新しいUpdateJobHandlerインスタンスを作成します
func NewUpdateJobHandler(usecase job.UpdateJobUsecase) *UpdateJobHandler {
	return &UpdateJobHandler{Usecase: usecase}
}

// Handle は求人情報更新リクエストを処理します
func (h *UpdateJobHandler) Handle(c *gin.Context) {
	// パスパラメータからIDを取得
	companyIDStr := c.Param("id")
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "企業IDは整数である必要があります",
			},
		})
		return
	}

	jobIDStr := c.Param("job_id")
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "求人IDは整数である必要があります",
			},
		})
		return
	}

	var jobData entity.JobPosting

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&jobData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "リクエストの形式が正しくありません",
			},
		})
		return
	}

	// IDと企業IDを設定
	jobData.ID = jobID
	jobData.CompanyID = companyID

	// バリデーション
	if jobData.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "バリデーションエラーが発生しました",
				"details": []gin.H{
					{
						"field":   "title",
						"message": "求人タイトルは必須です",
					},
				},
			},
		})
		return
	}

	// カスタムフィールドのバリデーション
	for i, field := range jobData.CustomFields {
		if field.FieldName == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "バリデーションエラーが発生しました",
					"details": []gin.H{
						{
							"field":   "custom_fields[" + strconv.Itoa(i) + "].field_name",
							"message": "項目名は必須です",
						},
					},
				},
			})
			return
		}
	}

	// ユースケースを呼び出し
	if err := h.Usecase.Execute(c.Request.Context(), &jobData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "SERVER_ERROR",
				"message": "サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, jobData)
}

// UpdateJob は求人情報を更新するハンドラー関数を返します
// 後方互換性のために残しています
func UpdateJob(usecase job.UpdateJobUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// パスパラメータからIDを取得
		companyIDStr := c.Param("id")
		companyID, err := strconv.Atoi(companyIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "企業IDは整数である必要があります",
				},
			})
			return
		}

		jobIDStr := c.Param("job_id")
		jobID, err := strconv.Atoi(jobIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_ID",
					"message": "求人IDは整数である必要があります",
				},
			})
			return
		}

		var jobData entity.JobPosting

		// リクエストボディをバインド
		if err := c.ShouldBindJSON(&jobData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "INVALID_REQUEST",
					"message": "リクエストの形式が正しくありません",
				},
			})
			return
		}

		// IDと企業IDを設定
		jobData.ID = jobID
		jobData.CompanyID = companyID

		// バリデーション
		if jobData.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "バリデーションエラーが発生しました",
					"details": []gin.H{
						{
							"field":   "title",
							"message": "求人タイトルは必須です",
						},
					},
				},
			})
			return
		}

		// カスタムフィールドのバリデーション
		for i, field := range jobData.CustomFields {
			if field.FieldName == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"code":    "VALIDATION_ERROR",
						"message": "バリデーションエラーが発生しました",
						"details": []gin.H{
							{
								"field":   "custom_fields[" + strconv.Itoa(i) + "].field_name",
								"message": "項目名は必須です",
							},
						},
					},
				})
				return
			}
		}

		// ユースケースを呼び出し
		if err := usecase.Execute(c.Request.Context(), &jobData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": gin.H{
					"code":    "SERVER_ERROR",
					"message": "サーバーエラーが発生しました",
				},
			})
			return
		}

		// 成功レスポンスを返す
		c.JSON(http.StatusOK, jobData)
	}
}
