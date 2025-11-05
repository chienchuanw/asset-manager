"use client";

import React from "react";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useCategories } from "@/hooks";
import { CashFlowType } from "@/types/cash-flow";
import { Loader2 } from "lucide-react";

interface CategorySelectProps {
  value: string;
  onValueChange: (value: string) => void;
  type: CashFlowType;
  placeholder?: string;
  disabled?: boolean;
  // 當提供此名稱時，元件在分類資料載入後會自動選擇該名稱的分類（若目前尚未選定）
  autoSelectName?: string;
}

/**
 * 分類選擇器元件
 *
 * 根據現金流類型（收入/支出）顯示對應的分類選項
 */
export function CategorySelect({
  value,
  onValueChange,
  type,
  placeholder = "選擇分類",
  disabled = false,
  autoSelectName,
}: CategorySelectProps) {
  // 取得分類列表
  const { data: categories, isLoading } = useCategories(type);

  // 當提供 autoSelectName 且目前尚未選擇分類時，於分類載入後自動選擇對應分類
  React.useEffect(() => {
    // 僅在提供 autoSelectName 且當前 value 為空或不在清單中時自動帶入
    // 這可避免父層 setValue 與子元件選單載入時間差造成的「閃跳後又被清空」
    if (!autoSelectName || !categories) return;

    const exists = value && categories.some((c) => c.id === value);
    if (exists) return;

    const target = categories.find((c) => c.name === autoSelectName);
    if (target) {
      onValueChange(target.id);
    }
  }, [autoSelectName, categories, value, onValueChange]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-10 border rounded-md bg-muted">
        <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <Select value={value} onValueChange={onValueChange} disabled={disabled}>
      <SelectTrigger>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {categories?.map((category) => (
          <SelectItem key={category.id} value={category.id}>
            {category.name}
            {category.is_system && (
              <span className="ml-2 text-xs text-muted-foreground">(系統)</span>
            )}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
