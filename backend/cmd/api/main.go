package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/takanoakira/ai-interview-practice/backend/internal/handler/company"
	"github.com/takanoakira/ai-interview-practice/backend/internal/handler/job_posting"
	companyRepo "github.com/takanoakira/ai-interview-practice/backend/internal/repository/company"
	jobPostingRepo "github.com/takanoakira/ai-interview-practice/backend/internal/repository/job_posting"

	"github.com/takanoakira/ai-interview-practice/backend/internal/routes"
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
	companyRepository := companyRepo.NewRepository(db)
	jobPostingRepository := jobPostingRepo.NewRepository(db)

	// ハンドラーの初期化
	companyHandler := company.NewHandler(companyRepository)
	jobPostingHandler := job_posting.NewHandler(jobPostingRepository)

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
	routes.SetupCompanyRoutes(router, companyHandler)
	routes.SetupJobPostingRoutes(router, jobPostingHandler)

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
