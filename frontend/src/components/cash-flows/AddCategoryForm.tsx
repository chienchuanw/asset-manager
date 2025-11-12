/**
 * 新增分類表單元件
 * Inline 表單用於快速新增分類
 */

import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from "@/components/ui/form";
import { useCreateCategory } from "@/hooks";
import {
  createCategorySchema,
  type CreateCategoryFormData,
  type CashFlowType,
} from "@/types/cash-flow";
import { toast } from "sonner";

interface AddCategoryFormProps {
  type: CashFlowType;
}

/**
 * 新增分類表單元件
 * 
 * 提供 inline 表單讓使用者快速新增分類
 */
export function AddCategoryForm({ type }: AddCategoryFormProps) {
  const [isAdding, setIsAdding] = useState(false);

  // 表單設定
  const form = useForm<CreateCategoryFormData>({
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
      setIsAdding(false);
    },
    onError: (error) => {
      toast.error(error.message || "建立失敗");
    },
  });

  // 送出表單
  const onSubmit = (data: CreateCategoryFormData) => {
    createMutation.mutate(data);
  };

  // 取消新增
  const handleCancel = () => {
    form.reset({ name: "", type });
    setIsAdding(false);
  };

  if (!isAdding) {
    return (
      <Button
        variant="outline"
        size="sm"
        onClick={() => setIsAdding(true)}
        className="w-full"
      >
        <Plus className="h-4 w-4 mr-2" />
        新增分類
      </Button>
    );
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-2">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormControl>
                <Input
                  placeholder="輸入分類名稱"
                  autoFocus
                  {...field}
                  disabled={createMutation.isPending}
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        <div className="flex gap-2">
          <Button
            type="submit"
            size="sm"
            disabled={createMutation.isPending}
            className="flex-1"
          >
            {createMutation.isPending ? "建立中..." : "確認"}
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            onClick={handleCancel}
            disabled={createMutation.isPending}
          >
            <X className="h-4 w-4" />
            <span className="sr-only">取消</span>
          </Button>
        </div>
      </form>
    </Form>
  );
}

