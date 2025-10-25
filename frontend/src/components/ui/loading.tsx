/**
 * Loading 組件
 * 統一的載入狀態顯示組件，支援多種場景使用
 */

import { Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";

interface LoadingProps {
  /**
   * 顯示變體
   * - page: 全頁面居中顯示
   * - inline: 行內顯示
   * - overlay: 覆蓋層顯示（帶背景模糊）
   */
  variant?: "page" | "inline" | "overlay";
  /**
   * 圖示大小
   */
  size?: "sm" | "md" | "lg";
  /**
   * 顯示文字
   */
  text?: string;
  /**
   * 自訂 className
   */
  className?: string;
}

/**
 * Loading 組件
 * 
 * @example
 * // 全頁面 Loading
 * <Loading variant="page" size="lg" text="載入資料中..." />
 * 
 * @example
 * // Inline Loading
 * <Loading variant="inline" size="sm" text="載入中..." />
 * 
 * @example
 * // Overlay Loading
 * <Loading variant="overlay" size="md" text="更新中..." />
 */
export function Loading({
  variant = "inline",
  size = "md",
  text,
  className,
}: LoadingProps) {
  // 圖示大小對應的 class
  const sizeClasses = {
    sm: "h-4 w-4",
    md: "h-8 w-8",
    lg: "h-12 w-12",
  };

  // 文字大小對應的 class
  const textSizeClasses = {
    sm: "text-xs",
    md: "text-sm",
    lg: "text-base",
  };

  // Loading 內容
  const content = (
    <div className={cn("flex flex-col items-center gap-3", className)}>
      <Loader2
        className={cn(
          "animate-spin text-muted-foreground",
          sizeClasses[size]
        )}
      />
      {text && (
        <p className={cn("text-muted-foreground", textSizeClasses[size])}>
          {text}
        </p>
      )}
    </div>
  );

  // 全頁面變體
  if (variant === "page") {
    return (
      <div className="flex items-center justify-center min-h-[400px]">
        {content}
      </div>
    );
  }

  // 覆蓋層變體
  if (variant === "overlay") {
    return (
      <div className="absolute inset-0 flex items-center justify-center bg-background/80 backdrop-blur-sm z-50">
        {content}
      </div>
    );
  }

  // 行內變體（預設）
  return content;
}

