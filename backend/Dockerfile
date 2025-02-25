# GolangベースのDockerイメージ
FROM golang:1.21-alpine

WORKDIR /app

# 必要なパッケージのインストール
RUN apk add --no-cache gcc musl-dev

# Go 1.21と互換性のあるairのバージョンをインストール
RUN go install github.com/cosmtrek/air@v1.49.0

# Go 1.21と互換性のあるgolang-migrateのバージョンをインストール
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

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