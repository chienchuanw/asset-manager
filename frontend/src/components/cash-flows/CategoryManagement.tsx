/**
 * 分類管理元件
 * 使用 Badge 顯示分類，提供簡約的管理介面
 * 支援拖拉排序功能
 */

import { useState } from "react";
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  horizontalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
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
import { useCategories, useReorderCategories } from "@/hooks";
import { CashFlowType, type CashFlowCategory } from "@/types/cash-flow";
import { EditCategoryDialog } from "./EditCategoryDialog";
import { DeleteCategoryDialog } from "./DeleteCategoryDialog";
import { AddCategoryDialog } from "./AddCategoryDialog";
import { Lock, Pencil, Trash2, Plus, GripVertical } from "lucide-react";
import { toast } from "sonner";

/**
 * 可拖拉排序的分類 Badge 元件
 */
interface SortableCategoryBadgeProps {
  category: CashFlowCategory;
  onEdit: (category: CashFlowCategory) => void;
  onDelete: (category: CashFlowCategory) => void;
}

function SortableCategoryBadge({
  category,
  onEdit,
  onDelete,
}: SortableCategoryBadgeProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: category.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <Badge
      ref={setNodeRef}
      style={style}
      variant={category.is_system ? "secondary" : "outline"}
      className={
        category.is_system
          ? "cursor-grab active:cursor-grabbing"
          : "group relative pr-2 hover:pr-16 transition-all cursor-grab active:cursor-grabbing"
      }
    >
      {/* 拖拉把手 */}
      <span {...attributes} {...listeners} className="mr-1 touch-none">
        <GripVertical className="h-2.5 w-2.5 text-muted-foreground" />
      </span>
      {category.is_system && <Lock className="mr-1 h-2.5 w-2.5" />}
      {category.name}
      {!category.is_system && (
        <div className="absolute right-1 hidden group-hover:flex gap-0.5">
          <Button
            variant="ghost"
            size="icon"
            className="h-5 w-5"
            onClick={(e) => {
              e.stopPropagation();
              onEdit(category);
            }}
          >
            <Pencil className="h-2.5 w-2.5" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-5 w-5 text-destructive"
            onClick={(e) => {
              e.stopPropagation();
              onDelete(category);
            }}
          >
            <Trash2 className="h-2.5 w-2.5" />
          </Button>
        </div>
      )}
    </Badge>
  );
}

/**
 * 分類管理元件
 *
 * 提供收入和支出分類的管理介面，支援拖拉排序
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
    data: incomeCategoriesData,
    isLoading: isLoadingIncome,
    error: incomeError,
  } = useCategories(CashFlowType.INCOME);

  // 取得支出分類
  const {
    data: expenseCategoriesData,
    isLoading: isLoadingExpense,
    error: expenseError,
  } = useCategories(CashFlowType.EXPENSE);

  // 重新排序 mutation
  const reorderMutation = useReorderCategories({
    onSuccess: () => {
      toast.success("排序更新成功");
    },
    onError: (error) => {
      toast.error(`排序更新失敗：${error.message}`);
    },
  });

  // 依照 sort_order 排序分類
  const incomeCategories = incomeCategoriesData?.slice().sort((a, b) => {
    return a.sort_order - b.sort_order;
  });

  const expenseCategories = expenseCategoriesData?.slice().sort((a, b) => {
    return a.sort_order - b.sort_order;
  });

  // 設定拖拉感應器
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  // 處理拖拉結束
  const handleDragEnd = (
    event: DragEndEvent,
    categories: CashFlowCategory[]
  ) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      const oldIndex = categories.findIndex((c) => c.id === active.id);
      const newIndex = categories.findIndex((c) => c.id === over.id);

      // 計算新的排序
      const reorderedCategories = arrayMove(categories, oldIndex, newIndex);

      // 建立排序資料
      const orders = reorderedCategories.map((category, index) => ({
        id: category.id,
        sort_order: index,
      }));

      // 呼叫 API 更新排序
      reorderMutation.mutate({ orders });
    }
  };

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
              ) : incomeCategories ? (
                <DndContext
                  sensors={sensors}
                  collisionDetection={closestCenter}
                  onDragEnd={(event) => handleDragEnd(event, incomeCategories)}
                >
                  <SortableContext
                    items={incomeCategories.map((c) => c.id)}
                    strategy={horizontalListSortingStrategy}
                  >
                    <div className="flex flex-wrap gap-2">
                      {incomeCategories.map((category) => (
                        <SortableCategoryBadge
                          key={category.id}
                          category={category}
                          onEdit={handleEdit}
                          onDelete={handleDelete}
                        />
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
                  </SortableContext>
                </DndContext>
              ) : null}
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
              ) : expenseCategories ? (
                <DndContext
                  sensors={sensors}
                  collisionDetection={closestCenter}
                  onDragEnd={(event) => handleDragEnd(event, expenseCategories)}
                >
                  <SortableContext
                    items={expenseCategories.map((c) => c.id)}
                    strategy={horizontalListSortingStrategy}
                  >
                    <div className="flex flex-wrap gap-2">
                      {expenseCategories.map((category) => (
                        <SortableCategoryBadge
                          key={category.id}
                          category={category}
                          onEdit={handleEdit}
                          onDelete={handleDelete}
                        />
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
                  </SortableContext>
                </DndContext>
              ) : null}
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
