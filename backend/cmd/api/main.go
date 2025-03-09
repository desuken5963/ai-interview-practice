package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	companyHandler "github.com/takanoakira/ai-interview-practice/backend/internal/handler/company"
	jobHandler "github.com/takanoakira/ai-interview-practice/backend/internal/handler/job"
	companyRepository "github.com/takanoakira/ai-interview-practice/backend/internal/repository/company"
	jobRepository "github.com/takanoakira/ai-interview-practice/backend/internal/repository/job"
	companyUseCase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/company"
	jobUseCase "github.com/takanoakira/ai-interview-practice/backend/internal/usecase/job"
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
	companyRepo := companyRepository.NewCompanyRepository(db)
	jobRepo := jobRepository.NewJobRepository(db)

	// ユースケースの初期化
	companyUC := companyUseCase.NewCompanyUseCase(companyRepo)
	jobUC := jobUseCase.NewJobUseCase(jobRepo, companyRepo)

	// ルーターの初期化
	router := gin.Default()

	// CORSの設定
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// ルートの登録
	companyHandler.RegisterRoutes(router, companyUC)
	jobHandler.RegisterRoutes(router, jobUC)

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
