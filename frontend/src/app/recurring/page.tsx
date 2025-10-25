/**
 * 訂閱分期管理頁面
 * 整合訂閱和分期的管理功能
 */

"use client";

import { useState } from "react";
import { AppLayout } from "@/components/layout/AppLayout";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
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
import { SubscriptionsList } from "@/components/dashboard/SubscriptionsList";
import { InstallmentsList } from "@/components/dashboard/InstallmentsList";
import { SubscriptionForm } from "@/components/dashboard/SubscriptionForm";
import { InstallmentForm } from "@/components/dashboard/InstallmentForm";
import { RecurringStatsCard } from "@/components/dashboard/RecurringStatsCard";
import {
  useSubscriptions,
  useCreateSubscription,
  useUpdateSubscription,
  useDeleteSubscription,
  useCancelSubscription,
  useInstallments,
  useCreateInstallment,
  useUpdateInstallment,
  useDeleteInstallment,
  useCategories,
} from "@/hooks";
import { PlusIcon } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import type {
  Subscription,
  CreateSubscriptionInput,
} from "@/types/subscription";
import type { Installment, CreateInstallmentInput } from "@/types/installment";

export default function RecurringPage() {
  const { toast } = useToast();
  const [activeTab, setActiveTab] = useState("subscriptions");

  // 訂閱相關狀態
  const [subscriptionDialogOpen, setSubscriptionDialogOpen] = useState(false);
  const [editingSubscription, setEditingSubscription] = useState<
    Subscription | undefined
  >();
  const [cancelingSubscription, setCancelingSubscription] = useState<
    Subscription | undefined
  >();
  const [deletingSubscriptionId, setDeletingSubscriptionId] = useState<
    string | undefined
  >();

  // 分期相關狀態
  const [installmentDialogOpen, setInstallmentDialogOpen] = useState(false);
  const [editingInstallment, setEditingInstallment] = useState<
    Installment | undefined
  >();
  const [deletingInstallmentId, setDeletingInstallmentId] = useState<
    string | undefined
  >();

  // 資料查詢
  const { data: subscriptions, isLoading: subscriptionsLoading } =
    useSubscriptions();
  const { data: installments, isLoading: installmentsLoading } =
    useInstallments();
  const { data: categories } = useCategories();

  // 訂閱 mutations
  const createSubscription = useCreateSubscription();
  const updateSubscription = useUpdateSubscription();
  const deleteSubscription = useDeleteSubscription();
  const cancelSubscription = useCancelSubscription();

  // 分期 mutations
  const createInstallment = useCreateInstallment();
  const updateInstallment = useUpdateInstallment();
  const deleteInstallment = useDeleteInstallment();

  // 訂閱處理函式
  const handleCreateSubscription = () => {
    setEditingSubscription(undefined);
    setSubscriptionDialogOpen(true);
  };

  const handleEditSubscription = (subscription: Subscription) => {
    setEditingSubscription(subscription);
    setSubscriptionDialogOpen(true);
  };

  const handleSubmitSubscription = async (data: CreateSubscriptionInput) => {
    try {
      if (editingSubscription) {
        await updateSubscription.mutateAsync({
          id: editingSubscription.id,
          data,
        });
        toast({
          title: "更新成功",
          description: "訂閱已更新",
        });
      } else {
        await createSubscription.mutateAsync(data);
        toast({
          title: "建立成功",
          description: "訂閱已建立",
        });
      }
      setSubscriptionDialogOpen(false);
      setEditingSubscription(undefined);
    } catch (error) {
      toast({
        title: "操作失敗",
        description: error instanceof Error ? error.message : "未知錯誤",
        variant: "destructive",
      });
    }
  };

  const handleCancelSubscription = (subscription: Subscription) => {
    setCancelingSubscription(subscription);
  };

  const confirmCancelSubscription = async () => {
    if (!cancelingSubscription) return;

    try {
      await cancelSubscription.mutateAsync({
        id: cancelingSubscription.id,
        data: { end_date: new Date().toISOString().split("T")[0] },
      });
      toast({
        title: "取消成功",
        description: "訂閱已取消",
      });
      setCancelingSubscription(undefined);
    } catch (error) {
      toast({
        title: "取消失敗",
        description: error instanceof Error ? error.message : "未知錯誤",
        variant: "destructive",
      });
    }
  };

  const handleDeleteSubscription = (id: string) => {
    setDeletingSubscriptionId(id);
  };

  const confirmDeleteSubscription = async () => {
    if (!deletingSubscriptionId) return;

    try {
      await deleteSubscription.mutateAsync(deletingSubscriptionId);
      toast({
        title: "刪除成功",
        description: "訂閱已刪除",
      });
      setDeletingSubscriptionId(undefined);
    } catch (error) {
      toast({
        title: "刪除失敗",
        description: error instanceof Error ? error.message : "未知錯誤",
        variant: "destructive",
      });
    }
  };

  // 分期處理函式
  const handleCreateInstallment = () => {
    setEditingInstallment(undefined);
    setInstallmentDialogOpen(true);
  };

  const handleEditInstallment = (installment: Installment) => {
    setEditingInstallment(installment);
    setInstallmentDialogOpen(true);
  };

  const handleSubmitInstallment = async (data: CreateInstallmentInput) => {
    try {
      if (editingInstallment) {
        await updateInstallment.mutateAsync({
          id: editingInstallment.id,
          data: {
            name: data.name,
            billing_day: data.billing_day,
            note: data.note,
          },
        });
        toast({
          title: "更新成功",
          description: "分期已更新",
        });
      } else {
        await createInstallment.mutateAsync(data);
        toast({
          title: "建立成功",
          description: "分期已建立",
        });
      }
      setInstallmentDialogOpen(false);
      setEditingInstallment(undefined);
    } catch (error) {
      toast({
        title: "操作失敗",
        description: error instanceof Error ? error.message : "未知錯誤",
        variant: "destructive",
      });
    }
  };

  const handleDeleteInstallment = (id: string) => {
    setDeletingInstallmentId(id);
  };

  const confirmDeleteInstallment = async () => {
    if (!deletingInstallmentId) return;

    try {
      await deleteInstallment.mutateAsync(deletingInstallmentId);
      toast({
        title: "刪除成功",
        description: "分期已刪除",
      });
      setDeletingInstallmentId(undefined);
    } catch (error) {
      toast({
        title: "刪除失敗",
        description: error instanceof Error ? error.message : "未知錯誤",
        variant: "destructive",
      });
    }
  };

  return (
    <AppLayout title="訂閱與分期" description="管理您的訂閱服務和分期付款">
      {/* Main Content */}
      <main className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* 統計卡片 */}
          <RecurringStatsCard
            subscriptions={subscriptions}
            installments={installments}
            isLoading={subscriptionsLoading || installmentsLoading}
          />

          {/* 分頁 */}
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <div className="flex items-center justify-between">
              <TabsList>
                <TabsTrigger value="subscriptions">訂閱服務</TabsTrigger>
                <TabsTrigger value="installments">分期付款</TabsTrigger>
              </TabsList>
              <Button
                onClick={
                  activeTab === "subscriptions"
                    ? handleCreateSubscription
                    : handleCreateInstallment
                }
              >
                <PlusIcon className="mr-2 h-4 w-4" />
                新增{activeTab === "subscriptions" ? "訂閱" : "分期"}
              </Button>
            </div>

            <TabsContent value="subscriptions" className="space-y-4">
              <SubscriptionsList
                subscriptions={subscriptions}
                isLoading={subscriptionsLoading}
                onEdit={handleEditSubscription}
                onDelete={handleDeleteSubscription}
                onCancel={handleCancelSubscription}
              />
            </TabsContent>

            <TabsContent value="installments" className="space-y-4">
              <InstallmentsList
                installments={installments}
                isLoading={installmentsLoading}
                onEdit={handleEditInstallment}
                onDelete={handleDeleteInstallment}
              />
            </TabsContent>
          </Tabs>
        </div>
      </main>

      {/* 訂閱表單對話框 */}
      <Dialog
        open={subscriptionDialogOpen}
        onOpenChange={setSubscriptionDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingSubscription ? "編輯訂閱" : "新增訂閱"}
            </DialogTitle>
          </DialogHeader>
          <SubscriptionForm
            subscription={editingSubscription}
            categories={categories}
            onSubmit={handleSubmitSubscription}
            onCancel={() => setSubscriptionDialogOpen(false)}
            isSubmitting={
              createSubscription.isPending || updateSubscription.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 分期表單對話框 */}
      <Dialog
        open={installmentDialogOpen}
        onOpenChange={setInstallmentDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingInstallment ? "編輯分期" : "新增分期"}
            </DialogTitle>
          </DialogHeader>
          <InstallmentForm
            installment={editingInstallment}
            categories={categories}
            onSubmit={handleSubmitInstallment}
            onCancel={() => setInstallmentDialogOpen(false)}
            isSubmitting={
              createInstallment.isPending || updateInstallment.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 取消訂閱確認對話框 */}
      <AlertDialog
        open={!!cancelingSubscription}
        onOpenChange={(open) => !open && setCancelingSubscription(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認取消訂閱</AlertDialogTitle>
            <AlertDialogDescription>
              確定要取消「{cancelingSubscription?.name}」的訂閱嗎？
              取消後將不再自動扣款。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction onClick={confirmCancelSubscription}>
              確認
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 刪除訂閱確認對話框 */}
      <AlertDialog
        open={!!deletingSubscriptionId}
        onOpenChange={(open) => !open && setDeletingSubscriptionId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認刪除</AlertDialogTitle>
            <AlertDialogDescription>
              確定要刪除這個訂閱嗎？此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteSubscription}
              className="bg-destructive text-destructive-foreground"
            >
              刪除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 刪除分期確認對話框 */}
      <AlertDialog
        open={!!deletingInstallmentId}
        onOpenChange={(open) => !open && setDeletingInstallmentId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認刪除</AlertDialogTitle>
            <AlertDialogDescription>
              確定要刪除這個分期嗎？此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteInstallment}
              className="bg-destructive text-destructive-foreground"
            >
              刪除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </AppLayout>
  );
}
