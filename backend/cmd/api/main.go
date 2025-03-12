package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	companyRepo "github.com/takanoakira/ai-interview-practice/backend/internal/repository/company"
	jobRepo "github.com/takanoakira/ai-interview-practice/backend/internal/repository/job_posting"
	"github.com/takanoakira/ai-interview-practice/backend/internal/routes"
	companyUsecase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
	jobUsecase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job_posting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// データベース接続
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	// 環境変数が設定されていない場合はデフォルト値を使用
	if dsn == "@tcp(:)/?charset=utf8mb4&parseTime=True&loc=Local" {
		dsn = "user:password@tcp(db:3306)/ai_interview?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// リポジトリの初期化
	companyRepository := companyRepo.NewCompanyRepository(db)
	jobPostingRepository := jobRepo.NewJobPostingRepository(db)

	// ユースケースの初期化
	// 企業関連のユースケース
	getCompaniesUC := companyUsecase.NewGetCompaniesUsecase(companyRepository)
	getCompanyUC := companyUsecase.NewGetCompanyUsecase(companyRepository)
	createCompanyUC := companyUsecase.NewCreateCompanyUsecase(companyRepository)
	updateCompanyUC := companyUsecase.NewUpdateCompanyUsecase(companyRepository)
	deleteCompanyUC := companyUsecase.NewDeleteCompanyUsecase(companyRepository)

	// 求人関連のユースケース
	jobPostingUC := jobUsecase.NewJobPostingUsecase(jobPostingRepository)

	// ルーターの設定
	router := gin.Default()

	// CORSの設定
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ハンドラーの登録
	routes.RegisterCompanyRoutes(
		router,
		getCompaniesUC,
		getCompanyUC,
		createCompanyUC,
		updateCompanyUC,
		deleteCompanyUC,
	)

	// 求人ハンドラーの登録
	routes.RegisterJobRoutes(
		router,
		jobPostingUC,
	)

	// サーバーの起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started at :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
