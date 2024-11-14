# ベースイメージとしてGo 1.23を使用（ビルド用ステージ）
FROM golang:1.23 as build

# 作業ディレクトリを設定
WORKDIR /app

# Goモジュールファイルをコピー
COPY go.mod ./

# 依存関係を取得
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド（出力ファイルはmain）
RUN go build -o main

# 実行用の軽量なDistrolessイメージを作成
FROM gcr.io/distroless/base-debian12

# 作業ディレクトリを設定
WORKDIR /app

# ビルドされたバイナリをコピー
COPY --from=build /app/main .

# ポート8080を公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./main"]
