# Debian/Ubuntuベースのイメージを使用
FROM golang:1.21

WORKDIR /api

# MySQLクライアントをインストール
RUN apt-get update && apt-get install -y default-mysql-client

# Go 1.21と互換性のあるairのバージョンをインストール
RUN go install github.com/cosmtrek/air@v1.49.0

# Go 1.21と互換性のあるgolang-migrateのバージョンをインストール
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

# OpenAPIコード生成ツールをインストール
RUN go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.0.0

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# airの設定ファイルをコピー
COPY .air.toml .

EXPOSE 8080

# airを使用してホットリロードを有効化
CMD ["air", "-c", ".air.toml"]