package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/takanoakira/ai-interview-practice/backend/internal/handler/company"
	"github.com/takanoakira/ai-interview-practice/backend/internal/handler/job"
	companyRepo "github.com/takanoakira/ai-interview-practice/backend/internal/repository/company"
	jobRepo "github.com/takanoakira/ai-interview-practice/backend/internal/repository/job"
	companyUsecase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
	jobUsecase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
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
	// 企業関連のリポジトリ
	getCompaniesRepo := companyRepo.NewGetCompaniesRepository(db)
	getCompanyRepo := companyRepo.NewGetCompanyRepository(db)
	createCompanyRepo := companyRepo.NewCreateCompanyRepository(db)
	updateCompanyRepo := companyRepo.NewUpdateCompanyRepository(db)
	deleteCompanyRepo := companyRepo.NewDeleteCompanyRepository(db)

	// 求人関連のリポジトリ
	getJobsRepo := jobRepo.NewGetJobsRepository(db)
	getJobRepo := jobRepo.NewGetJobRepository(db)
	createJobRepo := jobRepo.NewCreateJobRepository(db)
	updateJobRepo := jobRepo.NewUpdateJobRepository(db)
	deleteJobRepo := jobRepo.NewDeleteJobRepository(db)
	getJobsByCompanyRepo := jobRepo.NewGetJobsByCompanyRepository(db)

	// ユースケースの初期化
	// 企業関連のユースケース
	companyUC := companyUsecase.NewCompanyUseCase(
		createCompanyRepo,
		updateCompanyRepo,
		deleteCompanyRepo,
		getCompanyRepo,
		getCompaniesRepo,
	)

	// 求人関連のユースケース
	getJobsUC := jobUsecase.NewGetJobsUsecase(getJobsRepo)
	getJobUC := jobUsecase.NewGetJobUsecase(getJobRepo)
	createJobUC := jobUsecase.NewCreateJobUsecase(createJobRepo)
	updateJobUC := jobUsecase.NewUpdateJobUsecase(updateJobRepo)
	deleteJobUC := jobUsecase.NewDeleteJobUsecase(deleteJobRepo)
	getJobsByCompanyIDUC := jobUsecase.NewGetJobsByCompanyIDUsecase(getJobsByCompanyRepo)
	getCompanyWithJobsUC := jobUsecase.NewGetCompanyWithJobsUsecase(getJobsByCompanyRepo, getCompanyRepo)

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
	company.RegisterRoutes(router, companyUC)

	// 求人ハンドラーの登録
	job.RegisterRoutes(
		router,
		createJobUC,
		getJobUC,
		updateJobUC,
		deleteJobUC,
		getJobsUC,
		getJobsByCompanyIDUC,
		getCompanyWithJobsUC,
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
