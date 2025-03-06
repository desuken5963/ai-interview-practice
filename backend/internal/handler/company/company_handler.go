package company

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/takanoakira/ai-interview-practice/backend/internal/domain/entity"
	"github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
)

// CompanyHandler は企業情報に関するHTTPリクエストを処理するハンドラーです
type CompanyHandler struct {
	companyUseCase company.CompanyUseCase
}

// NewCompanyHandler は企業ハンドラーの新しいインスタンスを作成します
func NewCompanyHandler(companyUseCase company.CompanyUseCase) *CompanyHandler {
	return &CompanyHandler{
		companyUseCase: companyUseCase,
	}
}

// GetCompanies は企業情報の一覧を取得するハンドラーです
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	// クエリパラメータを取得
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// 文字列を整数に変換
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_PAGE",
				"message": "ページ番号は1以上の整数である必要があります",
			},
		})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_LIMIT",
				"message": "リミットは1から100の間の整数である必要があります",
			},
		})
		return
	}

	// ユースケースを呼び出して企業情報を取得
	response, err := h.companyUseCase.GetCompanies(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "内部サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, response)
}

// GetCompanyByID は指定されたIDの企業情報を取得するハンドラーです
func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	// パスパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "IDは整数である必要があります",
			},
		})
		return
	}

	// ユースケースを呼び出して企業情報を取得
	company, err := h.companyUseCase.GetCompanyByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "COMPANY_NOT_FOUND",
				"message": "指定されたIDの企業が見つかりません",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, company)
}

// CreateCompany は新しい企業情報を作成するハンドラーです
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company entity.Company

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "リクエストボディが不正です",
			},
		})
		return
	}

	// バリデーション
	if company.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "バリデーションエラーが発生しました",
				"details": []gin.H{
					{
						"field":   "name",
						"message": "企業名は必須です",
					},
				},
			},
		})
		return
	}

	// カスタムフィールドのバリデーション
	for i, field := range company.CustomFields {
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

	// ユースケースを呼び出して企業情報を作成
	if err := h.companyUseCase.CreateCompany(c.Request.Context(), &company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "内部サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusCreated, company)
}

// UpdateCompany は既存の企業情報を更新するハンドラーです
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	// パスパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "IDは整数である必要があります",
			},
		})
		return
	}

	var company entity.Company

	// リクエストボディをバインド
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "リクエストボディが不正です",
			},
		})
		return
	}

	// IDを設定
	company.ID = id

	// バリデーション
	if company.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "バリデーションエラーが発生しました",
				"details": []gin.H{
					{
						"field":   "name",
						"message": "企業名は必須です",
					},
				},
			},
		})
		return
	}

	// カスタムフィールドのバリデーション
	for i, field := range company.CustomFields {
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

	// ユースケースを呼び出して企業情報を更新
	if err := h.companyUseCase.UpdateCompany(c.Request.Context(), &company); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "内部サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.JSON(http.StatusOK, company)
}

// DeleteCompany は指定されたIDの企業情報を削除するハンドラーです
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	// パスパラメータからIDを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "IDは整数である必要があります",
			},
		})
		return
	}

	// ユースケースを呼び出して企業情報を削除
	if err := h.companyUseCase.DeleteCompany(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "内部サーバーエラーが発生しました",
			},
		})
		return
	}

	// 成功レスポンスを返す
	c.Status(http.StatusNoContent)
}

// RegisterRoutes はルーターに企業ハンドラーのルートを登録します
func (h *CompanyHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		companies := api.Group("/companies")
		{
			companies.GET("", h.GetCompanies)
			companies.GET("/:id", h.GetCompanyByID)
			companies.POST("", h.CreateCompany)
			companies.PUT("/:id", h.UpdateCompany)
			companies.DELETE("/:id", h.DeleteCompany)
		}
	}
}
