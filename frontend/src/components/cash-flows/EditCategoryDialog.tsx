/**
 * 編輯分類對話框元件
 * 用於編輯分類名稱
 */

import { useEffect } from "react";
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
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { useUpdateCategory } from "@/hooks";
import {
  updateCategorySchema,
  type UpdateCategoryFormData,
  type CashFlowCategory,
} from "@/types/cash-flow";
import { toast } from "sonner";

interface EditCategoryDialogProps {
  category: CashFlowCategory | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

/**
 * 編輯分類對話框元件
 * 
 * 提供表單讓使用者編輯分類名稱
 */
export function EditCategoryDialog({
  category,
  open,
  onOpenChange,
}: EditCategoryDialogProps) {
  // 表單設定
  const form = useForm<UpdateCategoryFormData>({
    resolver: zodResolver(updateCategorySchema),
    defaultValues: {
      name: "",
    },
  });

  // 當分類變更時，更新表單預設值
  useEffect(() => {
    if (category) {
      form.reset({
        name: category.name,
      });
    }
  }, [category, form]);

  // 更新分類 mutation
  const updateMutation = useUpdateCategory({
    onSuccess: () => {
      toast.success("分類更新成功");
      onOpenChange(false);
    },
    onError: (error) => {
      toast.error(error.message || "更新失敗");
    },
  });

  // 送出表單
  const onSubmit = (data: UpdateCategoryFormData) => {
    if (!category) return;
    updateMutation.mutate({
      id: category.id,
      data,
    });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>編輯分類</DialogTitle>
          <DialogDescription>
            修改分類名稱。系統預設分類無法編輯。
          </DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>分類名稱</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="輸入分類名稱"
                      {...field}
                      disabled={updateMutation.isPending}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={updateMutation.isPending}
              >
                取消
              </Button>
              <Button type="submit" disabled={updateMutation.isPending}>
                {updateMutation.isPending ? "更新中..." : "確認"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}

