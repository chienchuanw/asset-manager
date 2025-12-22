/**
 * i18n 配置
 * 定義支援的語言和預設語言
 */

// 支援的語言列表
export const locales = ["zh-TW", "en"] as const;

// 語言類型
export type Locale = (typeof locales)[number];

// 預設語言
export const defaultLocale: Locale = "zh-TW";

// 語言顯示名稱
export const localeNames: Record<Locale, string> = {
  "zh-TW": "繁體中文",
  en: "English",
};

// localStorage key
export const LOCALE_STORAGE_KEY = "locale";

