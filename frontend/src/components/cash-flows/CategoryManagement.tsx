/**
 * 分類管理元件
 * 包含收入/支出分頁、列表顯示、新增表單
 */

import { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { useCategories } from "@/hooks";
import { CashFlowType, type CashFlowCategory } from "@/types/cash-flow";
import { CategoryItem } from "./CategoryItem";
import { AddCategoryForm } from "./AddCategoryForm";
import { EditCategoryDialog } from "./EditCategoryDialog";
import { DeleteCategoryDialog } from "./DeleteCategoryDialog";

/**
 * 分類管理元件
 * 
 * 提供收入和支出分類的管理介面
 */
export function CategoryManagement() {
  const [editingCategory, setEditingCategory] =
    useState<CashFlowCategory | null>(null);
  const [deletingCategory, setDeletingCategory] =
    useState<CashFlowCategory | null>(null);

  // 取得收入分類
  const {
    data: incomeCategories,
    isLoading: isLoadingIncome,
    error: incomeError,
  } = useCategories(CashFlowType.INCOME);

  // 取得支出分類
  const {
    data: expenseCategories,
    isLoading: isLoadingExpense,
    error: expenseError,
  } = useCategories(CashFlowType.EXPENSE);

  // 處理編輯
  const handleEdit = (category: CashFlowCategory) => {
    setEditingCategory(category);
  };

  // 處理刪除
  const handleDelete = (category: CashFlowCategory) => {
    setDeletingCategory(category);
  };

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle>分類管理</CardTitle>
          <CardDescription>
            管理收入和支出分類。系統預設分類無法編輯或刪除。
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="income" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="income">收入分類</TabsTrigger>
              <TabsTrigger value="expense">支出分類</TabsTrigger>
            </TabsList>

            {/* 收入分類 */}
            <TabsContent value="income" className="space-y-4">
              {isLoadingIncome ? (
                <div className="space-y-2">
                  <Skeleton className="h-10 w-full" />
                  <Skeleton className="h-10 w-full" />
                  <Skeleton className="h-10 w-full" />
                </div>
              ) : incomeError ? (
                <div className="text-sm text-destructive">
                  載入失敗：{incomeError.message}
                </div>
              ) : (
                <>
                  <div className="space-y-1">
                    {incomeCategories?.map((category) => (
                      <CategoryItem
                        key={category.id}
                        category={category}
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                      />
                    ))}
                  </div>
                  <AddCategoryForm type={CashFlowType.INCOME} />
                </>
              )}
            </TabsContent>

            {/* 支出分類 */}
            <TabsContent value="expense" className="space-y-4">
              {isLoadingExpense ? (
                <div className="space-y-2">
                  <Skeleton className="h-10 w-full" />
                  <Skeleton className="h-10 w-full" />
                  <Skeleton className="h-10 w-full" />
                </div>
              ) : expenseError ? (
                <div className="text-sm text-destructive">
                  載入失敗：{expenseError.message}
                </div>
              ) : (
                <>
                  <div className="space-y-1">
                    {expenseCategories?.map((category) => (
                      <CategoryItem
                        key={category.id}
                        category={category}
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                      />
                    ))}
                  </div>
                  <AddCategoryForm type={CashFlowType.EXPENSE} />
                </>
              )}
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* 編輯對話框 */}
      <EditCategoryDialog
        category={editingCategory}
        open={!!editingCategory}
        onOpenChange={(open) => !open && setEditingCategory(null)}
      />

      {/* 刪除確認對話框 */}
      <DeleteCategoryDialog
        category={deletingCategory}
        open={!!deletingCategory}
        onOpenChange={(open) => !open && setDeletingCategory(null)}
      />
    </>
  );
}

