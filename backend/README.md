# AI面接練習アプリケーション バックエンドAPI

## 概要

このバックエンドAPIは、AI面接練習アプリケーションのサーバーサイド機能を提供します。Go言語で実装されており、MySQL 8.0をデータベースとして使用しています。

## 技術スタック

- **言語**: Go 1.21
- **データベース**: MySQL 8.0
- **開発環境**: Docker, Docker Compose
- **ホットリロード**: Air
- **マイグレーション**: golang-migrate

## 開発環境のセットアップ

### 前提条件

- Docker
- Docker Compose

### 環境構築手順

1. リポジトリをクローン
```bash
git clone <リポジトリURL>
cd <リポジトリ名>
```

2. コンテナのビルドと起動
```bash
docker compose build
docker compose up -d
```

3. APIの動作確認
```bash
curl http://localhost:8080/health
```

## データベース操作

コンテナ内で以下のMakefileコマンドを使用できます：

```bash
# APIコンテナに接続
docker compose exec api sh
cd /api

# マイグレーション実行
make migrate

# マイグレーションを1つ戻す
make migrate-down

# マイグレーションを全て戻す
make migrate-down-all

# マイグレーションをリセット（全て戻して再適用）
make migrate-reset

# マイグレーションのステータス確認
make migration-status

# 企業求人情報のテストデータ投入
make test-data-company-jobs

# MySQLに接続（対話モード）
make db-connect

# テーブル一覧表示
make db-show-tables
```

## 環境変数

環境変数は`.env`ファイルで管理されています。主な設定項目：

- `DB_USER`: データベースユーザー名
- `DB_PASSWORD`: データベースパスワード
- `DB_HOST`: データベースホスト名
- `DB_PORT`: データベースポート
- `DB_NAME`: データベース名

## APIエンドポイント

APIエンドポイントの詳細は開発中に追加されます。

## 開発ガイドライン

- コードの変更は自動的にホットリロードされます（Air使用）
- 新しいデータベーススキーマの変更は`migrations`ディレクトリに追加してください
- テストデータは`test_data`ディレクトリ内の適切な名前のSQLファイルに追加し、対応するMakefileコマンドを作成してください
  - 企業求人情報のテストデータは`test_data/company_jobs.sql`に定義されています

## トラブルシューティング

問題が発生した場合は、以下を試してください：

1. コンテナの再起動
```bash
docker compose restart
```

2. コンテナの再ビルド
```bash
docker compose down
docker compose build --no-cache
docker compose up -d
```