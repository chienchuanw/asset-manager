import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */

  // 啟用 standalone 模式,減少 Docker 映像檔大小
  output: "standalone",

  // 圖片優化設定
  images: {
    unoptimized: true, // 在 Docker 環境中停用圖片優化
  },
};

export default nextConfig;
