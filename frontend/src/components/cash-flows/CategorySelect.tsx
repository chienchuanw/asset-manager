"use client";

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
}: CategorySelectProps) {
  // 取得分類列表
  const { data: categories, isLoading } = useCategories(type);

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
              <span className="ml-2 text-xs text-muted-foreground">
                (系統)
              </span>
            )}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}

