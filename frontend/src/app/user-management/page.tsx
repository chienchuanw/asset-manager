/**
 * 使用者管理頁面
 * 整合銀行帳戶和信用卡的管理功能
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
import { BankAccountList } from "@/components/user-management/BankAccountList";
import { BankAccountForm } from "@/components/user-management/BankAccountForm";
import { CreditCardList } from "@/components/user-management/CreditCardList";
import { CreditCardForm } from "@/components/user-management/CreditCardForm";
import {
  useBankAccounts,
  useCreateBankAccount,
  useUpdateBankAccount,
  useDeleteBankAccount,
} from "@/hooks/useBankAccounts";
import {
  useCreditCards,
  useCreateCreditCard,
  useUpdateCreditCard,
  useDeleteCreditCard,
} from "@/hooks/useCreditCards";
import { PlusIcon } from "lucide-react";
import { useToast } from "@/hooks/use-toast";
import type {
  BankAccount,
  CreateBankAccountInput,
  CreditCard,
  CreateCreditCardInput,
} from "@/types/user-management";

export default function UserManagementPage() {
  const { toast } = useToast();
  const [activeTab, setActiveTab] = useState("bank-accounts");

  // 銀行帳戶相關狀態
  const [bankAccountDialogOpen, setBankAccountDialogOpen] = useState(false);
  const [editingBankAccount, setEditingBankAccount] = useState<
    BankAccount | undefined
  >();
  const [deletingBankAccountId, setDeletingBankAccountId] = useState<
    string | undefined
  >();

  // 信用卡相關狀態
  const [creditCardDialogOpen, setCreditCardDialogOpen] = useState(false);
  const [editingCreditCard, setEditingCreditCard] = useState<
    CreditCard | undefined
  >();
  const [deletingCreditCardId, setDeletingCreditCardId] = useState<
    string | undefined
  >();

  // 資料查詢
  const { data: bankAccounts, isLoading: bankAccountsLoading } =
    useBankAccounts();
  const { data: creditCards, isLoading: creditCardsLoading } = useCreditCards();

  // Mutations
  const createBankAccountMutation = useCreateBankAccount();
  const updateBankAccountMutation = useUpdateBankAccount();
  const deleteBankAccountMutation = useDeleteBankAccount();
  const createCreditCardMutation = useCreateCreditCard();
  const updateCreditCardMutation = useUpdateCreditCard();
  const deleteCreditCardMutation = useDeleteCreditCard();

  // ==================== 銀行帳戶處理函式 ====================

  const handleCreateBankAccount = () => {
    setEditingBankAccount(undefined);
    setBankAccountDialogOpen(true);
  };

  const handleEditBankAccount = (bankAccount: BankAccount) => {
    setEditingBankAccount(bankAccount);
    setBankAccountDialogOpen(true);
  };

  const handleSubmitBankAccount = async (data: CreateBankAccountInput) => {
    try {
      if (editingBankAccount) {
        await updateBankAccountMutation.mutateAsync({
          id: editingBankAccount.id,
          data,
        });
        toast({
          title: "更新成功",
          description: "銀行帳戶已更新",
        });
      } else {
        await createBankAccountMutation.mutateAsync(data);
        toast({
          title: "新增成功",
          description: "銀行帳戶已新增",
        });
      }
      setBankAccountDialogOpen(false);
      setEditingBankAccount(undefined);
    } catch (error: any) {
      toast({
        title: "操作失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  const handleDeleteBankAccount = async () => {
    if (!deletingBankAccountId) return;

    try {
      await deleteBankAccountMutation.mutateAsync(deletingBankAccountId);
      toast({
        title: "刪除成功",
        description: "銀行帳戶已刪除",
      });
      setDeletingBankAccountId(undefined);
    } catch (error: any) {
      toast({
        title: "刪除失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  // ==================== 信用卡處理函式 ====================

  const handleCreateCreditCard = () => {
    setEditingCreditCard(undefined);
    setCreditCardDialogOpen(true);
  };

  const handleEditCreditCard = (creditCard: CreditCard) => {
    setEditingCreditCard(creditCard);
    setCreditCardDialogOpen(true);
  };

  const handleSubmitCreditCard = async (data: CreateCreditCardInput) => {
    try {
      if (editingCreditCard) {
        await updateCreditCardMutation.mutateAsync({
          id: editingCreditCard.id,
          data,
        });
        toast({
          title: "更新成功",
          description: "信用卡已更新",
        });
      } else {
        await createCreditCardMutation.mutateAsync(data);
        toast({
          title: "新增成功",
          description: "信用卡已新增",
        });
      }
      setCreditCardDialogOpen(false);
      setEditingCreditCard(undefined);
    } catch (error: any) {
      toast({
        title: "操作失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  const handleDeleteCreditCard = async () => {
    if (!deletingCreditCardId) return;

    try {
      await deleteCreditCardMutation.mutateAsync(deletingCreditCardId);
      toast({
        title: "刪除成功",
        description: "信用卡已刪除",
      });
      setDeletingCreditCardId(undefined);
    } catch (error: any) {
      toast({
        title: "刪除失敗",
        description: error.message || "請稍後再試",
        variant: "destructive",
      });
    }
  };

  return (
    <AppLayout title="使用者管理" description="管理您的銀行帳戶和信用卡資訊">
      {/* Main Content */}
      <div className="flex-1 p-4 md:p-6 bg-gray-50">
        <div className="flex flex-col gap-6">
          {/* Tabs */}
          <Tabs value={activeTab} onValueChange={setActiveTab}>
            <div className="flex items-center justify-between">
              <TabsList>
                <TabsTrigger value="bank-accounts">銀行帳戶</TabsTrigger>
                <TabsTrigger value="credit-cards">信用卡</TabsTrigger>
              </TabsList>

              {activeTab === "bank-accounts" && (
                <Button onClick={handleCreateBankAccount}>
                  <PlusIcon className="mr-2 h-4 w-4" />
                  新增銀行帳戶
                </Button>
              )}

              {activeTab === "credit-cards" && (
                <Button onClick={handleCreateCreditCard}>
                  <PlusIcon className="mr-2 h-4 w-4" />
                  新增信用卡
                </Button>
              )}
            </div>

            {/* 銀行帳戶 Tab */}
            <TabsContent value="bank-accounts" className="space-y-4">
              <BankAccountList
                bankAccounts={bankAccounts}
                isLoading={bankAccountsLoading}
                onEdit={handleEditBankAccount}
                onDelete={setDeletingBankAccountId}
              />
            </TabsContent>

            {/* 信用卡 Tab */}
            <TabsContent value="credit-cards" className="space-y-4">
              <CreditCardList
                creditCards={creditCards}
                isLoading={creditCardsLoading}
                onEdit={handleEditCreditCard}
                onDelete={setDeletingCreditCardId}
              />
            </TabsContent>
          </Tabs>
        </div>
      </div>

      {/* 銀行帳戶 Dialog */}
      <Dialog
        open={bankAccountDialogOpen}
        onOpenChange={setBankAccountDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingBankAccount ? "編輯銀行帳戶" : "新增銀行帳戶"}
            </DialogTitle>
          </DialogHeader>
          <BankAccountForm
            bankAccount={editingBankAccount}
            onSubmit={handleSubmitBankAccount}
            onCancel={() => setBankAccountDialogOpen(false)}
            isSubmitting={
              createBankAccountMutation.isPending ||
              updateBankAccountMutation.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 信用卡 Dialog */}
      <Dialog
        open={creditCardDialogOpen}
        onOpenChange={setCreditCardDialogOpen}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>
              {editingCreditCard ? "編輯信用卡" : "新增信用卡"}
            </DialogTitle>
          </DialogHeader>
          <CreditCardForm
            creditCard={editingCreditCard}
            onSubmit={handleSubmitCreditCard}
            onCancel={() => setCreditCardDialogOpen(false)}
            isSubmitting={
              createCreditCardMutation.isPending ||
              updateCreditCardMutation.isPending
            }
          />
        </DialogContent>
      </Dialog>

      {/* 刪除銀行帳戶確認 Dialog */}
      <AlertDialog
        open={!!deletingBankAccountId}
        onOpenChange={(open) => !open && setDeletingBankAccountId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認刪除</AlertDialogTitle>
            <AlertDialogDescription>
              確定要刪除此銀行帳戶嗎？此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteBankAccount}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              刪除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      {/* 刪除信用卡確認 Dialog */}
      <AlertDialog
        open={!!deletingCreditCardId}
        onOpenChange={(open) => !open && setDeletingCreditCardId(undefined)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>確認刪除</AlertDialogTitle>
            <AlertDialogDescription>
              確定要刪除此信用卡嗎？此操作無法復原。
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>取消</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDeleteCreditCard}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              刪除
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </AppLayout>
  );
}
