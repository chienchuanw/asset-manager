"use client";

import { useLocale as useNextIntlLocale } from "next-intl";
import { useLocale } from "@/providers/LocaleProvider";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Globe } from "lucide-react";

/**
 * 語言切換元件
 * 提供下拉選單讓使用者在繁體中文和英文之間切換
 */
export function LanguageSwitcher() {
  const currentLocale = useNextIntlLocale();
  const { setLocale } = useLocale();

  // 語言選項
  const languages = [
    { code: "zh-TW", label: "繁體中文" },
    { code: "en", label: "English" },
  ] as const;

  // 處理語言切換
  const handleLanguageChange = (newLocale: "zh-TW" | "en") => {
    setLocale(newLocale);
  };

  // 取得當前語言的顯示名稱
  const currentLanguage = languages.find((lang) => lang.code === currentLocale);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="sm" className="gap-2">
          <Globe className="h-4 w-4" />
          <span className="hidden sm:inline">{currentLanguage?.label}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        {languages.map((lang) => (
          <DropdownMenuItem
            key={lang.code}
            onClick={() => handleLanguageChange(lang.code)}
            className={currentLocale === lang.code ? "bg-accent" : ""}
          >
            {lang.label}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
