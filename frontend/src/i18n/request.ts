/**
 * next-intl 請求配置
 * 為 Server Components 提供語言和翻譯訊息
 */

import { cookies } from "next/headers";
import { getRequestConfig } from "next-intl/server";
import { defaultLocale, locales, LOCALE_STORAGE_KEY, Locale } from "./config";

export default getRequestConfig(async () => {
  // 從 cookie 中獲取語言偏好，如果沒有則使用預設語言
  const cookieStore = await cookies();
  const localeFromCookie = cookieStore.get(LOCALE_STORAGE_KEY)?.value;

  // 驗證語言是否在支援列表中
  const locale: Locale =
    localeFromCookie && locales.includes(localeFromCookie as Locale)
      ? (localeFromCookie as Locale)
      : defaultLocale;

  return {
    locale,
    messages: (await import(`../../messages/${locale}.json`)).default,
  };
});

