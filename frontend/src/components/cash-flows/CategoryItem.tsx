/**
 * 分類項目元件
 * 顯示單個分類，包含名稱、系統標記、編輯/刪除按鈕
 */

import { CashFlowCategory } from "@/types/cash-flow";
import { Button } from "@/components/ui/button";
import { Pencil, Trash2, Lock } from "lucide-react";

interface CategoryItemProps {
  category: CashFlowCategory;
  onEdit: (category: CashFlowCategory) => void;
  onDelete: (category: CashFlowCategory) => void;
}

/**
 * 分類項目元件
 * 
 * 顯示分類名稱，系統分類會顯示鎖頭圖示並禁用編輯/刪除功能
 */
export function CategoryItem({
  category,
  onEdit,
  onDelete,
}: CategoryItemProps) {
  return (
    <div className="flex items-center justify-between py-2 px-3 rounded-md hover:bg-muted/50 group">
      <div className="flex items-center gap-2">
        {category.is_system && (
          <Lock className="h-4 w-4 text-muted-foreground" />
        )}
        <span
          className={category.is_system ? "text-muted-foreground" : ""}
        >
          {category.name}
        </span>
      </div>

      {!category.is_system && (
        <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onEdit(category)}
            className="h-8 w-8 p-0"
          >
            <Pencil className="h-4 w-4" />
            <span className="sr-only">編輯</span>
          </Button>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDelete(category)}
            className="h-8 w-8 p-0 text-destructive hover:text-destructive"
          >
            <Trash2 className="h-4 w-4" />
            <span className="sr-only">刪除</span>
          </Button>
        </div>
      )}
    </div>
  );
}

