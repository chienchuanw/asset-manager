"use client";

/**
 * LocaleProvider
 * 處理語言偏好的讀取和儲存，同時同步到 localStorage 和 cookie
 * cookie 用於 Server Component 讀取，localStorage 用於持久化儲存
 */

import {
  createContext,
  useContext,
  useEffect,
  useState,
  useCallback,
  ReactNode,
} from "react";
import { useRouter } from "next/navigation";
import {
  Locale,
  defaultLocale,
  locales,
  LOCALE_STORAGE_KEY,
} from "@/i18n/config";

interface LocaleContextType {
  locale: Locale;
  setLocale: (locale: Locale) => void;
  isLoading: boolean;
}

const LocaleContext = createContext<LocaleContextType | undefined>(undefined);

interface LocaleProviderProps {
  children: ReactNode;
  initialLocale?: Locale;
}

export function LocaleProvider({
  children,
  initialLocale,
}: LocaleProviderProps) {
  const router = useRouter();
  const [locale, setLocaleState] = useState<Locale>(
    initialLocale || defaultLocale
  );
  const [isLoading, setIsLoading] = useState(true);

  // 初始化時從 localStorage 讀取語言偏好
  useEffect(() => {
    const storedLocale = localStorage.getItem(LOCALE_STORAGE_KEY);
    if (storedLocale && locales.includes(storedLocale as Locale)) {
      setLocaleState(storedLocale as Locale);
      // 同步到 cookie
      document.cookie = `${LOCALE_STORAGE_KEY}=${storedLocale};path=/;max-age=31536000`;
    }
    setIsLoading(false);
  }, []);

  // 設定語言並同步到 localStorage 和 cookie
  const setLocale = useCallback(
    (newLocale: Locale) => {
      if (!locales.includes(newLocale)) return;

      setLocaleState(newLocale);

      // 儲存到 localStorage
      localStorage.setItem(LOCALE_STORAGE_KEY, newLocale);

      // 同步到 cookie，設定一年的有效期
      document.cookie = `${LOCALE_STORAGE_KEY}=${newLocale};path=/;max-age=31536000`;

      // 刷新頁面以套用新語言
      router.refresh();
    },
    [router]
  );

  return (
    <LocaleContext.Provider value={{ locale, setLocale, isLoading }}>
      {children}
    </LocaleContext.Provider>
  );
}

// 自訂 hook 用於取得語言設定
export function useLocale() {
  const context = useContext(LocaleContext);
  if (context === undefined) {
    throw new Error("useLocale must be used within a LocaleProvider");
  }
  return context;
}

