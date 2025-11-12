/**
 * 新增分類對話框元件
 * 提供表單介面來建立新的現金流分類
 */

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useCreateCategory } from "@/hooks";
import {
  createCategorySchema,
  type CreateCategoryInput,
  type CashFlowType,
} from "@/types/cash-flow";
import { toast } from "sonner";

interface AddCategoryDialogProps {
  type: CashFlowType;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

/**
 * 新增分類對話框元件
 * 
 * 提供表單介面來建立新的現金流分類
 */
export function AddCategoryDialog({
  type,
  open,
  onOpenChange,
}: AddCategoryDialogProps) {
  // 表單處理
  const form = useForm<CreateCategoryInput>({
    resolver: zodResolver(createCategorySchema),
    defaultValues: {
      name: "",
      type,
    },
  });

  // 建立分類 mutation
  const createMutation = useCreateCategory({
    onSuccess: () => {
      toast.success("分類建立成功");
      form.reset({ name: "", type });
      onOpenChange(false);
    },
    onError: (error) => {
      toast.error(error.message || "建立失敗");
    },
  });

  // 提交表單
  const onSubmit = (data: CreateCategoryInput) => {
    createMutation.mutate(data);
  };

  // 當 type 改變時，更新表單的 type 值
  const handleOpenChange = (newOpen: boolean) => {
    if (newOpen) {
      form.setValue("type", type);
    } else {
      form.reset({ name: "", type });
    }
    onOpenChange(newOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>新增分類</DialogTitle>
          <DialogDescription>
            建立新的{type === "income" ? "收入" : "支出"}分類
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">分類名稱</Label>
            <Input
              id="name"
              placeholder="請輸入分類名稱（最多 20 字元）"
              {...form.register("name")}
              disabled={createMutation.isPending}
            />
            {form.formState.errors.name && (
              <p className="text-sm text-destructive">
                {form.formState.errors.name.message}
              </p>
            )}
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => handleOpenChange(false)}
              disabled={createMutation.isPending}
            >
              取消
            </Button>
            <Button type="submit" disabled={createMutation.isPending}>
              {createMutation.isPending ? "建立中..." : "建立"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

