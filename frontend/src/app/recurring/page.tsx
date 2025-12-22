/**
 * 訂閱分期管理頁面
 * 整合訂閱和分期的管理功能
 */

"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
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
  // i18n hooks
  const t = useTranslations("recurring");
  const tCommon = useTranslations("common");

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
      // 將日期格式轉換為 ISO 8601 格式（後端期望的格式）
      const formattedData = {
        ...data,
        start_date: data.start_date
          ? new Date(data.start_date).toISOString()
          : new Date().toISOString(),
        end_date: data.end_date
          ? new Date(data.end_date).toISOString()
          : undefined,
      };

      if (editingSubscription) {
        await updateSubscription.mutateAsync({
          id: editingSubscription.id,
          data: formattedData,
        });
        toast({
          title: t("updateSuccess"),
          description: t("subscriptionUpdated"),
        });
      } else {
        await createSubscription.mutateAsync(formattedData);
        toast({
          title: t("createSuccess"),
          description: t("subscriptionCreated"),
        });
      }
      setSubscriptionDialogOpen(false);
      setEditingSubscription(undefined);
    } catch (error) {
      toast({
        title: t("operationFailed"),
        description: error instanceof Error ? error.message : t("unknownError"),
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
        data: { end_date: new Date().toISOString() },
      });
      toast({
        title: t("cancelSuccess"),
        description: t("subscriptionCanceled"),
      });
      setCancelingSubscription(undefined);
    } catch (error) {
      toast({
        title: t("cancelFailed"),
        description: error instanceof Error ? error.message : t("unknownError"),
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
        title: t("deleteSuccess"),
        description: t("subscriptionDeleted"),
      });
      setDeletingSubscriptionId(undefined);
    } catch (error) {
      toast({
        title: t("deleteFailed"),
        description: error instanceof Error ? error.message : t("unknownError"),
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
      // 將日期格式轉換為 ISO 8601 格式（後端期望的格式）
      const formattedData = {
        ...data,
        start_date: data.start_date
          ? new Date(data.start_date).toISOString()
          : new Date().toISOString(),
      };

      if (editingInstallment) {
        await updateInstallment.mutateAsync({
          id: editingInstallment.id,
          data: {
            name: formattedData.name,
            billing_day: formattedData.billing_day,
            note: formattedData.note,
          },
        });
        toast({
          title: t("updateSuccess"),
          description: t("installmentUpdated"),
        });
      } else {
        await createInstallment.mutateAsync(formattedData);
        toast({
          title: t("createSuccess"),
          description: t("installmentCreated"),
        });
      }
      setInstallmentDialogOpen(false);
      setEditingInstallment(undefined);
    } catch (error) {
      toast({
        title: t("operationFailed"),
        description: error instanceof Error ? error.message : t("unknownError"),
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
        title: t("deleteSuccess"),
        description: t("installmentDeleted"),
      });
      setDeletingInstallmentId(undefined);
    } catch (error) {
      toast({
        title: t("deleteFailed"),
        description: error instanceof Error ? error.message : t("unknownError"),
        variant: "destructive",
      });
    }
  };

  return (
    <AppLayout title={t("title")} description={t("description")}>
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
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
                <TabsTrigger value="subscriptions">
                  {t("subscriptionServices")}
                </TabsTrigger>
                <TabsTrigger value="installments">
                  {t("installmentPayments")}
                </TabsTrigger>
              </TabsList>
              <Button
                onClick={
                  activeTab === "subscriptions"
                    ? handleCreateSubscription
                    : handleCreateInstallment
                }
              >
                <PlusIcon className="mr-2 h-4 w-4" />
                {activeTab === "subscriptions"
                  ? t("addSubscription")
                  : t("addInstallment")}
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
      </div>

      {/* 訂閱表單對話框 */}
      <Dialog
        open={subscriptionDialogOpen}
        onOpenChange={setSubscriptionDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingSubscription
                ? t("editSubscription")
                : t("addSubscription")}
            </DialogTitle>
          </DialogHeader>
          {/* key 屬性確保每次編輯不同訂閱時表單重新掛載 */}
          <SubscriptionForm
            key={editingSubscription?.id ?? "new"}
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
              {editingInstallment ? t("editInstallment") : t("addInstallment")}
            </DialogTitle>
          </DialogHeader>
          {/* key 屬性確保每次編輯不同分期時表單重新掛載 */}
          <InstallmentForm
            key={editingInstallment?.id ?? "new"}
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
            <AlertDialogTitle>
              {t("confirmCancelSubscription")}
            </AlertDialogTitle>
            <AlertDialogDescription>
              {t("cancelSubscriptionConfirmMessage", {
                name: cancelingSubscription?.name ?? "",
              })}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={confirmCancelSubscription}>
              {tCommon("confirm")}
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
            <AlertDialogTitle>{t("confirmDelete")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("deleteSubscriptionConfirmMessage")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteSubscription}
              className="bg-destructive text-destructive-foreground"
            >
              {tCommon("delete")}
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
            <AlertDialogTitle>{t("confirmDelete")}</AlertDialogTitle>
            <AlertDialogDescription>
              {t("deleteInstallmentConfirmMessage")}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tCommon("cancel")}</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDeleteInstallment}
              className="bg-destructive text-destructive-foreground"
            >
              {tCommon("delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </AppLayout>
  );
}
