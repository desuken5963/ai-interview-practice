# .envファイルを読み込む
-include .env

# 環境変数を直接使用（デフォルト値は.envファイルに記述）
DB_URL = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)
TEST_DB_URL = mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)_test

.PHONY: migrate migrate-down migrate-reset db-connect db-show-tables migration-status test-data-company-jobs test test-unit test-integration test-e2e test-setup test-docker-up test-docker-down test-handler test-usecase test-repository test-domain dev

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

# テスト用Docker環境の起動
test-docker-up:
	@echo "テスト用Docker環境を起動します..."
	docker-compose --profile test up -d

# テスト用Docker環境の停止
test-docker-down:
	@echo "テスト用Docker環境を停止します..."
	docker-compose --profile test down

# テスト用データベースのセットアップ
test-setup: test-docker-up
	@echo "テスト用データベースをセットアップします..."
	docker-compose exec api_test bash -c "cd /api && migrate -database \"mysql://user:password@tcp(db_test:3306)/ai_interview_test?charset=utf8mb4&parseTime=True&loc=Local\" -path migrations up"

# ハンドラー層の単体テスト実行
test-handler: test-docker-up
	@echo "ハンドラー層の単体テストを実行します..."
	docker-compose exec api_test bash -c "cd /api && go test -v ./internal/handler/..."

# ユースケース層の単体テスト実行
test-usecase: test-docker-up
	@echo "ユースケース層の単体テストを実行します..."
	docker-compose exec api_test bash -c "cd /api && go test -v ./internal/usecase/..."

# リポジトリ層の単体テスト実行
test-repository: test-docker-up
	@echo "リポジトリ層の単体テストを実行します..."
	docker-compose exec api_test bash -c "cd /api && go test -v ./internal/infrastructure/..."

# ドメイン層の単体テスト実行
test-domain: test-docker-up
	@echo "ドメイン層の単体テストを実行します..."
	docker-compose exec api_test bash -c "cd /api && go test -v ./internal/domain/..."

# 単体テスト実行（全レイヤー）
test-unit: test-handler test-usecase test-repository test-domain
	@echo "全ての単体テストが完了しました"

# 統合テスト実行
test-integration: test-docker-up
	@echo "統合テストを実行します..."
	docker-compose exec api_test bash -c "cd /api && go test -v ./test/integration/..."

# E2Eテスト実行
test-e2e: test-docker-up
	@echo "E2Eテストを実行します..."
	docker-compose exec api_test bash -c "cd /api && go test -v ./test/e2e/..."

# 全テスト実行
test: test-docker-up test-unit test-integration test-e2e
	@echo "全てのテストが完了しました"
	docker-compose --profile test down 