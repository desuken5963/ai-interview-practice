import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  reactStrictMode: false,
  compiler: {
    // ハイドレーションエラーを抑制するための設定
    styledComponents: true,
  },
};

export default nextConfig;
