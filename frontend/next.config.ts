import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const nextConfig: NextConfig = {
  /* config options here */

  // 啟用 standalone 模式,減少 Docker 映像檔大小
  output: "standalone",

  // 圖片優化設定
  images: {
    unoptimized: true, // 在 Docker 環境中停用圖片優化
  },
};

// 設定 next-intl plugin，指定 i18n 請求配置檔案位置
const withNextIntl = createNextIntlPlugin("./src/i18n/request.ts");

export default withNextIntl(nextConfig);
