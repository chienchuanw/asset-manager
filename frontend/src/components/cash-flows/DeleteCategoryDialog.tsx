/**
 * 刪除分類確認對話框元件
 * 包含關聯記錄檢查和確認
 */

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { useDeleteCategory } from "@/hooks";
import type { CashFlowCategory } from "@/types/cash-flow";
import { toast } from "sonner";

interface DeleteCategoryDialogProps {
  category: CashFlowCategory | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

/**
 * 刪除分類確認對話框元件
 * 
 * 顯示刪除確認訊息，如果分類被使用則顯示錯誤訊息
 */
export function DeleteCategoryDialog({
  category,
  open,
  onOpenChange,
}: DeleteCategoryDialogProps) {
  // 刪除分類 mutation
  const deleteMutation = useDeleteCategory({
    onSuccess: () => {
      toast.success("分類刪除成功");
      onOpenChange(false);
    },
    onError: (error) => {
      // 檢查是否為「分類正在使用」的錯誤
      if (error.message?.includes("being used")) {
        toast.error("無法刪除此分類，因為已有現金流記錄使用此分類");
      } else if (error.message?.includes("system category")) {
        toast.error("無法刪除系統預設分類");
      } else {
        toast.error(error.message || "刪除失敗");
      }
    },
  });

  // 確認刪除
  const handleDelete = () => {
    if (!category) return;
    deleteMutation.mutate(category.id);
  };

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>確認刪除分類</AlertDialogTitle>
          <AlertDialogDescription>
            您確定要刪除「{category?.name}」分類嗎？
            <br />
            <br />
            如果此分類已被現金流記錄使用，將無法刪除。
            <br />
            此操作無法復原。
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={deleteMutation.isPending}>
            取消
          </AlertDialogCancel>
          <AlertDialogAction
            onClick={handleDelete}
            disabled={deleteMutation.isPending}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            {deleteMutation.isPending ? "刪除中..." : "確認刪除"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}

