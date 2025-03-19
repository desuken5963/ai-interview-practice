# .envファイルを読み込む
-include .env

# 環境変数を直接使用（デフォルト値は.envファイルに記述）
DB_URL = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)

.PHONY: migrate migrate-down migrate-reset db-connect db-show-tables migration-status test-data-company-jobs dev

# 開発環境起動
dev:
	@echo "開発環境を起動します..."
	docker-compose --profile dev up -d

# マイグレーション実行
migrate:
	@echo "マイグレーションを実行します..."
	migrate -database "$(DB_URL)" -path migrations up

# マイグレーションを1つ戻す
migrate-down:
	@echo "マイグレーションを1つ戻します..."
	migrate -database "$(DB_URL)" -path migrations down 1

# マイグレーションを全て戻す
migrate-down-all:
	@echo "マイグレーションを全て戻します..."
	migrate -database "$(DB_URL)" -path migrations down --all

# マイグレーションをリセット（全て戻して再適用）
migrate-reset:
	@echo "マイグレーションをリセットします..."
	migrate -database "$(DB_URL)" -path migrations down --all
	migrate -database "$(DB_URL)" -path migrations up

# 企業求人情報のテストデータ投入
test-data-company-jobs:
	@echo "企業求人情報のテストデータを投入します..."
	mysql -h$(DB_HOST) -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) < test/data/company_jobs.sql

# マイグレーションのステータス確認
migration-status:
	@echo "マイグレーションのステータスを確認します..."
	migrate -database "$(DB_URL)" -path migrations version
	@echo "詳細なステータス:"
	mysql -h$(DB_HOST) -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) -e "SELECT * FROM schema_migrations WHERE version < 1000000;"

# MySQL接続（対話モード）
db-connect:
	@echo "MySQLに接続します..."
	mysql -h$(DB_HOST) -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME)

# テーブル一覧表示
db-show-tables:
	@echo "テーブル一覧を表示します..."
	mysql -h$(DB_HOST) -u$(DB_USER) -p$(DB_PASSWORD) $(DB_NAME) -e "SHOW TABLES;" 