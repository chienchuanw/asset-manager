/**
 * 分類管理元件
 * 使用 Badge 顯示分類，提供簡約的管理介面
 */

import { useState } from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { useCategories } from "@/hooks";
import { CashFlowType, type CashFlowCategory } from "@/types/cash-flow";
import { EditCategoryDialog } from "./EditCategoryDialog";
import { DeleteCategoryDialog } from "./DeleteCategoryDialog";
import { AddCategoryDialog } from "./AddCategoryDialog";
import { Lock, Pencil, Trash2, Plus } from "lucide-react";

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
  const [showAddIncome, setShowAddIncome] = useState(false);
  const [showAddExpense, setShowAddExpense] = useState(false);

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
            <TabsContent value="income" className="space-y-3">
              {isLoadingIncome ? (
                <div className="flex flex-wrap gap-2">
                  <Skeleton className="h-7 w-20" />
                  <Skeleton className="h-7 w-24" />
                  <Skeleton className="h-7 w-16" />
                </div>
              ) : incomeError ? (
                <div className="text-sm text-destructive">
                  載入失敗：{incomeError.message}
                </div>
              ) : (
                <>
                  <div className="flex flex-wrap gap-2">
                    {incomeCategories?.map((category) => (
                      <Badge
                        key={category.id}
                        variant={category.is_system ? "secondary" : "outline"}
                        className={
                          category.is_system
                            ? ""
                            : "group relative pr-2 hover:pr-16 transition-all cursor-pointer"
                        }
                      >
                        {category.is_system && (
                          <Lock className="mr-1 h-2.5 w-2.5" />
                        )}
                        {category.name}
                        {!category.is_system && (
                          <div className="absolute right-1 hidden group-hover:flex gap-0.5">
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-5 w-5"
                              onClick={() => handleEdit(category)}
                            >
                              <Pencil className="h-2.5 w-2.5" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-5 w-5 text-destructive"
                              onClick={() => handleDelete(category)}
                            >
                              <Trash2 className="h-2.5 w-2.5" />
                            </Button>
                          </div>
                        )}
                      </Badge>
                    ))}
                    <Badge
                      variant="outline"
                      className="cursor-pointer hover:bg-accent"
                      onClick={() => setShowAddIncome(true)}
                    >
                      <Plus className="h-3 w-3 mr-1" />
                      新增收入分類
                    </Badge>
                  </div>
                </>
              )}
            </TabsContent>

            {/* 支出分類 */}
            <TabsContent value="expense" className="space-y-3">
              {isLoadingExpense ? (
                <div className="flex flex-wrap gap-2">
                  <Skeleton className="h-7 w-20" />
                  <Skeleton className="h-7 w-24" />
                  <Skeleton className="h-7 w-16" />
                </div>
              ) : expenseError ? (
                <div className="text-sm text-destructive">
                  載入失敗：{expenseError.message}
                </div>
              ) : (
                <>
                  <div className="flex flex-wrap gap-2">
                    {expenseCategories?.map((category) => (
                      <Badge
                        key={category.id}
                        variant={category.is_system ? "secondary" : "outline"}
                        className={
                          category.is_system
                            ? ""
                            : "group relative pr-2 hover:pr-16 transition-all cursor-pointer"
                        }
                      >
                        {category.is_system && (
                          <Lock className="mr-1 h-2.5 w-2.5" />
                        )}
                        {category.name}
                        {!category.is_system && (
                          <div className="absolute right-1 hidden group-hover:flex gap-0.5">
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-5 w-5"
                              onClick={() => handleEdit(category)}
                            >
                              <Pencil className="h-2.5 w-2.5" />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-5 w-5 text-destructive"
                              onClick={() => handleDelete(category)}
                            >
                              <Trash2 className="h-2.5 w-2.5" />
                            </Button>
                          </div>
                        )}
                      </Badge>
                    ))}
                    <Badge
                      variant="outline"
                      className="cursor-pointer hover:bg-accent"
                      onClick={() => setShowAddExpense(true)}
                    >
                      <Plus className="h-3 w-3 mr-1" />
                      新增支出分類
                    </Badge>
                  </div>
                </>
              )}
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>

      {/* 新增對話框 - 收入 */}
      <AddCategoryDialog
        type={CashFlowType.INCOME}
        open={showAddIncome}
        onOpenChange={setShowAddIncome}
      />

      {/* 新增對話框 - 支出 */}
      <AddCategoryDialog
        type={CashFlowType.EXPENSE}
        open={showAddExpense}
        onOpenChange={setShowAddExpense}
      />

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
