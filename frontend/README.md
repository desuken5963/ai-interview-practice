This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## バックエンドとの連携

このフロントエンドアプリケーションは、バックエンドAPIと連携して動作します。バックエンドAPIは以下の手順で起動してください。

1. バックエンドディレクトリに移動します。
```bash
cd ../backend
```

2. バックエンドサーバーを起動します。
```bash
make run
```

3. バックエンドサーバーが `http://localhost:8080` で起動します。

## 環境変数の設定

フロントエンドとバックエンドの連携には、環境変数の設定が必要です。`.env.local` ファイルを作成して、以下の環境変数を設定してください。

```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

バックエンドのURLが異なる場合は、適宜変更してください。

## API連携の仕組み

フロントエンドとバックエンドの連携は、`src/lib/api/client.ts` で定義されたAPIクライアントを通じて行われます。このクライアントは、以下の機能を提供します。

- 企業情報の取得、作成、更新、削除
- 求人情報の取得、作成、更新、削除

APIクライアントは、環境変数 `NEXT_PUBLIC_API_URL` で指定されたバックエンドAPIのURLに対してリクエストを送信します。

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.
